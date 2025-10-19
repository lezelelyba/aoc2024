package middleware

type OAuthReplyGithub struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type Token interface {
	Token() string
}

func (rep OAuthReplyGithub) Token() string {
	return rep.AccessToken
}
