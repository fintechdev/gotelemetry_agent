package oauth

import (
	"errors"
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"golang.org/x/oauth2"
	"log"
)

var entries map[string]config.OAuthConfigEntry

func Init(e map[string]config.OAuthConfigEntry) {
	entries = e

	aggregations.InitOAuthStorage()
}

func configForEntryWithName(name string) (*oauth2.Config, error) {
	entry, ok := entries[name]

	if !ok {
		return nil, fmt.Errorf("oAuth entry %s not found", name)
	}

	res := &oauth2.Config{
		ClientID:     entry.ClientID,
		ClientSecret: entry.ClientSecret,
		Scopes:       entry.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  entry.AuthorizationURL,
			TokenURL: entry.TokenURL,
		},
		RedirectURL: "https://qa-www.telemetryapp.com/oauth_response",
	}

	return res, nil
}

func tokenForEntryWithName(name string) (*oauth2.Token, error) {
	if _, ok := entries[name]; !ok {
		return nil, fmt.Errorf("oAuth entry %s not found", name)
	}

	return aggregations.ReadOAuthToken(name)
}

func writeTokenForEntryWithName(name string, token *oauth2.Token) error {
	if _, ok := entries[name]; !ok {
		return fmt.Errorf("oAuth entry %s not found", name)
	}

	return aggregations.WriteOAuthToken(name, token)
}

func RunCommand(cfg config.CLIConfigType, errorChannel chan error, completionChannel chan bool) {
	switch cfg.OAuthCommand {
	case config.OAuthCommands.None:
		// Do nothing
		break

	case config.OAuthCommands.Request:
		entry, err := configForEntryWithName(cfg.OAuthName)

		if err != nil {
			errorChannel <- err
			break
		}

		errorChannel <- gotelemetry.NewLogError("Please visit this URL and authorize the Agent:\n\n---\n%s\n---\n\n", entry.AuthCodeURL("s", oauth2.AccessTypeOffline, oauth2.ApprovalForce))
		errorChannel <- gotelemetry.NewLogError("When you are done, please run the agent with the oauth-exchange command to set the new token.\n\n")

	case config.OAuthCommands.Exchange:

		entry, err := configForEntryWithName(cfg.OAuthName)

		if err != nil {
			errorChannel <- err
			break
		}

		code := cfg.OAuthCode

		if code == "" {
			errorChannel <- errors.New("No authorization code found. Please provide one with -c.")
			break
		}

		token, err := entry.Exchange(oauth2.NoContext, code)

		if err != nil {
			errorChannel <- err
			break
		}

		if err := writeTokenForEntryWithName(cfg.OAuthName, token); err != nil {
			errorChannel <- err
			break
		}

		errorChannel <- gotelemetry.NewLogError("Token exchanged successfully. The entry %s can now be used to make authenticated calls.", cfg.OAuthName)

	default:
		errorChannel <- fmt.Errorf("Unknown oauth command %s", cfg.OAuthCommand)
	}

	completionChannel <- true
}
