package main

import (
	"net/http"
)

func main() {
	basic_response := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test Response.\r\n"))
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", basic_response)

	http.ListenAndServe(":80", mux)
}
