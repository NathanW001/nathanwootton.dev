package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type BlogPost struct {
	Title       string
	Subtitle    string
	Date        string // Thought this field is the time (and stored by the db in UNIX time), passing a time.Time value to the template will give the default tostring method which looks ugly so I've chosen to convert it myself
	Readingmins int
	Body        []string
}

func main() {
	// First, we connect to the DB and exit if there are errors
	db, err := sql.Open("sqlite", "db.sqlite3")
	if err != nil {
		log.Fatalf("Error connecting to SQLite database: %s.", err)
	}
	if err := db.PingContext(context.Background()); err != nil { // Though we open the database earlier, the connection isn't actually made until the first query so this is necessary
		log.Fatalf("Error connecting to SQLite database: %s.", err)
	}

	defer db.Close()
	log.Print("SQLite database connected successfully.")

	// Since in the case where the file doesn't exist it creates a new one, we must also guarentee that all of our tables are created within the DB. Also helps to spell them out explicitly for future reference.
	db.ExecContext(context.Background(),
		"CREATE TABLE IF NOT EXISTS BlogPosts (title TEXT PRIMARY KEY, url TEXT, subtitle TEXT, date INTEGER, readingmins INTEGER, body TEXT)",
	)
	db.ExecContext(context.Background(),
		"CREATE TABLE IF NOT EXISTS DailyLeetcode (title TEXT PRIMARY KEY, url TEXT, subtitle TEXT, date INTEGER, readingmins INTEGER, body TEXT, solution TEXT)",
	)

	// Next, we lay out the fuction responses all of our possible backend calls
	basic_response := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test Response.")
	}

	// HTML PAGE RESPONSES
	basic_html_response := func(w http.ResponseWriter, r *http.Request) {
		url_path := r.URL.EscapedPath()
		switch url_path {
		case "/":
			http.ServeFile(w, r, "../frontend/index.html")
		case "/about/":
			http.ServeFile(w, r, "../frontend/about.html")
		case "/blog/":
			http.ServeFile(w, r, "../frontend/blog.html")
		case "/dailyleetcode/":
			http.ServeFile(w, r, "../frontend/dailyleetcode.html")
		default:
			http.Error(w, "404 page not found", 404)
		}
	}

	blog_html_response := func(w http.ResponseWriter, r *http.Request) {
		url_path := r.URL.EscapedPath()
		article_url_name := strings.Split(url_path, "/")[2]

		post := BlogPost{}
		var temp_date_hold int64  // Stored in UNIX time in db, we have to convert it
		var temp_body_hold string // Similar, but stored using \n to break paragraphs so we need to split it up into the arr
		err := db.QueryRowContext(r.Context(), "SELECT title, subtitle, date, readingmins, body FROM BlogPosts WHERE url = $1", article_url_name).Scan(&(post.Title), &(post.Subtitle), &(temp_date_hold), &(post.Readingmins), &(temp_body_hold))
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No post found", 404)
				return
			} else {
				log.Printf("Server encountered error %s in blog_html_response", err)
				http.Error(w, "Internal Server Error", 500)
				return
			}
		}
		post.Date = time.Unix(temp_date_hold, 0).Format("Monday Jan 02, 2006")
		post.Body = strings.Split(temp_body_hold, `\n`)

		template_content, err := os.ReadFile("../frontend/blog_page.html")
		if err != nil {
			log.Printf("Server encountered error %s in blog_html_response", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		response_template, err := template.New("blog_response").Parse(string(template_content))
		if err != nil {
			log.Printf("Server encountered error %s in blog_html_response", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = response_template.Execute(w, post)
		if err != nil {
			log.Printf("Server encountered error %s in blog_html_response", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		// http.ServeFile(w, r, "../frontend/index.html")
	}

	leetcode_html_response := func(w http.ResponseWriter, r *http.Request) {

	}

	//CSS GENERIC RESPONSE
	css_response := func(w http.ResponseWriter, r *http.Request) {
		css_file := r.PathValue("filename")
		if !regexp.MustCompile(`^[\w-_]+\.css$`).MatchString(css_file) {
			http.Error(w, "invalid css file", 400)
			return
		}
		http.ServeFile(w, r, "../frontend/css/"+css_file)
	}

	//ASSET GENERIC RESPONSE
	asset_response := func(w http.ResponseWriter, r *http.Request) {
		asset_file := r.PathValue("filename")
		if !regexp.MustCompile(`^[\w-_]+\.[\w]+$`).MatchString(asset_file) {
			http.Error(w, "invalid css file", 400)
			return
		}
		http.ServeFile(w, r, "../frontend/assets/"+asset_file)
	}

	// Next, we create a new Mux and new HTTP Server to configure specific values
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:           ":42309",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Finally, we assign all of the relative URLs to their respective handler functions, after which we start the main server loop function. If the server ever returns (presumably because of an error) we log the error and exit the program.
	// Backend API for site content / article comments / signboard / etc..
	mux.HandleFunc("GET /data/all/info/", basic_response)
	mux.HandleFunc("GET /data/blog/info/", basic_response)
	mux.HandleFunc("GET /data/blog/content/", basic_response)
	mux.HandleFunc("GET /data/leetcode/info/", basic_response)
	mux.HandleFunc("GET /data/leetcode/content/", basic_response)

	mux.HandleFunc("GET /css/{filename}", css_response)
	mux.HandleFunc("GET /assets/{filename}", asset_response)

	// Frontend pages, HTML response
	mux.HandleFunc("GET /{$}", basic_html_response) // homepage

	mux.HandleFunc("GET /about/{$}", basic_html_response) // about page

	mux.HandleFunc("GET /blog/{$}", basic_html_response) // blog
	mux.HandleFunc("GET /blog/{blog_post}/{$}", blog_html_response)

	mux.HandleFunc("GET /dailyleetcode/{$}", basic_html_response) // daily leetcode
	mux.HandleFunc("GET /dailyleetcode/{date}/{$}", leetcode_html_response)

	log.Fatal(server.ListenAndServe())
}
