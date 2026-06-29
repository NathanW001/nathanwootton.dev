package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	basic_response := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test Response.\r\n"))
	}

	blog_path_response := func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("blog_post")
		fmt.Fprintf(w, "Test Response. Recieved wild card with value \"%s\"", path)
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
	mux.HandleFunc("GET /data/", basic_response)

	// Frontend pages, HTML response
	mux.HandleFunc("GET /{$}", basic_response) // homepage

	mux.HandleFunc("GET /about/{$}", basic_response) // about page

	mux.HandleFunc("GET /blog/{$}", basic_response) // blog
	mux.HandleFunc("GET /blog/{blog_post}", blog_path_response)

	mux.HandleFunc("GET /dailyleetcode/{$}", basic_response) // daily leetcode
	mux.HandleFunc("GET /dailyleetcode/{date}", basic_response)

	log.Fatal(server.ListenAndServe())
}
