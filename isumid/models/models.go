package models

import (
	"net/http"
	"time"
)

type FetchTotalTransactionsResponse struct {
	Count int `json:"count"`
}

type RecordedDataInput struct {
	Method     string
	Url        string
	Path       string
	ReqHeader  map[string][]string
	ReqBody    []byte
	StatusCode int
	ResHeader  map[string][]string
	ResBody    []byte
	// StartedAt is when the request entered the middleware.
	StartedAt time.Time
	// Duration is the time spent in the wrapped handler (next.ServeHTTP).
	Duration time.Duration
}

type Meta struct {
	Method      string              `json:"method"`
	Url         string              `json:"url"`
	ReqHeader   map[string][]string `json:"reqHeader"`
	StatusCode  int                 `json:"statusCode"`
	ResHeader   map[string][]string `json:"resHeader"`
	IsReqText   bool                `json:"isReqText"`
	IsResText   bool                `json:"isResText"`
	ReqLength   int                 `json:"reqLength"`
	ResLength   int                 `json:"resLength"`
	Ulid        string              `json:"ulid"`
	Path        string              `json:"path"`
	StartedAtUs int64               `json:"startedAtUs"`
	DurationUs  int64               `json:"durationUs"`
}

type RecordedTransaction struct {
	Meta
	ReqBody string `json:"reqBody"`
	ResBody string `json:"resBody"`
}

type ReproducerResponse struct {
	IsSameResBody    bool                `json:"isSameResBody"`
	IsSameResHeader  bool                `json:"isSameResHeader"`
	IsSameStatusCode bool                `json:"isSameStatusCode"`
	ActualResHeader  map[string][]string `json:"actualResHeader"`
	ActualResBody    string              `json:"actualResBody"`
	IsBodyText       bool                `json:"isBodyText"`
	StatusCode       int                 `json:"statusCode"`
	ActualResLength  int                 `json:"actualResLength"`
}

type IsRecordingResponse struct {
	IsRecording bool `json:"isRecording"`
}

type FetchBodyResponse struct {
	Body   []byte
	Header http.Header
}
