package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/report", cronJobResultHandler)
	http.ListenAndServe("127.0.0.1:8001", nil)
}

// cronJobResultHandler is called for incoming cron job results received from GoCron
func cronJobResultHandler(w http.ResponseWriter, r *http.Request) {
	if output, err := ioutil.ReadAll(r.Body); err == nil {
		fmt.Printf("%s\n\n", string(output))
	} else {
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
	}
}
