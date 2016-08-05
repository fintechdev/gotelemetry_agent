package main

import (
	"container/list"
	"fmt"
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
var VERSION = "3.0.3"
var configFile *config.File
var errorChannel chan error
var completionChannel chan bool
var apiStreamChannel chan string
var streamRunning bool
var logList *list.List

func handleErrors() {

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

					fmtMessage := fmt.Sprintf("%s: %s", prefix, err)

					// Maintain a list of the 100 most recent log items
					for logList.Len() >= 100 {
						logList.Remove(logList.Front())
					}
					logList.PushBack(fmtMessage)

					if streamRunning {
						apiStreamChannel <- fmtMessage
					}

					log.Print(fmtMessage)
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
	apiStreamChannel = make(chan string, 2)
	logList = list.New()

	go handleErrors()
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

			routes.SetAdditionalRoutes(apiStreamChannel, &streamRunning, logList)

		} else {
			if err := job.Init(configFile, errorChannel, completionChannel); err != nil {
				errorChannel <- gotelemetry.NewLogError("Initialization error: %s", err)
				completionChannel <- true
				return
			}
		}

	}
}
