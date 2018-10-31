# gocron

gocron is a job scheduler used as sample in Dynatrace blog post [Introducing Custom Services for Go Applications](https://www.google.com). It accepts job registration requests on port 8000 and will execute the application specified in the job definition periodically.

Register jobs with:
- `curl -X POST "http://localhost:8000/register?command=APPLICATION&schedule=SCHEDULE"`
- e.g.: `curl -X POST "http://localhost:8000/register?command=script.sh&schedule=@every%205s"`

The console output of finished jobs is sent as plain text to an HTTP server listening on port 8001.

- Build GoCron with: `go build GoCron.go`
- Build Server with: `go build server/Server.go`
