package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Hello - The handler struct for the homepage route "/"
type Hello struct {
	l *log.Logger
}

// NewHello(*log.Logger) - Creates a new HTTP handler of the type Hello and returns a pointer to the same.
// Expects a pointer to the logrus.Logger to be passed as the only parameter
func NewHello(l *log.Logger) *Hello {
	return &Hello{
		l: l,
	}
}

// Handler(http.ResponseWriter, *http.Request) - Creates a new HTTP Handler function for the Hello struct that handles
// the homepage route "/". Takes 2 arguments of the types http.ResponseWriter and a pointer to *http.Request, in that order
func (h *Hello) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := w.Write([]byte("Hello World!"))
	if err != nil {
		h.l.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
