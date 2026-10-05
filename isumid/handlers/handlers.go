package handlers

import (
	"compress/gzip"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kajikentaro/isucon-middleware/isumid/models"
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

	query := r.URL.Query().Get("query")

	searchResponse, err := h.service.Search(query, offset, length)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(searchResponse)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func (h Handler) FetchReqBody(w http.ResponseWriter, r *http.Request) {
	serveBody(w, r, h.service.FetchReqBody)
}

func (h Handler) FetchResBody(w http.ResponseWriter, r *http.Request) {
	serveBody(w, r, h.service.FetchResBody)
}

func (h Handler) FetchReproducedResBody(w http.ResponseWriter, r *http.Request) {
	serveBody(w, r, h.service.FetchReproducedResBody)
}

// serveBody writes a saved body with its saved header. The ulid is taken from the path.
func serveBody(w http.ResponseWriter, r *http.Request, fetch func(ulid string) (models.FetchBodyResponse, error)) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ulid, errorMessage := getUlidFromPath(r.URL.Path)
	if errorMessage != "" {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	saved, err := fetch(ulid)
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

const defaultExportLimitMB = 100

// parseLimitMB parses a query parameter of a size limit in MB, and returns it in bytes.
func parseLimitMB(r *http.Request, name string) (int64, error) {
	str := r.URL.Query().Get(name)
	if str == "" {
		return defaultExportLimitMB << 20, nil
	}
	mb, err := strconv.ParseInt(str, 10, 64)
	if err != nil || mb < 0 || mb > 1<<20 {
		return 0, fmt.Errorf("query parameter '%s' must be an integer between 0 and %d", name, 1<<20)
	}
	return mb << 20, nil
}

// Export returns a SQLite file of the recorded data with the bodies embedded, for analysis.
func (h Handler) Export(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	maxMetaBytes, err := parseLimitMB(r, "maxMetaMB")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	maxBodyBytes, err := parseLimitMB(r, "maxBodyMB")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := h.service.Export(maxMetaBytes, maxBodyBytes)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/vnd.sqlite3")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(file.Name())))
	w.Header().Set("Vary", "Accept-Encoding")

	if !acceptsGzip(r) {
		http.ServeContent(w, r, "", time.Now(), file)
		return
	}

	// compress while sending (the size is unknown in advance, so no Content-Length).
	// Browsers decompress it, so the saved file is a plain SQLite file.
	w.Header().Set("Content-Encoding", "gzip")
	gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(gz, file); err != nil {
		// the header is already sent, so the error cannot be returned to the client
		fmt.Fprintf(os.Stderr, "failed to send the exported file: %s\n", err)
		return
	}
	if err := gz.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to send the exported file: %s\n", err)
	}
}

// acceptsGzip reports whether the Accept-Encoding header of the request allows gzip.
func acceptsGzip(r *http.Request) bool {
	for _, value := range r.Header.Values("Accept-Encoding") {
		for _, coding := range strings.Split(value, ",") {
			name, params, _ := strings.Cut(coding, ";")
			if strings.TrimSpace(name) != "gzip" {
				continue
			}
			// "gzip;q=0" means gzip is not acceptable
			qStr, ok := strings.CutPrefix(strings.TrimSpace(params), "q=")
			if !ok {
				return true
			}
			q, err := strconv.ParseFloat(qStr, 64)
			return err == nil && q > 0
		}
	}
	return false
}

func (h Handler) ExportSize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	size, err := h.service.ExportSize()
	if err != nil {
		outputErr(w, err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(size)
}

//go:embed front-built/*
var assets embed.FS

// Frontend serves the built web UI. The UI refers to its assets and the APIs with relative paths,
// so it works under any prefix. Expects a request path with the prefix stripped.
func (h Handler) Frontend() http.Handler {
	frontBuilt, err := fs.Sub(assets, "front-built")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(frontBuilt))
}

// Expects a request path with the prefix stripped.
// ex) input: /path_name/12345abcde -> output: 12345abcde
func getUlidFromPath(path string) (string, errorMessage string) {
	const ANS_IDX = 3
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
