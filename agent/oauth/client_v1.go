package oauth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

type v1ClientData struct {
	Client               *oauth.Client      `json:"-"`
	TemporaryCredentials *oauth.Credentials `json:"temporary_credentials"`
	PermanentCredentials *oauth.Credentials `json:"permanent_credentials"`
	Realm                string             `json:"realm"`
}

type v1Client struct {
	name string
	data v1ClientData
}

var _ Client = &v1Client{}

func getV1Client(name string, entry config.OAuthConfigEntry) (Client, error) {
	credentials := oauth.Credentials{
		Token:  entry.ClientID,
		Secret: entry.ClientSecret,
	}

	header := http.Header{}

	for key, value := range entry.Header {
		header.Set(key, value)
	}

	var signatureMethod oauth.SignatureMethod
	var privateKey *rsa.PrivateKey

	switch entry.SignatureMethod {
	case "":
		// Do nothing

	case "hmac_sha1":
		signatureMethod = oauth.HMACSHA1

	case "rsa_sha1":
		return nil, errors.New("RSA oAuth1 signatures are not yet supported")

	case "plaintext":
		signatureMethod = oauth.PLAINTEXT

	default:
		return nil, fmt.Errorf("Invalid oAuth1 signature format %s", entry.SignatureMethod)
	}

	client := &oauth.Client{
		Credentials:                   credentials,
		TemporaryCredentialRequestURI: entry.CredentialsURL,
		ResourceOwnerAuthorizationURI: entry.AuthorizationURL,
		TokenRequestURI:               entry.TokenURL,
		Header:                        header,
		SignatureMethod:               signatureMethod,
		PrivateKey:                    privateKey,
	}

	res := &v1Client{
		name: name,
		data: v1ClientData{
			Client: client,
		},
	}

	aggregations.ReadOAuthToken(res.name, &res.data)

	return res, nil
}

func (v *v1Client) getAuthorizationURL() (string, error) {
	cred, err := v.data.Client.RequestTemporaryCredentials(nil, telemetryOAuthClientResponseURL, nil)

	if err != nil {
		return "", err
	}

	v.data.TemporaryCredentials = cred

	err = aggregations.WriteOAuthToken(v.name, v.data)

	return v.data.Client.AuthorizationURL(cred, nil), err
}

func (v *v1Client) ExchangeToken(code, verifier, realm string) error {
	cred, _, err := v.data.Client.RequestToken(nil, v.data.TemporaryCredentials, verifier)

	if err != nil {
		return err
	}

	v.data.PermanentCredentials = cred
	v.data.Realm = realm

	return aggregations.WriteOAuthToken(v.name, v.data)
}

func (v *v1Client) Do(req *http.Request) (*http.Response, error) {
	if req.Method != "GET" {
		return nil, errors.New("Only GET transactions are supported for oAuth1 clients.")
	}

	return v.data.Client.Get(nil, v.data.PermanentCredentials, req.URL.String(), req.Form)
}
