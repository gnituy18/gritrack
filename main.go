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
	host = "http://localhost"

	ErrUserNotLoggedIn = errors.New("user not logged in")
)

type User struct {
	Username string
	Birthday string
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

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./template/layout.html", "./template/index.html")
		if err != nil {
			log.Fatal(err)
		}

		if err = tmpl.Execute(w, nil); err != nil {
			log.Panic(err)
		}
	})

	http.HandleFunc("GET /snippet/{snippet}/{$}", func(w http.ResponseWriter, r *http.Request) {
		snippet := r.PathValue("snippet")
		if snippet == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles(fmt.Sprintf("./template/snippet/%s.html", snippet))
		if err != nil {
			log.Fatal(err)
		}
		tmpl, err = tmpl.Parse(fmt.Sprintf(`{{template "%s" . }}`, snippet))
		if err != nil {
			log.Fatal(err)
		}

		query := r.URL.Query()
		if err = tmpl.Execute(w, query); err != nil {
			log.Panic(err)
		}
	})

	http.HandleFunc("GET /sign-up/{$}", func(w http.ResponseWriter, r *http.Request) {
		if _, err := getSessionUser(r); err == nil {
			w.Header().Add("location", "/")
			w.WriteHeader(http.StatusFound)
			return
		}

		tmpl, err := template.ParseFiles("./template/layout.html", "./template/sign-up.html", "./template/snippet/sign-up-form.html")
		if err != nil {
			log.Fatal(err)
		}

		if err = tmpl.Execute(w, nil); err != nil {
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

		rows, err := tx.Query("SELECT username, email FROM user WHERE username = ? OR email = ?", username, email)
		if err != nil {
			tx.Rollback()
			log.Panic(err)
		}

		for rows.Next() {
			u, e := "", ""
			if err = rows.Scan(&u, &e); err != nil {
				tx.Rollback()
				log.Panic(err)
			}

			if u == username || e == email {
				tmpl, err := template.ParseFiles("./template/snippet/sign-up-form.html")
				if err != nil {
					log.Fatal(err)
				}

				tmpl, err = tmpl.Parse(`{{template "sign-up-form" . }}`)
				if err != nil {
					log.Fatal(err)
				}

				if err = tmpl.Execute(w, "Username or Email already exist."); err != nil {
					log.Panic(err)
				}

				tx.Rollback()
				return
			}

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
		username := ""
		if err := db.QueryRow("SELECT username FROM user WHERE email = ?", email).Scan(&username); err == sql.ErrNoRows {
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
		link := fmt.Sprintf("http://localhost:8080/log-in-with-token/?token=%s", id)
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

		cookie := http.Cookie{Name: "session", Value: token, Path: "/"}
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
