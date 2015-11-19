package oauth

import (
	"fmt"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
)

var entries map[string]config.OAuthConfigEntry

func Init(e map[string]config.OAuthConfigEntry) {
	entries = e

	aggregations.InitOAuthStorage()
}

func clientForEntryWithName(name string) (Client, error) {
	entry, ok := entries[name]

	if !ok {
		return nil, fmt.Errorf("oAuth entry %s not found", name)
	}

	switch entry.Version {
	case 1:
		return nil, fmt.Errorf("Unknown oAuth version %d", entry.Version)

	case 2:
		return getV2Client(name, entry)

	default:
		return nil, fmt.Errorf("Unknown oAuth version %d", entry.Version)
	}
}

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
		url, err := client.GetAuthorizationURL()

		if err != nil {
			errorChannel <- err
			break
		}

		errorChannel <- gotelemetry.NewLogError("Please visit this URL and authorize the Agent:\n\n---\n%s\n---\n\n", url)
		errorChannel <- gotelemetry.NewLogError("When you are done, please run the agent with the oauth-exchange command to set the new token.\n\n")

	case config.OAuthCommands.Exchange:
		err := client.ExchangeToken(cfg.OAuthCode, "", "")

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
