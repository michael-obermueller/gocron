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

type CronJob struct {
	Cmd string
}

// Run is called by the scheduler to execute the cron job
func (j *CronJob) Run() {
	report := fmt.Sprintf("Command %s output:\n", j.Cmd)
	if output, err := exec.Command(j.Cmd).CombinedOutput(); err == nil {
		report += string(output)
	} else {
		report += err.Error()
	}

	http.Post(reportToUrl, "text/plain", strings.NewReader(report))
}

// registerJobHandler is called for incoming job registration requests
func registerJobHandler(w http.ResponseWriter, r *http.Request) {
	cmd := filepath.Join(workingDir, r.URL.Query().Get("command"))
	schedule := r.URL.Query().Get("schedule")

	if err := jobRunner.AddJob(schedule, &CronJob{Cmd: cmd}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}