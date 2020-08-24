package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"
)

var characterSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

var hashes = []int{}

var CloseChan chan int

func POWWithCallBack(prefix []byte, difficulty int, completionEndpoint string) {
	solutionStr := POW(prefix, difficulty)
	if solutionStr == "" {
		return
	}
	v := url.Values{}
	v.Set("solution", solutionStr)
	_, err := http.Get(completionEndpoint + "?" + v.Encode())
	if err != nil {
		fmt.Println(err)
	}
}

func POW(prefix []byte, difficulty int) string {
	numberOfCPU := runtime.NumCPU()
	solutionChan := make(chan []byte, 1)
	for idx := 0; idx < numberOfCPU; idx++ {
		//pass in idx to ensure it stay the same as idx can change value
		go func(hashIndex int) {
			seed := uint64(time.Now().Local().UnixNano())
			randomBytes := make([]byte, 100)
			randomBytes = append(prefix, randomBytes...)
			for {
				select {
				case <-CloseChan:
					CloseChan <- 1
					return
				case <-time.After(time.Nanosecond):
					count := 0
					for count < 5000 {
						count++
						seed = RandomString(randomBytes, len(prefix), seed)
						if Hash(randomBytes, difficulty) {
							solutionChan <- randomBytes
							CloseChan <- 1
							return
						}
					}
				}
			}
		}(idx)
	}

	select {
	case solution := <-solutionChan:
		CloseChan <- 1
		return string(solution)
	case _ = <-CloseChan:
		return ""
	}
}

func RandomNumber(seed uint64) uint64 {
	seed ^= seed << 21
	seed ^= seed >> 35
	seed ^= seed << 4
	return seed
}

func RandomString(str []byte, offset int, seed uint64) uint64 {
	for i := offset; i < len(str); i++ {
		seed = RandomNumber(seed)
		str[i] = characterSet[seed%62]
	}
	return seed
}

func Hash(data []byte, bits int) bool {
	bs := sha256.Sum256(data)
	nbytes := bits / 8
	nbits := bits % 8
	idx := 0
	for ; idx < nbytes; idx++ {
		if bs[idx] > 0 {
			return false
		}
	}
	return (bs[idx] >> (8 - nbits)) == 0
}

func start(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	prefix := []byte(params["prefix"][0])
	difficulty, _ := strconv.ParseInt(params["difficulty"][0], 10, 32)
	CloseChan = make(chan int, 1)
	if completion, ok := params["completionEndpoint"]; ok {
		completion := completion[0]
		go POWWithCallBack(prefix, int(difficulty), completion)
	} else {
		s := POW(prefix, int(difficulty))
		if s == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Write([]byte(s))
	}
}

func cancel(w http.ResponseWriter, r *http.Request) {
	select {
	case CloseChan <- 1:
		w.WriteHeader(200)
	case <-time.After(time.Second):
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to cancel last request"))
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(time.Now().String()))
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/start", start) // ?prefix=&difficulty=&completionEndpoint=
	r.HandleFunc("/cancel", cancel)
	r.HandleFunc("/health-check", ping)
	return r
}

var port = flag.String("port", "8080", "The port number to listen")
var masterNode = flag.String("master", "", "connect to a pool of workers")

func pingMaster() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	host := "http://" + hostname + ":" + (*port)
	values := url.Values{}
	values.Set("host", host)
	parms := values.Encode()
	for {
		_, err = http.Get(*masterNode + "/add-worker?" + parms)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Minute)
	}
}

func main() {
	flag.Parse()
	if masterNode != nil && len(*masterNode) > 0 {
		go pingMaster()
	}
	http.ListenAndServe(":"+(*port), Router())
}
