package handlers

import (
	"embed"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

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

func (h Handler) FetchReqBody(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL: should be /req-body/[ulid]", http.StatusBadRequest)
		return
	}
	ulid := parts[2]

	saved, err := h.service.FetchReqBody(ulid)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Println(w.Header())

	for key, values := range saved.Header {
		w.Header()[key] = values
	}
	w.Write(saved.Body)
}

func (h Handler) FetchResBody(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL: should be /res-body/[ulid]", http.StatusBadRequest)
		return
	}
	ulid := parts[2]

	saved, err := h.service.FetchResBody(ulid)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}

	for key, values := range saved.Header {
		w.Header()[key] = values
	}
	w.Write(saved.Body)
}

func (h Handler) FetchReproducedResBody(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL: should be /reproduced-res-body/[ulid]", http.StatusBadRequest)
		return
	}
	ulid := parts[2]

	saved, err := h.service.FetchReproducedResBody(ulid)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}

	for key, values := range saved.Header {
		w.Header()[key] = values
	}
	w.Write(saved.Body)
}

//go:embed front-built/*
var assets embed.FS

func (h Handler) Frontend(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL: should be /isumid/[filepath]", http.StatusBadRequest)
		return
	}
	filePath := strings.Join(parts[2:], "/")

	data, err := assets.ReadFile("front-built/" + filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	extension := filepath.Ext(filePath)
	if mimeType := mime.TypeByExtension(extension); mimeType != "" {
		w.Header().Add("Content-Type", mimeType)
	}
	fmt.Fprintf(w, "%s", data)
}
