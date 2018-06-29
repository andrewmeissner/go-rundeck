package rundeck

import "os"

const (
	// DefaultAPIVersion24 is defaulted to api version 23
	DefaultAPIVersion24 = 24

	// EnvRundeckToken sets the name of the environment variable to read
	EnvRundeckToken = "RUNDECK_TOKEN"

	// EnvRundeckServerURL is the name of the environment variable for the server url
	EnvRundeckServerURL = "RUNDECK_SERVER_URL"

	localRundeckURL = "http://localhost:4440"
)

// Config is the basic configuration needed by the client to communicate with Rundeck
type Config struct {
	// ServerURL is expected in the given format, ie: http://localhost:4440, or https://my.rundeck.com.
	ServerURL string

	// APIVersion is the version of the api you want to use.
	APIVersion int

	// RundeckAuthToken is the authentication token used to communicate with Rundeck
	RundeckAuthToken string
}

// DefaultConfig implements a localhost basic configuration, relying on and assuming a valid api token
// set in the environment variable RUNDECK_TOKEN.
//
// The environment variable RUNDECK_SERVER_URL will be used if it is present,
// otherwise http://localhost:4440 will be used as the server url.
func DefaultConfig() *Config {
	serverURL := os.Getenv(EnvRundeckServerURL)
	if serverURL == "" {
		serverURL = localRundeckURL
	}

	return &Config{
		APIVersion:       DefaultAPIVersion24,
		RundeckAuthToken: os.Getenv(EnvRundeckToken),
		ServerURL:        serverURL,
	}
}
