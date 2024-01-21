package models

import "net/http"

type FetchTotalTransactionsResponse struct {
	Count int
}

type RecordedDataInput struct {
	Method    string
	Url       string
	ReqHeader map[string][]string
	ReqBody   []byte

	StatusCode int
	ResHeader  map[string][]string
	ResBody    []byte
}

type Meta struct {
	Method    string
	Url       string
	ReqHeader map[string][]string

	StatusCode int
	ResHeader  map[string][]string

	IsReqText bool
	IsResText bool
	ReqLength int
	ResLength int
	Ulid      string
}

type RecordedTransaction struct {
	Meta
	ReqBody string
	ResBody string
}

type ReproducerResponse struct {
	IsSameResBody    bool
	IsSameResHeader  bool
	IsSameStatusCode bool
	ActualResHeader  map[string][]string
	ActualResBody    string
	IsBodyText       bool
	StatusCode       int
	ActualResLength  int
}

type IsRecordingResponse struct {
	IsRecording bool
}

type FetchBodyResponse struct {
	Body   []byte
	Header http.Header
}
