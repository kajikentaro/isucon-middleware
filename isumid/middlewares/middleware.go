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
	"time"

	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
)

type Middleware struct {
	storage     storages.Storage
	isRecording bool
	autoStop    *settings.AutoSwitch
	autoStart   *settings.AutoSwitch
}

func New(storage storages.Storage, options *settings.Setting) Middleware {
	isRecording := true
	var autoStop *settings.AutoSwitch = nil
	var autoStart *settings.AutoSwitch = nil
	if options != nil {
		isRecording = options.RecordOnStart
		autoStop = options.AutoStop
		autoStart = options.AutoStart
	}

	return Middleware{
		storage:     storage,
		isRecording: isRecording,
		autoStop:    autoStop,
		autoStart:   autoStart,
	}
}

func (s *Middleware) StopRecording(w http.ResponseWriter, r *http.Request) {
	s.isRecording = false
}

func (s *Middleware) StartRecording(w http.ResponseWriter, r *http.Request) {
	s.isRecording = true
}

func (s *Middleware) Recorder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.autoStart != nil && s.autoStart.TriggerEndpoint == r.URL.Path {
			go func() {
				time.Sleep(time.Second * time.Duration(s.autoStart.AfterSec))
				s.isRecording = true
			}()
		}
		if s.autoStop != nil && s.autoStop.TriggerEndpoint == r.URL.Path {
			go func() {
				time.Sleep(time.Second * time.Duration(s.autoStop.AfterSec))
				s.isRecording = false
			}()
		}
		if !s.isRecording {
			next.ServeHTTP(w, r)
			return
		}

		// prepare to read request body
		var reqBodySniffer bytes.Buffer
		newBody := ReadCloser{
			Reader:        io.TeeReader(r.Body, &reqBodySniffer),
			originalClose: r.Body.Close,
		}
		r.Body = newBody

		// prepare to read response body
		var resBodySniffer bytes.Buffer
		statusCode := 200
		newW := ResponseWriter{original: w, writtenData: &resBodySniffer, statusCode: &statusCode}

		// go to original handler
		next.ServeHTTP(newW, r)

		reqBody, err := io.ReadAll(&reqBodySniffer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read ReqBodyData")
			return
		}

		resBody, err := io.ReadAll(&resBodySniffer)
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
			ResBody:    resBody,
		}
		err = s.storage.Save(saveData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to save recorded data")
			return
		}
	})
}

func (s *Middleware) Reproducer(next http.Handler) http.Handler {
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
			ActualResLength  int
		}{
			IsSameResBody:    bytes.Equal(actualResBody, savedResponseBody),
			IsSameResHeader:  reflect.DeepEqual(actualHeader, savedMeta.ResHeader),
			IsSameStatusCode: statusCode == savedMeta.StatusCode,
			IsBodyText:       storages.IsText(newResponse.Header(), actualResBody),
			ActualResHeader:  actualHeader,
			StatusCode:       statusCode,
			ActualResLength:  len(actualResBody),
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
