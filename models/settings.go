package models

type Setting struct {
	OutputDir string
}

type RecordedResponse struct {
	ResBody   []byte
	ResHeader string
	ReqBody   []byte
	ReqHeader string
}
