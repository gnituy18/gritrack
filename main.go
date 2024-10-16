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
	"strconv"
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
	pageTmpl map[string]*template.Template
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
	pageTmpl = map[string]*template.Template{
		"index":   template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/index.gotmpl")),
		"app":     template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/app.gotmpl")),
		"log-in":  template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/log-in.gotmpl")),
		"sign-up": template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/sign-up.gotmpl")),
	}

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
		var data any

		tmplateName := r.PathValue("template_name")
		switch tmplateName {
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

		executeTemplate(w, tmplateName, data, "")
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
		trackerName := r.PathValue("tracker")

		var isPublic bool
		var b string
		user := User{}
		if err := db.QueryRow("SELECT user.email, user.birthday, user.timezone, tracker.public FROM user JOIN tracker ON user.username = tracker.username WHERE user.username = ? AND tracker.name = ?", username, trackerName).Scan(&user.Email, &b, &user.TimeZone, &isPublic); err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}
		user.Username = username
		birthday, _ := time.Parse(time.DateOnly, b)
		user.Birthday = birthday

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

		endY := sessionUser.Today().Year()
		endM := sessionUser.Today().Month()
		startY := endY - 1
		startM := endM + 1
		tracker, err := sessionUser.Tracker(trackerName, startY, startM, endY, endM)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		executePage(w, r, "app", App{
			SessionUser: sessionUser,
			User:        &user,
			Tracker:     tracker,
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
			fmt.Sprintf("update-day-%s", date),
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
			executeTemplate(w, "sign-up-form", "Username or Email already exist.", "")
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

		executeTemplate(w, "account-created", nil, "")
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
		err = tmpl.ExecuteTemplate(&htmlBuffer, "log-in-email", link)
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

	style, err := os.ReadFile("asset/style.css")
	if err != nil {
		log.Fatal(err)
	}

	var cssBuf bytes.Buffer
	minifyWriter := minifier.Writer("text/css", &cssBuf)
	if _, err := minifyWriter.Write(style); err != nil {
		log.Fatal(err)
	}
	minifyWriter.Close()
	css := cssBuf.Bytes()

	http.HandleFunc("GET /style.css/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		w.Header().Add("Cache-Control", "public, max-age=60")
		if _, err := w.Write(css); err != nil {
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

	page := pageTmpl[name]
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

	if err := tmpl.ExecuteTemplate(minifyWriter, name, data); err != nil {
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
	loc, _ := time.LoadLocation(u.TimeZone)
	y, m, d := time.Now().In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func (u *User) TimeRelation(date time.Time) TimeRelation {
	if date.Before(u.Today()) {
		return Past
	} else if date.After(u.Today()) {
		return Future
	} else {
		return Today
	}
}

func (u *User) Tracker(name string, startYear int, startMonth time.Month, endYear int, endMonth time.Month) (*Tracker, error) {
	startDate := time.Date(startYear, startMonth, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(endYear, endMonth+1, 1, 0, 0, 0, 0, time.UTC)
	rows, err := db.Query("SELECT date, content FROM day WHERE username = ? AND tracker_name = ? AND date >= ? AND date < ?", u.Username, name, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dayContentMap := map[string]string{}
	for rows.Next() {
		var date string
		var content string
		rows.Scan(&date, &content)
		dayContentMap[date] = content
	}

	years := []Year{}
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 1, 0) {

		if len(years) == 0 || d.Year() != years[0].Value {
			years = append([]Year{{
				Value:  d.Year(),
				Months: []Month{},
			}}, years...)
		}

		years[0].Months = append([]Month{{
			Value: d.Month(),
			Days:  []Day{},
		}}, years[0].Months...)

		nextMonth := d.AddDate(0, 1, 0)
		for day := d; day.Before(nextMonth); day = day.AddDate(0, 0, 1) {
			content := ""
			if c, ok := dayContentMap[day.Format(time.DateOnly)]; ok {
				content = c
			}

			timeRelation := u.TimeRelation(day)

			years[0].Months[0].Days = append(years[0].Months[0].Days, Day{
				Date:         day,
				Content:      content,
				TimeRelation: timeRelation,
			})
		}
	}

	return &Tracker{
		Name:  name,
		Years: years,
	}, nil
}

func getSessionUser(r *http.Request) (*User, bool, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, false, nil
	}

	user := User{Trackers: []string{}}
	var b string
	if err := db.QueryRow("SELECT user.username, user.email, user.birthday, user.timezone FROM user JOIN session ON session.username = user.username WHERE session.id = ?", cookie.Value).Scan(&user.Username, &user.Email, &b, &user.TimeZone); err == sql.ErrNoRows {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	birthday, _ := time.Parse(time.DateOnly, b)
	user.Birthday = birthday

	rows, err := db.Query("SELECT name FROM tracker WHERE tracker.username = ? ORDER BY position", user.Username)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

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
	User        *User

	Tracker *Tracker
}

type Tracker struct {
	Name  string
	Years []Year
}

func (t *Tracker) String() string {
	return t.Name
}

type Year struct {
	Value  int
	Months []Month
}

func (y Year) String() string {
	return strconv.Itoa(y.Value)
}

type Month struct {
	Value time.Month
	Days  []Day
}

func (m Month) String() string {
	return m.Value.String()
}

type Day struct {
	Date         time.Time
	Content      string
	TimeRelation TimeRelation
}

func (d Day) Week() int {
	firstWeekday := int(time.Date(d.Date.Year(), d.Date.Month(), 1, 0, 0, 0, 0, time.UTC).Weekday())
	return (firstWeekday+d.Date.Day()-1)/7 + 1
}

func (d Day) Weekday() int {
	return int(d.Date.Weekday()) + 1
}

func (d Day) String() string {
	return d.Date.Format(time.DateOnly)
}

const (
	Past   TimeRelation = "Past"
	Today  TimeRelation = "Today"
	Future TimeRelation = "Future"
)

type TimeRelation string

func (tr TimeRelation) IsPast() bool {
	return tr == Past
}

func (tr TimeRelation) IsToday() bool {
	return tr == Today
}

func (tr TimeRelation) IsFuture() bool {
	return tr == Future
}
