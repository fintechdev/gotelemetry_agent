package oauth

import (
	"golang.org/x/oauth2"
	"net/http"
)

func Do(entryName string, req *http.Request) (*http.Response, error) {
	cfg, err := configForEntryWithName(entryName)

	if err != nil {
		return nil, err
	}

	token, err := tokenForEntryWithName(entryName)

	if err != nil {
		return nil, err
	}

	client := cfg.Client(oauth2.NoContext, token)

	res, err := client.Do(req)

	if err == nil {
		writeTokenForEntryWithName(entryName, token)
	}

	return res, err
}
