// Errors returned by the application
package weberrors

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Error sent by application for web responses, API responses
type AoCError struct {
	ErrorCode    int    `json:"errorcode"`
	ErrorMessage string `json:"errormessage"`
} //@name Error

// Constructor
func NewError(status int, message string) AoCError {
	return AoCError{ErrorCode: status, ErrorMessage: message}
}

// Handles errors - write response code and encapsulates err msg into json
func HandleError(w http.ResponseWriter, logger *log.Logger, err error, httpErrorCode int, errMsg string) error {
	if err != nil {
		rc := httpErrorCode

		errJson, _ := json.Marshal(NewError(rc, errMsg))

		logger.Println(errMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(rc)
		w.Write(errJson)
	}

	return err
}

// Helper function - converts bool to error
// true -> nil
// false -> generic error
func OkToError(ok bool) error {
	if ok {
		return nil
	}

	return errors.New("false")
}
