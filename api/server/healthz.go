package server

import "net/http"

func ready(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

func liveness(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
