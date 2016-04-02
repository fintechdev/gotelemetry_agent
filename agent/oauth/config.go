package oauth

import (
	"fmt"

	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

const telemetryOAuthClientResponseURL = "https://telemetrytv.com/oauth_response"

var entries map[string]config.OAuthConfigEntry

// Init sets the entries global to the config file parameter
func Init(e map[string]config.OAuthConfigEntry) {
	entries = e
}

func clientForEntryWithName(name string) (Client, error) {
	entry, ok := entries[name]

	if !ok {
		return nil, fmt.Errorf("oAuth entry %s not found", name)
	}

	switch entry.Version {
	case 1:
		return getV1Client(name, entry)

	case 2:
		return getV2Client(name, entry)

	default:
		return nil, fmt.Errorf("Unknown oAuth version %d", entry.Version)
	}
}

// RunCommand is used to execute commands set using the OAuth command line flag
func RunCommand(cfg config.CLIConfigType, errorChannel chan error, completionChannel chan bool) {
	client, err := clientForEntryWithName(cfg.OAuthName)

	if err != nil {
		errorChannel <- err
		return
	}

	switch cfg.OAuthCommand {
	case config.OAuthCommands.None:
		// Do nothing
		break

	case config.OAuthCommands.Request:
		url, err := client.getAuthorizationURL()

		if err != nil {
			errorChannel <- err
			break
		}

		errorChannel <- gotelemetry.NewLogError("Please visit this URL and authorize the Agent:\n\n---\n%s\n---\n\n", url)
		errorChannel <- gotelemetry.NewLogError("When you are done, please run the agent with the oauth-exchange command to set the new token.\n\n")

	case config.OAuthCommands.Exchange:
		err := client.ExchangeToken(cfg.OAuthCode, cfg.OAuthVerifier, cfg.OAuthRealmID)

		if err != nil {
			errorChannel <- err
			break
		}

		errorChannel <- gotelemetry.NewLogError("Token exchanged successfully. The entry %s can now be used to make authenticated calls.", cfg.OAuthName)

	default:
		errorChannel <- fmt.Errorf("Unknown oauth command %s", cfg.OAuthCommand)
	}

	completionChannel <- true
}
