package oauth

import (
	"errors"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"golang.org/x/oauth2"
	"net/http"
)

type v2Client struct {
	name string
	cfg  *oauth2.Config
	t    *oauth2.Token
}

var _ Client = &v2Client{}

func getV2Client(name string, entry config.OAuthConfigEntry) (Client, error) {
	endpoint := oauth2.Endpoint{
		AuthURL:  entry.AuthorizationURL,
		TokenURL: entry.TokenURL,
	}

	cfg := &oauth2.Config{
		ClientID:     entry.ClientID,
		ClientSecret: entry.ClientSecret,
		Scopes:       entry.Scopes,
		Endpoint:     endpoint,
		RedirectURL:  "https://telemetryapp.com/oauth_response",
	}

	res := &v2Client{
		name: name,
		cfg:  cfg,
		t:    &oauth2.Token{},
	}

	err := aggregations.ReadOAuthToken(res.name, &res.t)

	return res, err
}

func (v *v2Client) GetAuthorizationURL() (string, error) {
	return v.cfg.AuthCodeURL("s", oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (v *v2Client) ExchangeToken(code, verifier, realm string) error {
	if code == "" {
		return errors.New("No authorization code found. Please provide one with -c.")
	}

	token, err := v.cfg.Exchange(oauth2.NoContext, code)

	if err == nil {
		aggregations.WriteOAuthToken(v.name, token)
		v.t = token
	}

	return err
}

func (v *v2Client) Do(req *http.Request) (*http.Response, error) {
	client := v.cfg.Client(oauth2.NoContext, v.t)

	res, err := client.Do(req)

	if err == nil {
		return res, aggregations.WriteOAuthToken(v.name, v.t)
	}

	return res, err
}
