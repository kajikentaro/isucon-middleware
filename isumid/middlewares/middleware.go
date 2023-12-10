package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/kajikentaro/isucon-middleware/isumid/storages"
)

type Middleware struct {
	storage storages.Storage
}

func New(storage storages.Storage) Middleware {
	return Middleware{storage: storage}
}

func (s Middleware) Recorder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// prepare to read request body
		var reqBodyBuffer bytes.Buffer
		reader := io.TeeReader(r.Body, &reqBodyBuffer)
		sniffedReadCloser := readCloser{
			Reader:        reader,
			originalClose: r.Body.Close,
		}
		r.Body = sniffedReadCloser

		// prepare to read response body
		statusCode := 200
		sniffedResponseWriter := responseWriterSniffer{original: w, writtenData: &[]byte{}, statusCode: &statusCode}

		// go to original handler
		next.ServeHTTP(sniffedResponseWriter, r)

		ReqBodyData, err := io.ReadAll(&reqBodyBuffer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read ReqBodyData")
			return
		}

		reqOthers := storages.RequestOthers{
			Url:    r.URL.String(),
			Header: r.Header,
			Method: r.Method,
		}
		saveData := storages.RecordedDataInput{
			ReqBody:    ReqBodyData,
			ResBody:    *sniffedResponseWriter.writtenData,
			ReqOthers:  reqOthers,
			ResHeader:  w.Header(),
			StatusCode: statusCode,
		}
		err = s.storage.Save(saveData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to save recorded data")
			return
		}
	})
}

func (s Middleware) Executer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get ulid from path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 || parts[2] == "" {
			http.Error(w, "invalid URL: should be /execute/[ulid]", http.StatusBadRequest)
			return
		}
		ulid := parts[2]

		// fetch recorded data
		saved, err := s.storage.FetchDetail(ulid)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read save data: %#v", err), http.StatusBadRequest)
			return
		}

		// prepqre request
		request, err := http.NewRequest(saved.ReqOthers.Method, saved.ReqOthers.Url, bytes.NewReader(saved.ReqBody))
		request.Header = saved.ReqOthers.Header
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to prepare request: %#v", err), http.StatusBadRequest)
			return
		}

		// prepare to read response body
		statusCode := 200
		responseWriter := responseWriterOwn{header: &http.Header{}, writtenData: &bytes.Buffer{}, statusCode: &statusCode}

		// go to original handler
		next.ServeHTTP(responseWriter, request)

		actualResBody := responseWriter.writtenData.Bytes()
		isSameResBody := bytes.Equal(actualResBody, saved.ResBody)
		isSameResHeader := reflect.DeepEqual(responseWriter.Header(), saved.ResHeader)
		isSameStatusCode := statusCode == saved.StatusCode
		res := struct {
			IsSameResBody    bool
			IsSameResHeader  bool
			IsSameStatusCode bool
			ActualResHeader  http.Header
			ActualResBody    string
		}{
			IsSameResBody:    isSameResBody,
			IsSameResHeader:  isSameResHeader,
			IsSameStatusCode: isSameStatusCode,
			ActualResHeader:  responseWriter.Header(),
		}

		if storages.IsText(responseWriter.Header()) {
			res.ActualResBody = string(actualResBody)
		}

		json, err := json.Marshal(res)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to stringify json: %#v", err), http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(json)
	})
}
