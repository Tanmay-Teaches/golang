package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var workersMutex = sync.RWMutex{}

//Keep track on how many failures for health check
var workers = map[string]int{}

var solutionChan = make(chan string, 1)

func start(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	prefix := params["prefix"][0]
	difficulty := params["difficulty"][0]
	v := url.Values{}
	v.Set("prefix", prefix)
	v.Set("difficulty", difficulty)
	v.Set("completionEndpoint", "http://"+*host+":"+*port+"/completion")
	workerParams := v.Encode()
	workersMutex.RLock()
	for worker, healthCheck := range workers {
		if healthCheck < 0 {
			continue
		}
		_, err := http.Get(worker + "/start?" + workerParams)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	workersMutex.RUnlock()
	solution := <-solutionChan
	w.Write([]byte(solution))
}

func solutionFound(w http.ResponseWriter, r *http.Request) {
	solution := r.URL.Query()["solution"][0]
	solutionChan <- solution
	workersMutex.RLock()
	defer workersMutex.RUnlock()
	for worker, healthCheck := range workers {
		if healthCheck != 0 {
			continue
		}
		_, _ = http.Get(worker + "/cancel")
	}
}

func addWorker(w http.ResponseWriter, r *http.Request) {
	worker := r.URL.Query()["host"][0]
	workersMutex.Lock()
	defer workersMutex.Unlock()
	if _, ok := workers[worker]; !ok {
		fmt.Println("Added: " + worker)
	}
	workers[worker] = 0
}

func pingWorker() {
	for {
		workersMutex.Lock()
		for worker, healthCheck := range workers {
			_ = worker
			if healthCheck < 0 {
				//Assume worker is dead.
				continue
			}
			// check on worker
			r, err := http.Get(worker + "/health-check")
			if err == nil && r.StatusCode == 200 {
				workers[worker] = 0
			} else {
				workers[worker] = healthCheck + 1
			}
			if healthCheck > 5 {
				workers[worker] = -1
			}
		}
		workersMutex.Unlock()
		time.Sleep(time.Second * 10)
	}
}

var port = flag.String("port", "8079", "The port number to listen")
var host = flag.String("host", "localhost", "host ip or service name")

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/start", start)
	r.HandleFunc("/completion", solutionFound)
	r.HandleFunc("/add-worker", addWorker)
	go pingWorker()
	http.ListenAndServe(*host+":"+(*port), r)
}
