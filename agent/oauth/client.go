package oauth

import (
	"net/http"
)

// Client is an interface to allow different OAuth versions to be used the same way
type Client interface {
	getAuthorizationURL() (string, error)
	ExchangeToken(code, verifier, realm string) error
	Do(req *http.Request) (*http.Response, error)
}

// Do is used to execute OAuth 1.0a and 2.0 requests using the same interface
func Do(entryName string, req *http.Request) (*http.Response, error) {
	client, err := clientForEntryWithName(entryName)

	if err != nil {
		return nil, err
	}

	return client.Do(req)
}
