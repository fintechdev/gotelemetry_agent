package oauth

import (
	"errors"
	"net/http"
	"time"

	"github.com/telemetryapp/gotelemetry_agent/agent/database"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"golang.org/x/oauth2"
)

type v2Client struct {
	name string
	cfg  *oauth2.Config
	t    *oauth2.Token
	ttl  string
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
		RedirectURL:  telemetryOAuthClientResponseURL,
	}

	res := &v2Client{
		name: name,
		cfg:  cfg,
		t:    &oauth2.Token{},
		ttl:  entry.TTL,
	}

	err := database.ReadOAuthToken(res.name, &res.t)

	return res, err
}

func (v *v2Client) getAuthorizationURL() (string, error) {
	return v.cfg.AuthCodeURL("s", oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (v *v2Client) ExchangeToken(code, verifier, realm string) error {
	if code == "" {
		return errors.New("No authorization code found. Please provide one with -o.")
	}

	token, err := v.cfg.Exchange(oauth2.NoContext, code)

	if err == nil {
		// Add an expiration time if the TTL property has been set
		if v.ttl != "" {
			var ttl time.Duration
			ttl, err = config.ParseTimeInterval(v.ttl)
			if err != nil {
				return err
			}
			token.Expiry = time.Now().Add(ttl)
		}

		database.WriteOAuthToken(v.name, token)
		v.t = token
	}

	return err
}

func (v *v2Client) Do(req *http.Request) (*http.Response, error) {
	var err error

	// Refresh the token
	v.t, err = v.cfg.TokenSource(oauth2.NoContext, v.t).Token()

	if err != nil {
		return nil, err
	}

	// Update the token expiration time if it has been set
	if v.ttl != "" {
		var ttl time.Duration
		ttl, err = config.ParseTimeInterval(v.ttl)
		if err != nil {
			return nil, err
		}
		v.t.Expiry = time.Now().Add(ttl)
	}

	// Write the updated token to the database
	err = database.WriteOAuthToken(v.name, v.t)

	if err != nil {
		return nil, err
	}

	// Execute the query
	client := v.cfg.Client(oauth2.NoContext, v.t)
	res, err := client.Do(req)

	return res, err
}
