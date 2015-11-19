package oauth

import (
	"net/http"
)

type Client interface {
	GetAuthorizationURL() (string, error)
	ExchangeToken(code, verifier, realm string) error
	Do(req *http.Request) (*http.Response, error)
}

func Do(entryName string, req *http.Request) (*http.Response, error) {
	client, err := clientForEntryWithName(entryName)

	if err != nil {
		return nil, err
	}

	return client.Do(req)
}
