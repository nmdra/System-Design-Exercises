package main

import "net/http"

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", process())
	return mux
}
