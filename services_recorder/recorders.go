package services_recorder

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kajikentaro/request-record-middleware/storages"
)

type Service struct {
	storage storages.Storage
}

func New(storage storages.Storage) Service {
	return Service{storage: storage}
}

func (r Service) Middleware(reqHeader http.Header, reqBody io.Reader, resHeader http.Header, resBody *[]byte, statusCode int) {
	ReqBodyData, err := io.ReadAll(reqBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	saveData := storages.SaveDataInput{
		ReqBody:    ReqBodyData,
		ResBody:    *resBody,
		ReqHeader:  reqHeader,
		ResHeader:  resHeader,
		StatusCode: statusCode,
	}
	err = r.storage.Save(saveData)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
