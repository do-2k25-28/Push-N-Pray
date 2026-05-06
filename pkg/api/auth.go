package api

import (
	"encoding/base64"
	"net/http"
)

type Authenticator interface {
	Apply(req *http.Request)
}

type BearerToken string

func (t BearerToken) Apply(req *http.Request) {
	if t == "" {
		return
	}
	req.Header.Set("Authorization", "Bearer "+string(t))
}

type BasicAuth struct {
	Username string
	Token    string
}

func (b BasicAuth) Apply(req *http.Request) {
	if b.Username == "" || b.Token == "" {
		return
	}
	creds := b.Username + ":" + b.Token
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(creds)))
}
