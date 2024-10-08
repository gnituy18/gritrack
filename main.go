package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	_ "modernc.org/sqlite"
)

var (
	port = "8080"
	host = "http://localhost"

	db       *sql.DB
	tmpl     *template.Template
	minifier *minify.M
)

func main() {
	if v, ok := os.LookupEnv("LOG_FILE"); ok {
		logFile, err := os.OpenFile(v, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logFile)
	}

	if v, ok := os.LookupEnv("PORT"); ok {
		port = v
	}

	if v, ok := os.LookupEnv("HOST"); ok {
		host = v
	}

	var err error
	if db, err = sql.Open("sqlite", "./db"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Exec("PRAGMA foreign_keys = ON")

	go func() {
		for {
			t := time.Now().Add(-7 * 24 * time.Hour)
			if _, err := db.Exec("DELETE FROM session WHERE created_at < ?", t.Format(time.DateTime)); err != nil {
				time.Sleep(time.Minute)
				continue
			}

			time.Sleep(time.Hour)
		}
	}()

	tmpl = template.Must(template.New("base").Funcs(sprig.FuncMap()).ParseGlob("./template/*.gotmpl"))

	minifier = minify.New()
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/css", css.Minify)

	http.HandleFunc("GET /template/app/{template_name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		query := r.URL.Query()
		tmplName := r.PathValue("template_name")

		var data any
		switch tmplName {
		case "day":
			fallthrough
		case "day-detail":
			tracker := query.Get("tracker")
			date := query.Get("date")
			t, err := time.Parse(time.DateOnly, date)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var content string
			if err := db.QueryRow("SELECT content FROM day WHERE username = ? AND tracker_name = ? AND date = ?", sessionUser.Username, tracker, date).Scan(&content); err != nil && err != sql.ErrNoRows {
				w.WriteHeader(http.StatusInternalServerError)
				log.Panic(err)
			}

			data = map[string]any{
				"tracker": tracker,
				"day": Day{
					Date:         t,
					Content:      content,
					TimeRelation: sessionUser.TimeRelation(t),
				},
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

		executeTemplate(w, tmplName, data, "")
	})

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			executePage(w, r, "index", nil)
			return
		}

		var username, tracker string
		if err := db.QueryRow("SELECT user.username, tracker.name FROM user JOIN session ON user.username = session.username JOIN tracker ON user.username = tracker.username WHERE session.id = ? AND tracker.position = 1", cookie.Value).Scan(&username, &tracker); err == sql.ErrNoRows {
			executePage(w, r, "index", nil)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		http.Redirect(w, r, fmt.Sprintf("/%s/%s/", username, tracker), http.StatusTemporaryRedirect)
	})

	http.HandleFunc("GET /{username}/{tracker}/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		username := r.PathValue("username")
		tracker := r.PathValue("tracker")

		var isPublic bool
		var birthday string
		user := User{}
		if err := db.QueryRow("SELECT user.email, user.birthday, user.timezone, tracker.public FROM user JOIN tracker ON user.username = tracker.username WHERE user.username = ? AND tracker.name = ?", username, tracker).Scan(&user.Email, &birthday, &user.TimeZone, &isPublic); err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}
		user.Username = username
		t, _ := time.Parse(time.DateOnly, birthday)
		user.Birthday = t

		if !isPublic {
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if username != sessionUser.Username {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		rows, err := db.Query("SELECT date, content FROM day WHERE username = ? AND tracker_name = ?", username, tracker)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		m := map[string]string{}
		for rows.Next() {
			var date string
			var content string
			rows.Scan(&date, &content)
			m[date] = content
		}

		today := user.Today()
		startDate := time.Date(today.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := time.Date(today.Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
		currentDate := startDate
		years := []Year{
			{
				Months: []Month{},
			},
		}
		for currentDate.Before(endOfYear) {
			if currentDate.Day() == 1 {
				if currentDate.Month() == 1 {
					years = append([]Year{
						{
							Months: []Month{},
						},
					}, years...)
				}

				years[0].Months = append([]Month{{Days: []Day{}}}, years[0].Months...)
			}

			content := ""
			if c, ok := m[currentDate.Format(time.DateOnly)]; ok {
				content = c
			}

			timeRelation := user.TimeRelation(currentDate)

			years[0].Months[0].Days = append(years[0].Months[0].Days, Day{
				Date:         currentDate,
				Content:      content,
				TimeRelation: timeRelation,
			})

			currentDate = currentDate.Add(24 * time.Hour)
		}

		executePage(w, r, "app", App{
			SessionUser: sessionUser,
			User:        &user,
			Tracker:     tracker,
			Years:       years,
		})
	})

	http.HandleFunc("PUT /{tracker}/{date}/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tracker := r.PathValue("tracker")
		date := r.PathValue("date")
		content := r.FormValue("content")

		if !sessionUser.HasTracker(tracker) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		t, err := time.Parse(time.DateOnly, date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if t.After(sessionUser.Today()) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if _, err := db.Exec("INSERT INTO day (username, tracker_name, date, content) VALUES (?, ?, ?, ?) ON CONFLICT DO UPDATE SET content = ?", sessionUser.Username, tracker, date, content, content); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		executeTemplate(w, "day-detail", map[string]any{
			"tracker": tracker,
			"day": Day{
				Date:         t,
				Content:      content,
				TimeRelation: sessionUser.TimeRelation(t),
			},
		},
			fmt.Sprintf("update-d-%s", date),
		)
		return
	})

	http.HandleFunc("GET /sign-up/{$}", func(w http.ResponseWriter, r *http.Request) {
		if _, ok, err := getSessionUser(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if ok {
			w.Header().Add("location", "/")
			w.WriteHeader(http.StatusFound)
			return
		}

		executePage(w, r, "sign-up", nil)
	})

	http.HandleFunc("POST /sign-up/{$}", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		email := r.FormValue("email")
		birthday := r.FormValue("birthday")
		tz := r.FormValue("timezone")

		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Panic(err)
		}

		var count int
		if err := tx.QueryRow("SELECT count(*) FROM user WHERE username = ? OR email = ?", username, email).Scan(&count); err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if count != 0 {
			tx.Rollback()

			if err := template.Must(tmpl.Clone()).ExecuteTemplate(w, "sign-up-form", "Username or Email already exist."); err != nil {
				log.Panic(err)
			}

			return
		}

		if _, err := tx.Exec("INSERT INTO user (username, email, birthday, timezone) VALUES (?, ?, ?, ?)", username, email, birthday, tz); err != nil {
			tx.Rollback()
			log.Panic(err)
		}
		if _, err := tx.Exec("INSERT INTO tracker (username, name, position) VALUES (?, ?, 1)", username, "Your_First_Tracker"); err != nil {
			tx.Rollback()
			log.Panic(err)
		}

		if err := tx.Commit(); err != nil {
			log.Panic(err)
		}

		template.Must(tmpl.Clone()).ExecuteTemplate(w, "account-created", nil)
	})

	http.HandleFunc("GET /log-in/{$}", func(w http.ResponseWriter, r *http.Request) {
		if _, ok, err := getSessionUser(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if ok {
			w.Header().Add("location", "/")
			w.WriteHeader(http.StatusFound)
			return
		}

		executePage(w, r, "log-in", nil)
	})

	http.HandleFunc("POST /send-log-in-email/{$}", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		var username string
		if err := db.QueryRow("SELECT username FROM user WHERE email = ?", email).Scan(&username); err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user not exist"))
			return
		} else if err != nil {
			log.Panic(err)
		}

		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Panic(err)
			return
		}

		bs := make([]byte, 32)
		if _, err = rand.Read(bs); err != nil {
			log.Panic(err)
		}
		id := base64.URLEncoding.EncodeToString(bs)

		if _, err := db.Exec("INSERT INTO session (id, username) VALUES (?, ?)", id, username); err != nil {
			log.Panic(err)
		}

		from := "no-reply@gritrack.com"
		title := "Log in to Gritrack"
		var htmlBuffer bytes.Buffer

		portStr := ""
		if port != "80" {
			portStr = ":" + port
		}

		link := fmt.Sprintf("%s%s/log-in-with-token/?token=%s", host, portStr, id)
		err = template.Must(tmpl.Clone()).ExecuteTemplate(&htmlBuffer, "log-in-email", link)
		if err != nil {
			log.Panic(err)
		}

		body := htmlBuffer.String()

		client := ses.NewFromConfig(cfg)
		_, err = client.SendEmail(ctx, &ses.SendEmailInput{
			Destination: &types.Destination{
				ToAddresses: []string{email},
			},
			Message: &types.Message{
				Subject: &types.Content{
					Data: &title,
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: &body,
					},
				},
			},
			Source: &from,
		})

		w.Write([]byte("Log in email sent. Check your email."))
	})

	http.HandleFunc("GET /log-in-with-token/{$}", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		token := query.Get("token")
		username := query.Get("username")

		if err := db.QueryRow("SELECT * FROM session WHERE id = ? AND username = ?", token, username).Err(); err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Panic(err)
		}

		cookie := http.Cookie{Name: "session", Value: token, Path: "/", Expires: time.Now().Add(7 * 24 * time.Hour)}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	http.HandleFunc("GET /style.css/{$}", func(w http.ResponseWriter, r *http.Request) {
		asset, err := os.ReadFile("asset/style.css")
		if err != nil {
			log.Fatal(err)
		}

		var writer io.Writer = w

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			writer = gzipWriter
		}

		minifyWriter := minifier.Writer("text/css", writer)
		defer minifyWriter.Close()

		w.Header().Add("Content-Type", "text/css")
		w.Header().Add("Cache-Control", "public, max-age=60")
		if _, err := minifyWriter.Write(asset); err != nil {
			log.Panic(err)
		}
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func executePage(w http.ResponseWriter, r *http.Request, name string, data any) {
	var writer io.Writer = w

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		writer = gzipWriter
	}

	minifyWriter := minifier.Writer("text/html", writer)
	defer minifyWriter.Close()

	page := template.Must(template.Must(tmpl.Clone()).ParseFiles(fmt.Sprintf("./page/%s.gotmpl", name)))
	if err := page.ExecuteTemplate(minifyWriter, "page", data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}
}

func executeTemplate(w http.ResponseWriter, name string, data any, trigger string) {
	minifyWriter := minifier.Writer("text/html", w)
	defer minifyWriter.Close()

	if len(trigger) > 0 {
		w.Header().Add("HX-Trigger", trigger)
	}

	if err := template.Must(tmpl.Clone()).ExecuteTemplate(minifyWriter, name, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}
}

type User struct {
	Username string
	Birthday time.Time
	Email    string
	TimeZone string

	Trackers []string
}

func (u *User) HasTracker(tracker string) bool {
	for _, t := range u.Trackers {
		if t == tracker {
			return true
		}
	}

	return false
}

func (u *User) Today() time.Time {
	return DayStartUTC(u.TimeZone)
}

func (u *User) TimeRelation(date time.Time) TimeRelation {
	if date.Before(u.Birthday) {
		return Prenatal
	} else if date.Before(u.Today()) {
		return Past
	} else if date.After(u.Today()) {
		return Future
	} else {
		return Today
	}
}

func DayStartUTC(tz string) time.Time {
	loc, _ := time.LoadLocation(tz)
	y, m, d := time.Now().In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func getSessionUser(r *http.Request) (*User, bool, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, false, nil
	}

	user := User{Trackers: []string{}}
	var birthday string
	if err := db.QueryRow("SELECT user.username, user.email, user.birthday, user.timezone FROM user JOIN session ON session.username = user.username WHERE session.id = ?", cookie.Value).Scan(&user.Username, &user.Email, &birthday, &user.TimeZone); err == sql.ErrNoRows {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	t, _ := time.Parse(time.DateOnly, birthday)
	user.Birthday = t

	rows, err := db.Query("SELECT name FROM tracker WHERE tracker.username = ? ORDER BY position", user.Username)
	if err != nil {
		return nil, false, err
	}

	for rows.Next() {
		var tracker string
		if err := rows.Scan(&tracker); err != nil {
			return nil, false, err
		}

		user.Trackers = append(user.Trackers, tracker)
	}

	return &user, true, nil
}

type App struct {
	SessionUser *User

	User    *User
	Tracker string
	Years   []Year
}

type Year struct {
	Months []Month
}

type Month struct {
	Days []Day
}

func (m Month) String() string {
	return m.Days[0].String()[:7]
}

type TimeRelation string

const (
	Prenatal TimeRelation = "Prenatal"
	Past     TimeRelation = "Past"
	Today    TimeRelation = "Today"
	Future   TimeRelation = "Future"
)

type Day struct {
	Date    time.Time
	Content string

	TimeRelation TimeRelation
}

func (d Day) String() string {
	return d.Date.Format(time.DateOnly)
}

func (d Day) IsPrenatal() bool {
	return d.TimeRelation == Prenatal
}

func (d Day) IsPast() bool {
	return d.TimeRelation == Past
}

func (d Day) IsToday() bool {
	return d.TimeRelation == Today
}

func (d Day) IsFuture() bool {
	return d.TimeRelation == Future
}
