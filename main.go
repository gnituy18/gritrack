package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	_ "modernc.org/sqlite"
)

var (
	loc  = time.UTC
	port = "8080"
	host = "http://localhost:8080"

	ErrUserNotLoggedIn = errors.New("user not logged in")
)

type User struct {
	Username string
	Birthday string
}

type Track struct {
	Birthday time.Time
	Today    time.Time
	Days     [][]time.Time
}

type TmplPayload struct {
	User  *User
	Query *url.Values
	Track *Track
}

func main() {
	if lf, ok := os.LookupEnv("LOG_FILE"); ok {
		logFile, err := os.OpenFile(lf, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logFile)
	}

	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	if h, ok := os.LookupEnv("HOST"); ok {
		host = h
	}

	go func() {
		for {
			db := openDB()
			t := time.Now().Add(-7 * 24 * time.Hour)
			if _, err := db.Exec("DELETE FROM session WHERE created_at < ?", t.Format(time.DateTime)); err != nil {
				db.Close()
				time.Sleep(time.Minute)
				continue
			}

			db.Close()
			time.Sleep(time.Hour)
		}
	}()

	http.HandleFunc("GET /snippet/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		query := r.URL.Query()

		tmpl := snippet(name)
		var data any
		switch name {
		case "track":
			today, err := time.Parse(time.DateOnly, query.Get("today"))
			if err != nil {
				log.Panic(err)
			}

			birthday, err := time.Parse(time.DateOnly, "1993-12-20")
			if err != nil {
				log.Panic(err)
			}

			startDate := time.Date(birthday.Year(), birthday.Month(), 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(birthday.Year()+90, birthday.Month(), 1, 0, 0, 0, 0, time.UTC)
			currentDate := startDate
			days := [][]time.Time{}
			for currentDate.Before(endDate) {
				if currentDate.Day() == 1 {
					days = append(days, []time.Time{})
				}

				days[len(days)-1] = append(days[len(days)-1], currentDate)
				currentDate = currentDate.Add(24 * time.Hour)
			}

			data = Track{
				Today:    today,
				Birthday: birthday,
				Days:     days,
			}

		default:
			data = struct{}{}
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Panic(err)
		}
	})

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		user, err := getSessionUser(r)
		if err != nil && err != ErrUserNotLoggedIn {
			log.Panic(err)
		}

		payload := TmplPayload{
			User: user,
		}

		tmpl := template.Must(template.ParseFiles("./template/layout.html", "./template/index.html"))
		if err := tmpl.Execute(w, payload); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}
	})

	http.HandleFunc("GET /sign-up/{$}", func(w http.ResponseWriter, r *http.Request) {
		if _, err := getSessionUser(r); err == nil {
			w.Header().Add("location", "/")
			w.WriteHeader(http.StatusFound)
			return
		}

		tmpl := template.Must(template.ParseFiles("./template/layout.html", "./template/sign-up.html", "./template/snippet/sign-up-form.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			log.Panic(err)
		}
	})

	http.HandleFunc("POST /sign-up/{$}", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		email := r.FormValue("email")
		birthday := r.FormValue("birthday")

		db := openDB()
		defer db.Close()

		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Panic(err)
		}

		var count int
		if err = tx.QueryRow("SELECT count(*) FROM user WHERE username = ? OR email = ?", username, email).Scan(&count); err == sql.ErrNoRows {
			tx.Rollback()

			if err = snippet("sign-up-form").Execute(w, "Username or Email already exist."); err != nil {
				log.Panic(err)
			}

			return
		} else if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := tx.Exec("INSERT INTO user (username, email, birthday, public) VALUES (?, ?, ?, ?)", username, email, birthday, false); err != nil {
			tx.Rollback()
			log.Panic(err)
		}
		if _, err := tx.Exec("INSERT INTO goal (username, name) VALUES (?, ?)", username, "Your Goal!"); err != nil {
			tx.Rollback()
			log.Panic(err)
		}

		if err := tx.Commit(); err != nil {
			log.Panic(err)
		}

		tmpl, err := template.ParseFiles("./template/account-created.html")
		if err != nil {
			log.Fatal(err)
		}

		if err = tmpl.Execute(w, nil); err != nil {
			log.Panic(err)
		}
	})

	http.HandleFunc("GET /log-in/{$}", func(w http.ResponseWriter, r *http.Request) {
		if _, err := getSessionUser(r); err == nil {
			w.Header().Add("location", "/")
			w.WriteHeader(http.StatusFound)
			return
		}

		tmpl, err := template.ParseFiles("./template/layout.html", "./template/log-in.html")
		if err != nil {
			log.Fatal(err)
		}

		if err = tmpl.Execute(w, nil); err != nil {
			log.Panic(err)
		}
	})

	http.HandleFunc("POST /send-log-in-email/{$}", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		db := openDB()
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
		tmpl, err := template.ParseFiles("./template/log-in-email.html")
		if err != nil {
			log.Fatal(err)
		}
		link := fmt.Sprintf("%s/log-in-with-token/?token=%s", host, id)
		err = tmpl.Execute(&htmlBuffer, link)
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

		w.Write([]byte("log in email sent."))
	})

	http.HandleFunc("GET /log-in-with-token/{$}", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		token := query.Get("token")
		username := query.Get("username")

		db := openDB()
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

	http.HandleFunc("GET /styles.css/{$}", func(w http.ResponseWriter, r *http.Request) {
		asset, err := os.ReadFile("assets/styles.css")
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("Content-Type", "text/css")
		if _, err := w.Write(asset); err != nil {
			log.Panic(err)
		}
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func openDB() *sql.DB {
	db, err := sql.Open("sqlite", "./db")
	if err != nil {
		log.Panic(err)
	}

	return db
}

func snippet(name string) *template.Template {
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("./template/snippet/%s.html", name)))
	return template.Must(tmpl.Parse(fmt.Sprintf(`{{template "%s" . }}`, name)))
}

func getSessionUser(r *http.Request) (*User, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, ErrUserNotLoggedIn
	}

	db := openDB()
	defer db.Close()

	user := User{}
	if err := db.QueryRow("SELECT user.username, user.birthday FROM session JOIN user ON session.username = user.username WHERE session.id = ?", cookie.Value).Scan(&user.Username, &user.Birthday); err == sql.ErrNoRows {
		return nil, ErrUserNotLoggedIn
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
