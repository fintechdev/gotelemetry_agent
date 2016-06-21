package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/database"
	"github.com/telemetryapp/gotelemetry_agent/agent/graphite"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
	"github.com/telemetryapp/gotelemetry_agent/agent/oauth"
	"github.com/telemetryapp/gotelemetry_agent/agent/routes"
)

// VERSION number automatically populated by goxc config file
var VERSION = "3.0.2"

var configFile *config.File
var errorChannel chan error
var completionChannel chan bool

func handleErrors(errorChannel chan error) {

	for {
		select {
		case err, ok := <-errorChannel:
			if !ok {
				return
			}

			if e, ok := err.(*gotelemetry.Error); ok {
				logLevel := e.GetLogLevel()

				if logLevel >= config.CLIConfig.LogLevel {
					prefix := "Error"

					switch logLevel {
					case gotelemetry.LogLevelLog:
						prefix = "Info "

					case gotelemetry.LogLevelDebug:
						prefix = "Debug"
					}

					log.Printf("%s: %s", prefix, err)
				}

				continue
			}

			log.Printf("Error: %s", err.Error())
		}
	}
}

func main() {
	var err error

	config.Init(VERSION)

	configFile, err = config.NewConfigFile()

	if err != nil {
		log.Fatalf("Initialization error: %s", err)
	}

	errorChannel = make(chan error, 0)
	completionChannel = make(chan bool, 1)

	go handleErrors(errorChannel)
	go run()

	for {
		select {
		case <-completionChannel:
			goto Done
		}
	}

Done:

	for len(errorChannel) > 0 {
	}

	// Give error channel a moment to complete jobs in progress
	time.Sleep(100 * time.Millisecond)

	log.Println("No more jobs to run; exiting.")
}

func run() {
	if err := database.Init(configFile, errorChannel); err != nil {
		errorChannel <- gotelemetry.NewLogError("Initialization error: %s", err)
		completionChannel <- true
		return
	}

	if err := database.MergeDatabaseWithConfigFile(configFile); err != nil {
		errorChannel <- gotelemetry.NewLogError("Initialization error: %s", err)
		completionChannel <- true
		return
	}

	if err := graphite.Init(configFile, errorChannel); err != nil {
		errorChannel <- gotelemetry.NewLogError("Initialization error: %s", err)
		completionChannel <- true
		return
	}

	oauth.Init(configFile.OAuthConfig())

	if config.CLIConfig.IsPiping {
		payload, err := ioutil.ReadAll(os.Stdin)

		if err != nil {
			errorChannel <- err
		}

		agent.ProcessPipeRequest(configFile, errorChannel, completionChannel, payload)
	} else if config.CLIConfig.IsNotifying {
		agent.ProcessNotificationRequest(configFile, errorChannel, completionChannel, config.CLIConfig.NotificationChannel, config.CLIConfig.NotificationFlow, config.CLIConfig.Notification)
	} else if config.CLIConfig.OAuthCommand != config.OAuthCommands.None {
		oauth.RunCommand(config.CLIConfig, errorChannel, completionChannel)
	} else {

		serverListening, err := routes.Init(configFile, errorChannel)
		if err != nil {
			errorChannel <- gotelemetry.NewLogError("Initialization error: %s", err)
			completionChannel <- true
			return
		}

		if serverListening {
			if err := job.Init(configFile, errorChannel, nil); err != nil {
				errorChannel <- gotelemetry.NewLogError("Initialization error: %s", err)
				completionChannel <- true
				return
			}
			routes.SetAdditionalRoutes()
		} else {
			if err := job.Init(configFile, errorChannel, completionChannel); err != nil {
				errorChannel <- gotelemetry.NewLogError("Initialization error: %s", err)
				completionChannel <- true
				return
			}
		}

	}
}
