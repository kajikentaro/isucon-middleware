package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kajikentaro/isucon-middleware/isumid/services"
)

type Handler struct {
	service services.Service
}

func New(service services.Service) Handler {
	return Handler{service: service}
}

func outputErr(w http.ResponseWriter, err error, statusCode int) {
	message := fmt.Sprintf("request-record-middleware: %#v", err)
	http.Error(w, message, statusCode)
}

func (h Handler) FetchAll(w http.ResponseWriter, r *http.Request) {
	saved, err := h.service.FetchAll()
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(saved)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
	}
	w.Write(res)
}

func (h Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL: should be /fetch/[ulid]", http.StatusBadRequest)
		return
	}
	ulid := parts[2]

	saved, err := h.service.Fetch(ulid)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(saved)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
	}
	w.Write(res)
}
