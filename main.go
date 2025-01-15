package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
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
	gosimpleSlug "github.com/gosimple/slug"
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
	pageTmpl map[string]*template.Template
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
	db.Exec("PRAGMA journal_mode = WAL")

	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for {
			<-ticker.C
			cutoff := time.Now().Add(-7 * 24 * time.Hour).Unix()
			if _, err := db.Exec("DELETE FROM user_sessions WHERE created_at < ?", cutoff); err != nil {
				log.Printf("deleting user sessions failed: %v\n", err)
				continue
			}
		}
	}()

	tmpl = template.Must(template.New("base").Funcs(sprig.FuncMap()).ParseGlob("./template/*.gotmpl"))
	pageTmpl = map[string]*template.Template{
		"index":            template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/index.gotmpl")),
		"app":              template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/app.gotmpl")),
		"owner-tracker":    template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/owner-tracker.gotmpl")),
		"create-tracker":   template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/create-tracker.gotmpl")),
		"log-in":           template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/log-in.gotmpl")),
		"sign-up":          template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/sign-up.gotmpl")),
		"settings":         template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/settings.gotmpl")),
		"settings-tracker": template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/settings-tracker.gotmpl")),
		"email-sent":       template.Must(template.Must(tmpl.Clone()).ParseFiles("./page/email-sent.gotmpl")),
	}

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
		case "months":
			username := query.Get("username")
			slug := query.Get("slug")

			user := User{}
			tracker := Tracker{}
			if err := db.QueryRow("SELECT users.email, users.timezone, trackers.description, trackers.position, trackers.public FROM users JOIN trackers ON users.username = trackers.username WHERE users.username = ? AND trackers.slug = ?", username, slug).Scan(&user.Email, &user.TimeZone, &tracker.Description, &tracker.Position, &tracker.Public); err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Panic(err)
			}
			user.Username = username
			tracker.Slug = slug

			if !tracker.Public {
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if username != sessionUser.Username {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}

			query := r.URL.Query()
			from := query.Get("from")

			fromY := sessionUser.Today().Year()
			fromM := sessionUser.Today().Month()
			if from != "" {
				t, err := time.Parse(time.DateOnly, from+"-01")
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				fromY = t.Year()
				fromM = t.Month()
			}

			toY := fromY
			toM := fromM - time.Month(3) + 1
			to := query.Get("to")
			if to != "" {
				t, err := time.Parse(time.DateOnly, to+"-01")
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				toY = t.Year()
				toM = t.Month()
			}

			trackerEntries, err := sessionUser.TrackerEntries(slug, fromY, fromM, toY, toM)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Panic(err)
			}

			data = map[string]any{
				"tracker": tracker,
				"entries": trackerEntries,
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

		executeTemplates(w, data, "", tmplateName)
	})

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			executePage(w, r, "index", nil)
			return
		}

		var username string
		if err := db.QueryRow(`
			SELECT users.username
			FROM users
			JOIN user_sessions ON users.username = user_sessions.username
			WHERE user_sessions.id = ?
		`, cookie.Value).Scan(&username); err == sql.ErrNoRows {
			executePage(w, r, "index", nil)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		http.Redirect(w, r, fmt.Sprintf("/%s/", username), http.StatusFound)
	})

	// user handler
	http.HandleFunc("GET /{username}/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if sessionUser.Username == r.PathValue("username") {
			daysArr, err := sessionUser.PastDays(8)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			executePage(w, r, "app", map[string]any{
				"sessionUser": sessionUser,
				"daysArr":     daysArr,
			})
			return
		}

		w.WriteHeader(http.StatusForbidden)
	})

	http.HandleFunc("GET /{username}/{slug}/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		username := r.PathValue("username")
		slug := r.PathValue("slug")

		tracker := sessionUser.Tracker(slug)

		if !tracker.Public {
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if username != sessionUser.Username {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		query := r.URL.Query()
		to := query.Get("to")

		toY := sessionUser.Today().Year()
		toM := sessionUser.Today().Month()
		if to != "" {
			t, err := time.Parse(time.DateOnly, to+"-01")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			toY = t.Year()
			toM = t.Month()
		}

		fromY := toY
		fromM := toM - time.Month(3) + 1
		from := query.Get("from")
		if from != "" {
			t, err := time.Parse(time.DateOnly, from+"-01")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fromY = t.Year()
			fromM = t.Month()
		}

		trackerEntries, err := sessionUser.TrackerEntries(slug, fromY, fromM, toY, toM)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		executePage(w, r, "owner-tracker", map[string]any{
			"sessionUser": sessionUser,
			"tracker":     &tracker,
			"entries":     trackerEntries,
		})
	})

	http.HandleFunc("GET /settings/{slug}/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		slug := r.PathValue("slug")

		tracker := sessionUser.Tracker(slug)
		if tracker == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		executePage(w, r, "settings-tracker", map[string]any{
			"sessionUser": sessionUser,
			"tracker":     tracker,
		})
	})

	http.HandleFunc("PATCH /settings/{slug}/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		slug := r.PathValue("slug")

		tracker := sessionUser.Tracker(slug)
		if tracker == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		displayName := r.FormValue("display_name")
		description := r.FormValue("description")

		if _, err := db.Exec(`
			UPDATE trackers
			SET display_name = ?, description = ?
			WHERE username = ?
			AND slug = ?
			`, displayName, description, sessionUser.Username, slug); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		w.Header().Add("HX-Trigger", "success")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("PATCH /settings/{slug}/slug/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		slug := r.PathValue("slug")

		tracker := sessionUser.Tracker(slug)
		if tracker == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		newSlug := r.FormValue("slug")

		if _, err := db.Exec(`
			UPDATE trackers
			SET slug = ?
			WHERE username = ?
			AND slug = ?
			`, newSlug, sessionUser.Username, slug); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		w.Header().Add("HX-Redirect", fmt.Sprintf("/settings/%s/", newSlug))
		w.WriteHeader(http.StatusSeeOther)
	})

	http.HandleFunc("DELETE /settings/{slug}/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		slug := r.PathValue("slug")

		tracker := sessionUser.Tracker(slug)
		if tracker == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if _, err := db.Exec(`
			DELETE FROM trackers
			WHERE username = ?
			AND slug = ?
			`, sessionUser.Username, slug); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		w.Header().Add("HX-Redirect", fmt.Sprintf("/%s/", sessionUser.Username))
		w.WriteHeader(http.StatusSeeOther)
	})

	http.HandleFunc("GET /day-detail/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		query := r.URL.Query()
		slug := query.Get("slug")
		date := query.Get("date")

		if sessionUser.Tracker(slug) == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if _, err := time.Parse(time.DateOnly, date); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var emoji string
		var content string
		if err := db.QueryRow("SELECT emoji, content FROM tracker_entries WHERE username = ? AND slug = ? AND date = ?", sessionUser.Username, slug, date).Scan(&emoji, &content); err != nil && err != sql.ErrNoRows {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		m := map[string]string{
			"emoji":   emoji,
			"content": content,
		}

		json.NewEncoder(w).Encode(m)
	})

	http.HandleFunc("PUT /day-detail/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		slug := r.FormValue("slug")
		date := r.FormValue("date")
		emoji := r.FormValue("emoji")
		content := r.FormValue("content")

		if sessionUser.Tracker(slug) == nil {
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

		if _, err := db.Exec("INSERT INTO tracker_entries (username, slug, date, emoji, content) VALUES (?, ?, ?, ?, ?) ON CONFLICT DO UPDATE SET emoji = ?, content = ?", sessionUser.Username, slug, date, emoji, content, emoji, content); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		tracker := Tracker{
			Slug: slug,
		}
		day := Day{
			Date: t,
		}
		if err := db.QueryRow(`
				SELECT
				trackers.display_name,
				trackers.position,
				trackers.public,
				tracker_entries.emoji,
				tracker_entries.content 
				FROM tracker_entries 
				INNER JOIN trackers
				ON tracker_entries.username = trackers.username
				AND tracker_entries.slug = trackers.slug
				WHERE tracker_entries.username = ? 
				AND tracker_entries.slug = ? 
				AND tracker_entries.date = ?
			`, sessionUser.Username, slug, date).Scan(
			&tracker.DisplayName,
			&tracker.Position,
			&tracker.Public,
			&day.Emoji,
			&day.Content,
		); err != nil && err != sql.ErrNoRows {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		day.TimeRelation = sessionUser.TimeRelation(t)

		data := map[string]any{
			"tracker": tracker,
			"day":     day,
			"oob":     true,
		}

		executeTemplates(w, data, "", "day", "today-preview")
		return
	})

	http.HandleFunc("GET /create-tracker/{$}", func(w http.ResponseWriter, r *http.Request) {
		if _, ok, err := getSessionUser(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		executePage(w, r, "create-tracker", nil)
	})

	http.HandleFunc("POST /create-tracker/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		displayName := strings.TrimSpace(r.FormValue("display_name"))
		if displayName == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		baseSlug := gosimpleSlug.Make(displayName)
		slug := baseSlug

		suffix := 1

		for {
			_, err := db.Exec(
				"INSERT INTO trackers (username, slug, display_name, position) VALUES (?, ?, ?, (SELECT COALESCE(MAX(position), 0) + 1 FROM trackers WHERE username = ?))",
				sessionUser.Username, slug, displayName, sessionUser.Username,
			)

			if err == nil {
				break
			}

			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				slug = fmt.Sprintf("%s-%d", baseSlug, suffix)
				suffix++
				continue
			}

			log.Panic(err)
		}

		w.Header().Add("HX-Redirect", fmt.Sprintf("/%s/", sessionUser.Username))
		w.WriteHeader(http.StatusSeeOther)
	})

	http.HandleFunc("GET /settings/{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionUser, ok, err := getSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		} else if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		executePage(w, r, "settings", sessionUser)
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
		tz := r.FormValue("timezone")

		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Panic(err)
		}

		var count int
		if err := tx.QueryRow("SELECT count(*) FROM users WHERE username = ? OR email = ?", username, email).Scan(&count); err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if count != 0 {
			tx.Rollback()
			w.Write([]byte("⚠️ Username or Email already exist."))
			return
		}

		if _, err := tx.Exec("INSERT INTO users (username, email, timezone) VALUES (?, ?, ?)", username, email, tz); err != nil {
			tx.Rollback()
			log.Panic(err)
		}

		if err := tx.Commit(); err != nil {
			log.Panic(err)
		}

		if _, err := sendLogInEmail(email, username); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		w.Header().Add("HX-Location", "/account-created/")
		w.WriteHeader(http.StatusSeeOther)
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
		if err := db.QueryRow("SELECT username FROM users WHERE email = ?", email).Scan(&username); err == sql.ErrNoRows {
			w.Write([]byte("You don't have an account yet."))
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		if _, err := sendLogInEmail(email, username); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		w.Header().Add("HX-Location", "/log-in-email-sent/")
		w.WriteHeader(http.StatusSeeOther)
	})

	http.HandleFunc("GET /account-created/{$}", func(w http.ResponseWriter, r *http.Request) {
		executePage(w, r, "email-sent", "🎉 Account Successfully Created! 🎉")
	})

	http.HandleFunc("GET /log-in-email-sent/{$}", func(w http.ResponseWriter, r *http.Request) {
		executePage(w, r, "email-sent", "📧 Log In Email Sent! 📧")
	})

	http.HandleFunc("GET /log-in-with-token/{$}", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		token := query.Get("token")
		username := query.Get("username")

		if err := db.QueryRow("SELECT * FROM user_sessions WHERE id = ? AND username = ?", token, username).Err(); err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Panic(err)
		}

		cookie := http.Cookie{Name: "session", Value: token, Path: "/", Expires: time.Now().Add(7 * 24 * time.Hour)}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	style, err := os.ReadFile("asset/style.css")
	if err != nil {
		log.Fatal(err)
	}

	minifier = minify.New()
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/css", css.Minify)

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

	http.HandleFunc("GET /logo.svg/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")

		svg := `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="9 9 62 62" width="62" height="62">
  <rect x="0" y="0" width="24" height="24" rx="4" ry="4" />
  <rect x="28" y="0" width="24" height="24" rx="4" ry="4" />
  <rect x="56" y="0" width="24" height="24" rx="4" ry="4" />
  <rect x="0" y="28" width="24" height="24" rx="4" ry="4"  />
  <rect x="28" y="28" width="24" height="24" rx="4" ry="4"  />
  <rect x="56" y="28" width="24" height="24" rx="4" ry="4"  />
  <rect x="0" y="56" width="24" height="24" rx="4" ry="4"  />
  <rect x="28" y="56" width="24" height="24" rx="4" ry="4"  />
  <rect x="56" y="56" width="24" height="24" rx="4" ry="4" />
</svg>
`
		w.Write([]byte(svg))
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

func executeTemplates(w http.ResponseWriter, data any, trigger string, names ...string) {
	minifyWriter := minifier.Writer("text/html", w)
	defer minifyWriter.Close()

	if len(trigger) > 0 {
		w.Header().Add("HX-Trigger", trigger)
	}

	for _, name := range names {
		if err := tmpl.ExecuteTemplate(minifyWriter, name, data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}
	}
}

func sendLogInEmail(email, username string) (*ses.SendEmailOutput, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	bs := make([]byte, 32)
	if _, err = rand.Read(bs); err != nil {
		log.Panic(err)
	}
	sessionId := base64.URLEncoding.EncodeToString(bs)

	if _, err := db.Exec("INSERT INTO user_sessions (id, username) VALUES (?, ?)", sessionId, username); err != nil {
		log.Panic(err)
	}

	from := "no-reply@gritrack.com"
	title := "Log in to Gritrack"
	var htmlBuffer bytes.Buffer

	portStr := ""
	if port != "80" {
		portStr = ":" + port
	}

	link := fmt.Sprintf("%s%s/log-in-with-token/?token=%s", host, portStr, sessionId)
	err = tmpl.ExecuteTemplate(&htmlBuffer, "log-in-email", link)
	if err != nil {
		log.Panic(err)
	}
	body := htmlBuffer.String()

	client := ses.NewFromConfig(cfg)
	return client.SendEmail(ctx, &ses.SendEmailInput{
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
}

func getSessionUser(r *http.Request) (*User, bool, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, false, nil
	}

	var user *User
	rows, err := db.Query(`
		SELECT
		users.username,
		users.email,
		users.timezone,
		users.public,
		trackers.slug,
		trackers.display_name,
		trackers.description,
		trackers.position,
		trackers.public
		FROM users
		INNER JOIN user_sessions ON users.username = user_sessions.username
		LEFT JOIN trackers ON trackers.username = users.username
		WHERE user_sessions.id = ?
		ORDER BY trackers.position`, cookie.Value)

	if err != nil {
		return nil, false, err
	}

	for rows.Next() {
		var username, email, timeZone string
		var userPublic bool
		var slug, displayName, description string
		var position int
		var trackerPublic bool

		rows.Scan(
			&username,
			&email,
			&timeZone,
			&userPublic,
			&slug,
			&displayName,
			&description,
			&position,
			&trackerPublic,
		)

		if user == nil {
			user = &User{
				Username: username,
				Email:    email,
				TimeZone: timeZone,
				Public:   userPublic,
				Trackers: []Tracker{},
			}
		}

		if slug != "" {
			user.Trackers = append(user.Trackers, Tracker{
				Slug:        slug,
				DisplayName: displayName,
				Description: description,
				Position:    position,
				Public:      trackerPublic,
			})
		}
	}

	if user == nil {
		return nil, false, nil
	}

	return user, true, nil
}

type User struct {
	Username string
	Email    string
	TimeZone string
	Public   bool

	Trackers []Tracker
}

func (u *User) Tracker(slug string) *Tracker {
	for _, t := range u.Trackers {
		if t.Slug == slug {
			return &t
		}
	}

	return nil
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

func (u *User) PastDays(days int) (map[string][]*Day, error) {
	endDate := u.Today()
	startDate := endDate.AddDate(0, 0, -(days - 1))
	rows, err := db.Query(`
		SELECT
		slug,
		date,
		emoji,
		content
		FROM tracker_entries
		WHERE username = ? AND date >= ? AND date <= ?
		ORDER BY slug, date
		`, u.Username, startDate.Format(time.DateOnly), endDate.Format(time.DateOnly))
	if err != nil {
		return nil, err
	}

	daysMap := map[string]map[string]*Day{}
	daysArr := map[string][]*Day{}
	for _, t := range u.Trackers {
		daysMap[t.Slug] = map[string]*Day{}
		daysArr[t.Slug] = []*Day{}
	}

	for rows.Next() {
		var slug, date, content, emoji string
		rows.Scan(&slug, &date, &emoji, &content)
		t, err := time.Parse(time.DateOnly, date)
		if err != nil {
			return nil, err
		}
		timeRelation := u.TimeRelation(t)
		daysMap[slug][date] = &Day{
			Date:         t,
			Emoji:        emoji,
			Content:      content,
			TimeRelation: timeRelation,
		}
	}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		for _, t := range u.Trackers {
			if d := daysMap[t.Slug][d.Format(time.DateOnly)]; d != nil {
				daysArr[t.Slug] = append(daysArr[t.Slug], d)
				continue
			}

			daysArr[t.Slug] = append(daysArr[t.Slug], &Day{
				Date:         d,
				TimeRelation: u.TimeRelation(d),
			})
		}
	}

	return daysArr, nil
}

func (u *User) TrackerEntries(slug string, fromYear int, fromMonth time.Month, toYear int, toMonth time.Month) (*TrackerEntries, error) {
	startDate := time.Date(fromYear, fromMonth, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(toYear, toMonth+1, 1, 0, 0, 0, 0, time.UTC)
	rows, err := db.Query("SELECT date, emoji, content FROM tracker_entries WHERE username = ? AND slug = ? AND date >= ? AND date < ?", u.Username, slug, startDate.Format(time.DateOnly), endDate.Format(time.DateOnly))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	daysMap := map[string]struct {
		Content string
		Emoji   string
	}{}
	for rows.Next() {
		var date string
		var content string
		var emoji string
		rows.Scan(&date, &emoji, &content)
		daysMap[date] = struct {
			Content string
			Emoji   string
		}{Emoji: emoji, Content: content}
	}

	months := []Month{}
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 1, 0) {
		nextMonth := d.AddDate(0, 1, 0)
		m := Month{
			Year:  d.Year(),
			Month: d.Month(),
			Days:  []Day{},
		}
		for day := d; day.Before(nextMonth); day = day.AddDate(0, 0, 1) {
			content := ""
			emoji := ""
			if data, ok := daysMap[day.Format(time.DateOnly)]; ok {
				content = data.Content
				emoji = data.Emoji
			}

			timeRelation := u.TimeRelation(day)

			m.Days = append(m.Days, Day{
				Date:         day,
				Content:      content,
				Emoji:        emoji,
				TimeRelation: timeRelation,
			})
		}

		months = append([]Month{m}, months...)
	}

	return &TrackerEntries{
		Months: months,
	}, nil
}

type Tracker struct {
	Slug        string
	DisplayName string
	Description string
	Position    int
	Public      bool
}

func (t *Tracker) String() string {
	return t.DisplayName
}

type TrackerEntries struct {
	Months []Month
}

type Month struct {
	Year  int
	Month time.Month
	Days  []Day
}

func (m Month) Weeks() int {
	return m.Days[len(m.Days)-1].Week()
}

func (m Month) FormatYYYYMM() string {
	return fmt.Sprintf("%d-%02d", m.Year, m.Month)
}

func (m Month) FormatTwoDigitMonth() string {
	return fmt.Sprintf("%02d", m.Month)
}

func (m Month) FormatMonthName() string {
	return m.Month.String()[0:3]
}

type Day struct {
	Date         time.Time
	Emoji        string
	Content      string
	TimeRelation TimeRelation
}

func (d Day) WeekdayString() string {
	return d.Date.Weekday().String()[0:3]
}

func (d Day) Week() int {
	firstWeekday := int(time.Date(d.Date.Year(), d.Date.Month(), 1, 0, 0, 0, 0, time.UTC).Weekday())
	return (firstWeekday+d.Date.Day()-1)/7 + 2
}

func (d Day) Weekday() int {
	return int(d.Date.Weekday()) + 1
}

func (d Day) String() string {
	return d.Date.Format(time.DateOnly)
}

func (d Day) Set() bool {
	return !(d.Content == "" && d.Emoji == "")
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
