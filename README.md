# go-rundeck

## Testing
Use the supplied Vagrantfile to spin up a local instance of Rundeck.  Login to http://localhost:4440 using `admin` and `admin` as the username and password.  Create an API token and set that to an environment variable called `RUNDECK_TOKEN`.

This environment variables must be set in the same session as running the tests, otherwise the tests will fail to authenticate with the containerized instance of Rundeck.