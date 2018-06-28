package rundeck

import "os"

const (
	// DefaultAPIVersion23 is defaulted to api version 23
	DefaultAPIVersion23 = 23

	// EnvRundeckToken sets the name of the environment variable to read
	EnvRundeckToken = "RDECK_TOKEN"
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
// set in the environment variable RDECK_TOKEN
func DefaultConfig() *Config {
	return &Config{
		APIVersion:       DefaultAPIVersion23,
		RundeckAuthToken: os.Getenv(EnvRundeckToken),
		ServerURL:        "http://127.0.0.1:4440",
	}
}
