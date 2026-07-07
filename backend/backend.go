package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

type BlogPost struct {
	Title    string
	Subtitle string
	Date     time.Time
	Readmins int
	Body     []string
}

func main() {
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
		// TODO: connect to sqlite and get information
		// TODO: add html/template with blog_page.html
		http.ServeFile(w, r, "../frontend/blog_page.html")
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

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:           ":42309",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

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
