package handlers

import (
	"embed"
	"encoding/json"
	"errors"
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

func (h Handler) FetchList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Query().Get("offset") == "" {
		outputErr(w, errors.New("query parameter 'offset' is not defined"), http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		outputErr(w, err, http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("length") == "" {
		outputErr(w, errors.New("query parameter 'length' is not defined"), http.StatusBadRequest)
		return
	}
	length, err := strconv.Atoi(r.URL.Query().Get("length"))
	if err != nil {
		outputErr(w, err, http.StatusBadRequest)
		return
	}

	saved, err := h.service.FetchList(offset, length)
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
	const ANS_IDX = 4
	parts := strings.Split(path, "/")
	if len(parts) > ANS_IDX {
		err := fmt.Sprintf("invalid URL: %s, should be %s/[ulid]", path, strings.Join(parts[:ANS_IDX-1], "/"))
		return "", err
	}
	if len(parts) < ANS_IDX {
		return "", fmt.Sprintf("invalid URL: %s", path)
	}
	if len(parts) == ANS_IDX && parts[ANS_IDX-1] != "" {
		return parts[ANS_IDX-1], ""
	}
	err := fmt.Sprintf("invalid URL: %s, should be %s/[ulid]", path, strings.Join(parts[:ANS_IDX-1], "/"))
	return "", err
}

// ex) input: /isumid/path_name/ab/cd/ef -> output: ab/cd/ef
func getFilePathFromUrlPath(path string) (string, errorMessage string) {
	const ANS_IDX = 3
	parts := strings.Split(path, "/")
	if len(parts) < ANS_IDX {
		return "", fmt.Sprintf("invalid URL: %s", path)
	}
	if len(parts) >= ANS_IDX && parts[ANS_IDX-1] != "" {
		return strings.Join(parts[ANS_IDX-1:], "/"), ""
	}
	err := fmt.Sprintf("invalid URL: %s, should be %s/[file path]", path, strings.Join(parts[:ANS_IDX-1], "/"))
	return "", err
}
