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
		var sniffer bytes.Buffer
		newBody := ReadCloser{
			Reader:        io.TeeReader(r.Body, &sniffer),
			originalClose: r.Body.Close,
		}
		r.Body = newBody

		// prepare to read response body
		statusCode := 200
		newW := ResponseWriter{original: w, writtenData: &[]byte{}, statusCode: &statusCode}

		// go to original handler
		next.ServeHTTP(newW, r)

		reqBody, err := io.ReadAll(&sniffer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read ReqBodyData")
			return
		}

		saveData := storages.RecordedDataInput{
			Method:     r.Method,
			Url:        r.URL.String(),
			ReqHeader:  r.Header,
			ReqBody:    reqBody,
			StatusCode: statusCode,
			ResHeader:  newW.Header(),
			ResBody:    *newW.writtenData,
		}
		err = s.storage.Save(saveData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to save recorded data")
			return
		}
	})
}

func (s Middleware) Reproducer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get ulid from path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 || parts[3] == "" {
			http.Error(w, "invalid URL: should be /isumid/reproduce/[ulid]", http.StatusBadRequest)
			return
		}
		ulid := parts[3]

		// fetch recorded data
		savedRequestBody, err := s.storage.FetchReqBody(ulid)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read saved request body: %#v", err), http.StatusInternalServerError)
			return
		}

		savedMeta, err := s.storage.FetchMeta(ulid)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read saved meta data: %#v", err), http.StatusInternalServerError)
			return
		}

		// prepqre request
		newRequest, err := http.NewRequest(savedMeta.Method, savedMeta.Url, bytes.NewReader(savedRequestBody))
		newRequest.Header = savedMeta.ReqHeader
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to prepare request: %#v", err), http.StatusBadRequest)
			return
		}

		// prepare mocked ResponseWriter
		statusCode := 200
		newResponse := responseWriterOwn{header: &http.Header{}, writtenData: &bytes.Buffer{}, statusCode: &statusCode}

		// go to original handler
		next.ServeHTTP(newResponse, newRequest)

		savedResponseBody, err := s.storage.FetchResBody(ulid)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read saved request body: %#v", err), http.StatusInternalServerError)
			return
		}

		// DeepEqual fail unless convert map[string][]string
		var actualHeader map[string][]string = newResponse.Header()
		actualResBody := newResponse.writtenData.Bytes()
		s.storage.SaveReproduced(ulid, actualResBody, actualHeader)

		res := struct {
			IsSameResBody    bool
			IsSameResHeader  bool
			IsSameStatusCode bool
			ActualResHeader  map[string][]string
			ActualResBody    string
			IsBodyText       bool
			StatusCode       int
		}{
			IsSameResBody:    bytes.Equal(actualResBody, savedResponseBody),
			IsSameResHeader:  reflect.DeepEqual(actualHeader, savedMeta.ResHeader),
			IsSameStatusCode: statusCode == savedMeta.StatusCode,
			IsBodyText:       storages.IsText(newResponse.Header()),
			ActualResHeader:  actualHeader,
			StatusCode:       statusCode,
		}

		if res.IsBodyText {
			res.ActualResBody = string(actualResBody)
		}

		json, err := json.Marshal(res)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to stringify json: %#v", err), http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})
}
