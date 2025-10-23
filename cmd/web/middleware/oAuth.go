// Package handles authentication with OAuth providers
package middleware

import (
	"encoding/json"
	"errors"
)

// Success reply from Github
type OAuthGithubOK struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// Error reply from Github
type OAuthGithubError struct {
	Err            string `json:"error"`
	ErrDescription string `json:"error_description"`
	ErrURI         string `json:"error_uri"`
}

// Combined reply from Github
type OAuthGithubReply struct {
	Response *OAuthGithubOK
	Err      *OAuthGithubError
}

// Interface for OAuth Response
type OAuthResponse interface {
	error
	Token() (string, error)
	IsError() bool
}

// Returns token
func (r OAuthGithubReply) Token() (string, error) {

	// if it' succesfull response
	if r.Response != nil {
		// provide token
		return r.Response.AccessToken, nil
	}

	return "", errors.New("reply is an error: unable to get token")
}

// Response is Error
func (r OAuthGithubReply) IsError() bool {
	return !(r.Err == nil)
}

// Error interface implementation for Reply
func (r OAuthGithubReply) Error() string {
	if r.Err != nil {
		return r.Err.Error()
	}

	return ""
}

// Error interface implementation for Error component of the Reply
func (e OAuthGithubError) Error() string {
	return e.Err
}

// Custom Unmarshaller for GitHub Response
func (r *OAuthGithubReply) UnmarshalJSON(data []byte) error {

	// unmarshal as a map
	var probe map[string]json.RawMessage

	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	// check if response is error, based on the error field
	if _, ok := probe["error"]; ok {
		var errObj OAuthGithubError

		// unmarshall as Error
		if err := json.Unmarshal(data, &errObj); err != nil {
			return err
		}

		r.Err = &errObj
		return nil
	}

	// otherwise unmarshal as a valid response
	var resObj OAuthGithubOK

	if err := json.Unmarshal(data, &resObj); err != nil {
		return err
	}

	r.Response = &resObj

	return nil
}
