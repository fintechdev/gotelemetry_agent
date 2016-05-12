package config

import (
	"log"
	"os"
	"regexp"

	"github.com/alecthomas/kingpin"
	"github.com/telemetryapp/gotelemetry"
)

// OAuthCommand is set to the type of command that will be executed by oauth.RunCommand
type OAuthCommand string

// OAuthCommands are the states that an OAuth command can be set to
var OAuthCommands = struct {
	None     OAuthCommand
	Request  OAuthCommand
	Exchange OAuthCommand
}{
	None:     "",
	Request:  "request",
	Exchange: "exchange",
}

// CLIConfigType manages the various settings that are initialized at Agent launch
type CLIConfigType struct {
	APIURL              string
	ChannelTag          string
	ConfigFileLocation  string
	DatabasePath        string
	DatabaseTTL         string
	AuthenticationKey   string
	AuthenticationPort  string
	CertFile            string
	KeyFile             string
	LogLevel            gotelemetry.LogLevel
	Filter              *regexp.Regexp
	ForceRunOnce        bool
	IsPiping            bool
	UseJSONPatch        bool
	UsePOST             bool
	IsNotifying         bool
	DebugMode           bool
	NotificationChannel string
	NotificationFlow    string
	Notification        gotelemetry.Notification
	WantsFunctionHelp   bool
	FunctionHelpName    string
	OAuthCommand        OAuthCommand
	OAuthName           string
	OAuthCode           string
	OAuthVerifier       string
	OAuthRealmID        string
}

// CLIConfig is accessed throughout the Agent to check startup configurations
var CLIConfig CLIConfigType

func banner(version string) {
	println()
	println("Telemetry Agent v" + version)
	println()
	println("Copyright © 2012-2016 Telemetry Inc.")
	println()
	println("For license information, see the LICENSE file")
	println("---------------------------------------------")
	println()
}

// Init the Agent by initializing flags and displaying on screen data
func Init(version string) {
	gotelemetry.UserAgentString = "Telemetry Agent v" + version
	banner(version)

	app := kingpin.New("telemetry_agent", "The Telemetry Agent")

	app.Version(version)

	app.Flag("config", "Path to the configuration file for this agent.").Short('c').StringVar(&CLIConfig.ConfigFileLocation)

	app.Flag("path", "Path to the database file for this agent.").Short('p').StringVar(&CLIConfig.DatabasePath)
	app.Flag("ttl", "The maximum lifespan of all series data in the Database.").StringVar(&CLIConfig.DatabaseTTL)
	app.Flag("auth_key", "The Authentication Key used for TelemetryTV to connect to the Agent.").Short('k').StringVar(&CLIConfig.AuthenticationKey)
	app.Flag("listen", "The port that the Agent's API will use to listen for TelemetryTV.").Short('l').StringVar(&CLIConfig.AuthenticationPort)
	app.Flag("certfile", "").StringVar(&CLIConfig.CertFile)
	app.Flag("keyfile", "").StringVar(&CLIConfig.KeyFile)

	app.Flag("apiurl", "Set the URL to the Telemetry API").Short('a').Default("https://api.telemetrytv.com").StringVar(&CLIConfig.APIURL)
	logLevel := app.Flag("verbosity", "Set the verbosity level (`debug`, `info`, `error`).").Short('v').Default("info").Enum("debug", "info", "error")
	filter := app.Flag("filter", "Run only the jobs whose IDs (or tags if no ID is specified) match the given regular expression").Default(".").String()
	app.Flag("debug", "Run scripts in debug mode. No API calls will be made. All output will be printed to the console.").BoolVar(&CLIConfig.DebugMode)

	once := app.Command("once", "Run all jobs exactly once and exit.")

	pipe := app.Command("pipe", "Accept a Rails-style HTTP PATCH Telemetry payload from stdin, send it to the API, and then exit.")
	pipe.Flag("channel", "The tag of the channel to which the update is sent.").StringVar(&CLIConfig.ChannelTag)
	pipe.Flag("jsonpatch", "With --pipe, submit the package as a JSON-Patch request instead. Ignored otherwise.").BoolVar(&CLIConfig.UseJSONPatch)
	pipe.Flag("post", "With --pipe, submit the package as a POST request instead. Ignored otherwise.").BoolVar(&CLIConfig.UsePOST)

	notify := app.Command("notify", "Send a channel notification.")
	notify.Flag("channel", "The tag of the channel to which the notification is sent. Either channel or flow is required.").StringVar(&CLIConfig.NotificationChannel)
	notify.Flag("flow", "The Tag of the Flow to whose channel the notification is sent. Either channel or flow is required.").StringVar(&CLIConfig.NotificationFlow)
	notify.Flag("title", "The title of the notification.").Required().StringVar(&CLIConfig.Notification.Title)
	notify.Flag("message", "The message to be displayed in the notification.").Required().StringVar(&CLIConfig.Notification.Message)
	notify.Flag("icon", "An icon to be displayed in the notification.").StringVar(&CLIConfig.Notification.Icon)
	notify.Flag("duration", "The amount of seconds for which the notification must be displayed.").Default("1").IntVar(&CLIConfig.Notification.Duration)
	notify.Flag("sound", "A URL to a notification sound (use `default` for Telemetry's default notification sound).").StringVar(&CLIConfig.Notification.SoundURL)

	oauthRequest := app.Command("oauth-request", "Request an oAuth authorization token")
	oauthRequest.Flag("name", "The name of the oAuth entry").Short('n').StringVar(&CLIConfig.OAuthName)

	oauthExchange := app.Command("oauth-exchange", "Exchange an oAuth authorization code")
	oauthExchange.Flag("name", "The name of the oAuth entry").Short('n').StringVar(&CLIConfig.OAuthName)
	oauthExchange.Flag("code", "The authorization code received from the provider").Short('o').StringVar(&CLIConfig.OAuthCode)
	oauthExchange.Flag("verifier", "The verifier code received from the provider").Short('e').StringVar(&CLIConfig.OAuthVerifier)
	oauthExchange.Flag("realm", "The realm ID received from the provider").Short('r').StringVar(&CLIConfig.OAuthRealmID)

	run := app.Command("run", "Runs the jobs scheduled in the configuration file provided.").Default()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case once.FullCommand():
		CLIConfig.ForceRunOnce = true

	case pipe.FullCommand():
		CLIConfig.IsPiping = true

	case notify.FullCommand():
		CLIConfig.IsNotifying = true

	case oauthRequest.FullCommand():
		CLIConfig.OAuthCommand = OAuthCommands.Request

	case oauthExchange.FullCommand():
		CLIConfig.OAuthCommand = OAuthCommands.Exchange

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
