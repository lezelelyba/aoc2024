package weberrors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type APIError struct {
	ErrorCode    int    `json:"errorcode"`
	ErrorMessage string `json:"errormessage"`
}

func NewError(status int, message string) APIError {
	return APIError{ErrorCode: status, ErrorMessage: message}
}

func HandleError(w http.ResponseWriter, logger *log.Logger, err error, httpErrorCode int, errMsg string) error {
	if err != nil {
		rc := httpErrorCode

		errJson, _ := json.Marshal(NewError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)
	}

	return err
}

func OkToError(ok bool) error {
	if ok {
		return nil
	}

	return fmt.Errorf("false")
}
