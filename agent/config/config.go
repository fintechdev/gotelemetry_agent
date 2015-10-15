package config

import (
	"github.com/alecthomas/kingpin"
	"github.com/telemetryapp/gotelemetry"
	"log"
	"os"
	"regexp"
)

type CLIConfigType struct {
	APIURL              string
	ConfigFileLocation  string
	LogLevel            gotelemetry.LogLevel
	Filter              *regexp.Regexp
	ForceRunOnce        bool
	IsPiping            bool
	UseJSONPatch        bool
	UsePOST             bool
	IsNotifying         bool
	NotificationChannel string
	NotificationFlow    string
	Notification        gotelemetry.Notification
	WantsFunctionHelp   bool
	FunctionHelpName    string
}

const AgentVersion = "2.2.5"

var CLIConfig CLIConfigType

func banner() {
	println()
	println("Telemetry Agent v " + AgentVersion)
	println("Copyright © 2012-2015 Telemetry, Inc.")
	println()
	println("For license information, see the LICENSE file")
	println("---------------------------------------------")
	println()
}

func Init() {
	gotelemetry.UserAgentString = "Telemetry Agent v " + AgentVersion
	banner()

	app := kingpin.New("telemetry_agent", "The Telemetry Agent")

	app.Version(AgentVersion)

	app.Flag("config", "Path to the configuration file for this agent.").Short('c').Default("./gotelemetry_agent.yaml").StringVar(&CLIConfig.ConfigFileLocation)

	app.Flag("apiurl", "Set the URL to the Telemetry API").Short('a').Default("https://api.telemetryapp.com").StringVar(&CLIConfig.APIURL)
	logLevel := app.Flag("verbosity", "Set the verbosity level (`debug`, `info`, `error`).").Short('v').Default("info").Enum("debug", "info", "error")
	filter := app.Flag("filter", "Run only the jobs whose IDs (or tags if no ID is specified) match the given regular expression").Default(".").String()

	once := app.Command("once", "Run all jobs exactly once and exit.")

	pipe := app.Command("pipe", "Accept a Rails-style HTTP PATCH Telemetry payload from stdin, send it to the API, and then exit.")
	pipe.Flag("jsonpatch", "With --pipe, submit the package as a JSON-Patch request instead. Ignored otherwise.").BoolVar(&CLIConfig.UseJSONPatch)
	pipe.Flag("post", "With --pipe, submit the package as a POST request instead. Ignored otherwise.").BoolVar(&CLIConfig.UsePOST)

	notify := app.Command("notify", "Send a channel notification.")
	notify.Flag("channel", "The ID of the channel to which the notification is sent. Either channel or flow is required.").StringVar(&CLIConfig.NotificationChannel)
	notify.Flag("flow", "The Tag of the Flow to whose channel the notification is sent. Either channel or flow is required.").StringVar(&CLIConfig.NotificationFlow)
	notify.Flag("title", "The title of the notification.").Required().StringVar(&CLIConfig.Notification.Title)
	notify.Flag("message", "The message to be displayed in the notification.").Required().StringVar(&CLIConfig.Notification.Message)
	notify.Flag("icon", "An icon to be displayed in the notification.").StringVar(&CLIConfig.Notification.Icon)
	notify.Flag("duration", "The amount of seconds for which the notification must be displayed.").Default("1").IntVar(&CLIConfig.Notification.Duration)
	notify.Flag("sound", "A URL to a notification sound (use `default` for Telemetry's default notification sound).").StringVar(&CLIConfig.Notification.SoundURL)

	run := app.Command("run", "Runs the jobs scheduled in the configuration file provided.")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case once.FullCommand():
		CLIConfig.ForceRunOnce = true

	case pipe.FullCommand():
		CLIConfig.IsPiping = true

	case notify.FullCommand():
		CLIConfig.IsNotifying = true

	case run.FullCommand():
	default:
		// Do nothing, runs normally
	}

	switch *logLevel {
	case "debug":
		CLIConfig.LogLevel = gotelemetry.LogLevelDebug

	case "info":
		CLIConfig.LogLevel = gotelemetry.LogLevelLog

	case "error":
		CLIConfig.LogLevel = gotelemetry.LogLevelError

	default:
		log.Fatalf("Invalid verbosity level `%s`", *logLevel)
	}

	rx, err := regexp.Compile(*filter)

	if err != nil {
		log.Fatalf("Invalid regular expression provided for -filter: %s", err)
	}

	CLIConfig.Filter = rx
}
