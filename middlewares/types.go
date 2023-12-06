package middlewares

import (
	"bytes"
	"io"
	"net/http"
)

type readCloser struct {
	io.Reader
	originalClose func() error
}

func (n readCloser) Close() error {
	return n.originalClose()
}

type responseWriterSniffer struct {
	original    http.ResponseWriter
	writtenData *[]byte
	statusCode  *int
}

func (r responseWriterSniffer) Header() http.Header {
	return r.original.Header()
}
func (r responseWriterSniffer) Write(in []byte) (int, error) {
	*r.writtenData = in
	return r.original.Write(in)
}
func (r responseWriterSniffer) WriteHeader(statusCode int) {
	*r.statusCode = statusCode
	r.original.WriteHeader(statusCode)
}

type responseWriterOwn struct {
	header      *http.Header
	writtenData *bytes.Buffer
	statusCode  *int
}

func (r responseWriterOwn) Header() http.Header {
	return *r.header
}
func (r responseWriterOwn) Write(in []byte) (int, error) {
	return r.writtenData.Write(in)
}
func (r responseWriterOwn) WriteHeader(statusCode int) {
	*r.statusCode = statusCode
}
