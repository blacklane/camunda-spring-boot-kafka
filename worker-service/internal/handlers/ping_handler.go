package handlers

import "net/http"

// PingHandler handle the /ping
func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		resp := "ok\n"
		w.Write([]byte(resp))
	}
}
