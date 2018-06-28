package rundeck

// Config is the basic configuration needed by the client to communicate with Rundeck
type Config struct {
	// ServerURL is expected in the given format, ie: http://localhost:4440, or https://my.rundeck.com.
	ServerURL string

	// APIVersion is the version of the api you want to use.
	APIVersion int

	// RundeckAuthToken is the authentication token used to communicate with Rundeck
	RundeckAuthToken string
}
