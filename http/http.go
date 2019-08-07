package http

import "net/http"

func Healthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("X-Custom-Header", "Awesome")
}
