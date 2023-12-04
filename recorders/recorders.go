package recorders

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kajikentaro/request-record-middleware/storages"
)

type Recorder struct {
	storage storages.Storage
}

func New(storage storages.Storage) Recorder {
	return Recorder{storage: storage}
}

func (r Recorder) Middleware(reqHeader http.Header, reqBody io.Reader, resHeader http.Header, resBody *[]byte, statusCode int) {
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
