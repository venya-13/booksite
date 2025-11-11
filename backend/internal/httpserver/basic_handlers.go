package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Google Auth Demo!")
	fmt.Fprintln(w, "Click here to <a href='/login'>Login with Google</a>")
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"status": "ok",
		"uptime": time.Since(startTime).String(),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

var startTime = time.Now()
