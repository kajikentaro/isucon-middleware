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

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			outputErr(w, errors.New("query parameter 'offset' must be an integer"), http.StatusBadRequest)
			return
		}
	}

	length := 100
	if lengthStr := r.URL.Query().Get("length"); lengthStr != "" {
		var err error
		length, err = strconv.Atoi(lengthStr)
		if err != nil {
			outputErr(w, errors.New("query parameter 'length' must be an integer"), http.StatusBadRequest)
			return
		}
	}

	query := (r.URL.Query().Get("query"))

	searchResponse, err := h.service.Search(query, offset, length)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(searchResponse)
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

func (h Handler) Remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ulid, errorMessage := getUlidFromPath(r.URL.Path)
	if errorMessage != "" {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	err := h.service.Remove(ulid)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) RemoveAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := h.service.RemoveAll()
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) FetchTotalTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	totalTransactions, err := h.service.FetchTotalTransactions()
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(totalTransactions)
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
