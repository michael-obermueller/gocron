package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/robfig/cron"
)

var jobRunner = cron.New()
var reportToUrl = "http://127.0.0.1:8001/report"
var workingDir string

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	workingDir = wd

	jobRunner.Start()
	http.HandleFunc("/register", registerJobHandler)
	http.ListenAndServe("127.0.0.1:8000", nil)
}

// CronJob 
type CronJob struct {
	Cmd       string
	OnSuccess func(cmd, output string)
	OnError   func(cmd string, err error)
}

// Run is called by the scheduler to execute the cron job
func (j *CronJob) Run() {
	if output, err := exec.Command(j.Cmd).CombinedOutput(); err == nil {
		j.OnSuccess(j.Cmd, string(output))
	} else {
		j.OnError(j.Cmd, err)
	}
}

// registerJobHandler is called for incoming job registration requests
func registerJobHandler(w http.ResponseWriter, r *http.Request) {
	cmd := filepath.Join(workingDir, r.URL.Query().Get("command"))
	schedule := r.URL.Query().Get("schedule")

	err := jobRunner.AddJob(schedule, &CronJob{Cmd: cmd,
		OnSuccess: func(cmd, output string) {
			body := fmt.Sprintf("Command %s succeeded, output:\n%s", cmd, output)
			http.Post(reportToUrl, "text/plain", strings.NewReader(body))
		},
		OnError: func(cmd string, err error) {
			body := fmt.Sprintf("Command %s failed: %s", cmd, err.Error())
			http.Post(reportToUrl, "text/plain", strings.NewReader(body))
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
