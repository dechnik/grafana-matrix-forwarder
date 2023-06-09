package cfg

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// ResolveMode determines how the application will handle resolved alerts
type ResolveMode string

// AppSettings includes all application parameters
type AppSettings struct {
	VersionMode     bool
	UserID          string
	UserPassword    string
	HomeserverURL   string
	ServerHost      string
	MetricRounding  int
	ServerPort      int
	LogPayload      bool
	ResolveMode     ResolveMode
	PersistAlertMap bool
	AuthScheme      string
	AuthCredentials string
}

const (
	ResolveWithReaction    ResolveMode = "reaction"
	ResolveWithMessage     ResolveMode = "message"
	ResolveWithReply       ResolveMode = "reply"
	minServerPort                      = 1000
	maxServerPort                      = 65535
	defaultServerPort                  = 6000
	defaultServerHost                  = "0.0.0.0"
	defaultHomeServerUrl               = "matrix.org"
	defaultResolveMode                 = ResolveWithMessage
	defaultMetricRounding              = 3
	defaultPersistAlertMap             = true
)

// Parse the AppSettings data from the command line
func Parse() AppSettings {
	appSettings := &AppSettings{}
	appSettings.setDefaults()
	appSettings.updateSettingsFromEnvironment()
	appSettings.updateSettingsFromCommandLine()

	appSettings.validateConfiguration()
	return *appSettings
}

func (settings *AppSettings) setDefaults() {
	settings.ServerPort = defaultServerPort
	settings.ServerHost = defaultServerHost
	settings.HomeserverURL = defaultHomeServerUrl
	settings.ResolveMode = defaultResolveMode
	settings.MetricRounding = defaultMetricRounding
	settings.PersistAlertMap = defaultPersistAlertMap
}

func (settings *AppSettings) setResolveMode(resolveModeStr string) {
	resolveModeStrLower := strings.ToLower(resolveModeStr)
	if resolveModeStrLower == string(ResolveWithReaction) {
		settings.ResolveMode = ResolveWithReaction
	} else if resolveModeStrLower == string(ResolveWithMessage) {
		settings.ResolveMode = ResolveWithMessage
	} else if resolveModeStrLower == string(ResolveWithReply) {
		settings.ResolveMode = ResolveWithReply
	} else {
		log.Printf("invalid resolve mode provided (%s) - defaulting to %s", resolveModeStr, ResolveWithMessage)
		settings.ResolveMode = ResolveWithMessage
	}
}

func (settings *AppSettings) validateConfiguration() {
	var flagsValid = true
	if !settings.VersionMode {
		if settings.UserID == "" {
			fmt.Printf("missing parameter '-%s' or '%s'\n", userFlagName, userEnvName)
			flagsValid = false
		}
		if settings.UserPassword == "" {
			fmt.Printf("missing parameter '-%s' or '%s'\n", passwordFlagName, passwordEnvName)
			flagsValid = false
		}
		if settings.HomeserverURL == "" {
			fmt.Printf("missing parameter '-%s' or '%s'\n", homeServerFlagName, homeServerEnvName)
			flagsValid = false
		}
		if settings.ServerPort < minServerPort || settings.ServerPort > maxServerPort {
			fmt.Printf("invalid server port, must be within %d and %d (found %d)\n",
				minServerPort, maxServerPort, settings.ServerPort)
			flagsValid = false
		}
		if (settings.AuthScheme == "") != (settings.AuthCredentials == "") {
			fmt.Println("invalid auth setup - both scheme and credentials should be set")
			flagsValid = false
		}
		if strings.ToLower(settings.AuthScheme) != "" && strings.ToLower(settings.AuthScheme) != "bearer" {
			fmt.Println("unsupported auth scheme (supported: bearer)")
			flagsValid = false
		}
	}
	if !flagsValid {
		flag.Usage()
		os.Exit(1)
	}
}
