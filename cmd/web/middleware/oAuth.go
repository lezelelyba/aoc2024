package middleware

import (
	"encoding/json"
	"errors"
)

type OAuthGithubOK struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type OAuthGithubError struct {
	Err            string `json:"error"`
	ErrDescription string `json:"error_description"`
	ErrURI         string `json:"error_uri"`
}

type OAuthGithubReply struct {
	Response *OAuthGithubOK
	Err      *OAuthGithubError
}

type OAuthResponse interface {
	error
	Token() (string, error)
	IsError() bool
}

func (r OAuthGithubReply) Token() (string, error) {

	if r.Response != nil {
		return r.Response.AccessToken, nil
	}

	return "", errors.New("reply is an error: unable to get token")
}

func (r OAuthGithubReply) IsError() bool {
	return !(r.Err == nil)
}

func (r OAuthGithubReply) Error() string {
	if r.Err != nil {
		return r.Err.Error()
	}

	return ""
}

func (e OAuthGithubError) Error() string {
	return e.Err
}

func (r *OAuthGithubReply) UnmarshalJSON(data []byte) error {

	var probe map[string]json.RawMessage

	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	if _, ok := probe["error"]; ok {
		var errObj OAuthGithubError

		if err := json.Unmarshal(data, &errObj); err != nil {
			return err
		}

		r.Err = &errObj
		return nil
	}

	var resObj OAuthGithubOK

	if err := json.Unmarshal(data, &resObj); err != nil {
		return err
	}

	r.Response = &resObj

	return nil
}
