package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

var workers = []string{"http://localhost:8080"}
var solutionChan = make(chan string, 1)

func start(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	prefix := params["prefix"][0]
	difficulty := params["difficulty"][0]
	v := url.Values{}
	v.Set("prefix", prefix)
	v.Set("difficulty", difficulty)
	v.Set("completionEndpoint", "http://localhost:8081/completion")
	workerParams := v.Encode()
	for _, worker := range workers {
		_, _ = http.Get(worker + "/start?" + workerParams)
	}
	fmt.Fprintf(w, <-solutionChan)
}

func solutionFound(w http.ResponseWriter, r *http.Request) {
	solution := r.URL.Query()["solution"][0]
	solutionChan <- solution
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/start", start) // ?prefix=&difficulty=
	r.HandleFunc("/completion", solutionFound)
	http.ListenAndServe(":8081", r)
}
