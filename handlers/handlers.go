package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kajikentaro/request-record-middleware/services_endpoint"
	"github.com/kajikentaro/request-record-middleware/services_recorder"
)

type Handler struct {
	service  services_endpoint.Service
	recorder services_recorder.Service
}

type readCloser struct {
	io.Reader
	originalClose func() error
}

func (n readCloser) Close() error {
	return n.originalClose()
}

func New(service services_endpoint.Service, recorder services_recorder.Service) Handler {
	return Handler{service: service, recorder: recorder}
}

func outputErr(w http.ResponseWriter, err error) {
	message := fmt.Sprintf("%#v", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

func (h Handler) FetchAll(w http.ResponseWriter, r *http.Request) {
	saved, err := h.service.FetchAll()
	if err != nil {
		outputErr(w, err)
		return
	}

	res, err := json.Marshal(saved)
	if err != nil {
		outputErr(w, err)
	}
	w.Write(res)
}

func (h Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	ulid := parts[2]

	saved, err := h.service.Fetch(ulid)
	if err != nil {
		outputErr(w, err)
		return
	}

	res, err := json.Marshal(saved)
	if err != nil {
		outputErr(w, err)
	}
	w.Write(res)
}

type responseWriter struct {
	original    http.ResponseWriter
	writtenData *[]byte
	statusCode  *int
}

func (r responseWriter) Header() http.Header {
	return r.original.Header()
}
func (r responseWriter) Write(in []byte) (int, error) {
	*r.writtenData = in
	return r.original.Write(in)
}
func (r responseWriter) WriteHeader(statusCode int) {
	*r.statusCode = statusCode
	r.original.WriteHeader(statusCode)
}

func (h Handler) RecorderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		reader := io.TeeReader(r.Body, &buf)
		sniffedReadCloser := readCloser{
			Reader:        reader,
			originalClose: r.Body.Close,
		}
		r.Body = sniffedReadCloser

		var statusCode int
		sniffedResponseWriter := responseWriter{original: w, writtenData: &[]byte{}, statusCode: &statusCode}

		next.ServeHTTP(sniffedResponseWriter, r)

		h.recorder.Middleware(r.Header, &buf, w.Header(), sniffedResponseWriter.writtenData, *sniffedResponseWriter.statusCode)
	})
}
