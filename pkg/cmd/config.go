package cmd

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/prometheus/common/version"
)

// CLIConfig defines all fields used to configure the Github Application. It will
// be ingested by Kong to generate the CLI.
type CLIConfig struct {
	Github struct {
		APIVersion string   `name:"github.api_version" help:"Github API version" env:"GQA_GITHUB_API_VERSION" default:"v3" enum:"v3,v4"`
		APIUrl     *url.URL `name:"github.addr" help:"Github API url" env:"GQA_GITHUB_ADDR" default:"https://api.github.com/"`

		Application struct {
			IntegrationID int      `name:"github.app_id" help:"Github Application ID" env:"GQA_GITHUB_APP_ID" required:""`
			WebhookSecret string   `name:"github.webhook_secret" help:"Github Webhook secret" env:"GQA_GITHUB_WEBHOOK_SECRET" required:""`
			Pkey          string   `name:"github.pkey" help:"Github Application private key path" env:"GQA_GITHUB_PKEY"`
			PkeyFile      *os.File `name:"github.pkey_file" help:"Github Application private key path"`
		} `embed:""`

		UserAgent string `name:"user_agent" help:"User agent used on Github requests" env:"GQA_USER_AGENT" default:"github-quick-action/${version}"`
	} `embed:""`

	ListenAddr string `name:"listen.addr" help:"Webhook listening address" env:"GQA_LISTEN_ADDR" default:"localhost:3000"`
	ListenPath string `name:"listen.path" help:"Webhook listening path" env:"GQA_LISTEN_PATH" default:"/api/v1/webhook"`
	LogLevel   string `name:"log.level" help:"Log level verbosity" env:"GQA_LOG_LEVEL" default:"info" enum:"trace,debug,info,warn,error,fatal,panic"`

	Version kong.VersionFlag
}

// FromEnvironment fills the CLIConfig from the environment. It doesn't use Kong in order to reduce each
// call on serverless platform.
func FromEnvironment() (CLIConfig, error) {
	config := CLIConfig{}
	for key, fncs := range envVarsDefinitions {
		value, set := os.LookupEnv(key)

		var err error
		if set {
			err = fncs.set(&config, value)
		} else {
			err = fncs.defaults(&config)
		}

		if err != nil {
			return CLIConfig{}, fmt.Errorf("invalid variable %s: %w", key, err)
		}
	}

	return config, nil
}

// Validate checks if some specific fields are properly configured.
func (c CLIConfig) Validate() error {
	// at least pkey or pkey_file should be defined
	if c.Github.Application.Pkey == "" && c.Github.Application.PkeyFile == nil {
		return fmt.Errorf("github.pkey or github.pkey_file is required")
	}

	return nil
}

// GithubAppConfig returns a githubapp.Config for the given CLI
func (c CLIConfig) GithubAppConfig() (githubapp.Config, error) {
	config := githubapp.Config{}
	switch c.Github.APIVersion {
	case "v3":
		config.V3APIURL = c.Github.APIUrl.String()
	case "v4":
		config.V4APIURL = c.Github.APIUrl.String()
	}

	config.App.IntegrationID = int64(c.Github.Application.IntegrationID)
	config.App.WebhookSecret = c.Github.Application.WebhookSecret

	config.App.PrivateKey = c.Github.Application.Pkey
	if config.App.PrivateKey == "" {
		pkey, err := io.ReadAll(c.Github.Application.PkeyFile)
		if err != nil {
			return githubapp.Config{}, fmt.Errorf("failed to read Github private key: %w", err)
		}
		config.App.PrivateKey = string(pkey)
	}

	return config, nil
}

// envVarsDefinitions defines hardcoded method to generate CLIConfig from environment variables. This is required
// to reduce as much as possible CPU time consumed by each call on serverless platform (Kong is a bit too heavier).
var envVarsDefinitions = map[string]struct {
	defaults func(*CLIConfig) error
	set      func(*CLIConfig, string) error
}{
	"GQA_GITHUB_API_VERSION": {
		defaults: func(config *CLIConfig) (err error) { config.Github.APIVersion = "v3"; return },
		set:      func(config *CLIConfig, s string) (err error) { config.Github.APIVersion = s; return },
	},
	"GQA_GITHUB_ADDR": {
		defaults: func(config *CLIConfig) (err error) {
			config.Github.APIUrl, err = url.Parse("https://api.github.com")
			return
		},
		set: func(config *CLIConfig, s string) (err error) { config.Github.APIUrl, err = url.Parse(s); return },
	},
	"GQA_USER_AGENT": {
		defaults: func(config *CLIConfig) (err error) {
			config.Github.UserAgent = "github-quick-action/" + version.Info()
			return
		},
		set: func(config *CLIConfig, s string) (err error) { config.Github.UserAgent = s; return },
	},

	"GQA_GITHUB_APP_ID": {
		defaults: func(_ *CLIConfig) error { return fmt.Errorf("variable is required") },
		set: func(c *CLIConfig, s string) (err error) {
			c.Github.Application.IntegrationID, err = strconv.Atoi(s)
			return
		},
	},
	"GQA_GITHUB_WEBHOOK_SECRET": {
		defaults: func(_ *CLIConfig) error { return fmt.Errorf("variable is required") },
		set:      func(config *CLIConfig, s string) (err error) { config.Github.Application.WebhookSecret = s; return },
	},
	"GQA_GITHUB_PKEY": {
		defaults: func(_ *CLIConfig) error { return fmt.Errorf("variable is required") },
		set:      func(config *CLIConfig, s string) (err error) { config.Github.Application.Pkey = s; return },
	},

	"GQA_LISTEN_ADDR": {
		defaults: func(config *CLIConfig) (err error) { config.ListenAddr = "localhost:3000"; return },
		set:      func(config *CLIConfig, s string) (err error) { config.ListenAddr = s; return },
	},
	"GQA_LISTEN_PATH": {
		defaults: func(config *CLIConfig) (err error) { config.ListenPath = "/api/v1/webhook"; return },
		set:      func(config *CLIConfig, s string) (err error) { config.ListenPath = s; return },
	},
	"GQA_LOG_LEVEL": {
		defaults: func(config *CLIConfig) (err error) { config.LogLevel = "info"; return },
		set:      func(config *CLIConfig, s string) (err error) { config.LogLevel = s; return },
	},
}
