package handlers

import (
	"embed"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
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

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		outputErr(w, err, http.StatusBadRequest)
	}
	length, err := strconv.Atoi(r.URL.Query().Get("length"))
	if err != nil {
		outputErr(w, err, http.StatusBadRequest)
	}

	saved, err := h.service.FetchAll(offset, length)
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

	ulid, errorMessage := getUlidFromPath(r.URL.Path)
	if errorMessage != "" {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

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

	ulid, errorMessage := getUlidFromPath(r.URL.Path)
	if errorMessage != "" {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

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

	ulid, errorMessage := getUlidFromPath(r.URL.Path)
	if errorMessage != "" {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

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
	filePath, errorMessage := getFilePathFromUrlPath(r.URL.Path)
	if errorMessage != "" {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

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

// ex) input: /isumid/path_name/12345abcde -> output: 12345abcde
func getUlidFromPath(path string) (string, errorMessage string) {
	parts := strings.Split(path, "/")
	if len(parts) > 4 {
		err := fmt.Sprintf("invalid URL: %s, should be %s/[ulid]", path, strings.Join(parts[:3], "/"))
		return "", err
	}
	if len(parts) < 4 {
		return "", fmt.Sprintf("invalid URL: %s", path)
	}
	if len(parts) == 4 && parts[3] != "" {
		return parts[3], ""
	}
	err := fmt.Sprintf("invalid URL: %s, should be %s/[ulid]", path, strings.Join(parts[:3], "/"))
	return "", err
}

// ex) input: /isumid/path_name/ab/cd/ef -> output: ab/cd/ef
func getFilePathFromUrlPath(path string) (string, errorMessage string) {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return "", fmt.Sprintf("invalid URL: %s", path)
	}
	if len(parts) >= 4 && parts[3] != "" {
		return strings.Join(parts[3:], "/"), ""
	}
	err := fmt.Sprintf("invalid URL: %s, should be %s/[file path]", path, strings.Join(parts[:3], "/"))
	return "", err
}
