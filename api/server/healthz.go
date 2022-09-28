package server

import "net/http"

func ready(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func liveness(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
