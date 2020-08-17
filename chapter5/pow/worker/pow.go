package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var characterSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

var hashes = []int{}

var CloseChan chan int

func start(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	prefix := []byte(params["prefix"][0])
	difficulty, _ := strconv.ParseInt(params["difficulty"][0], 10, 32)
	if completion, ok := params["completionEndpoint"]; ok {
		completion := completion[0]
		CloseChan = make(chan int, 1)
		go POWWithCallBack(prefix, int(difficulty), completion)
	} else {
		s := POW(prefix, int(difficulty))
		w.WriteHeader(200)
		w.Write([]byte(s))
	}
}

func cancel(w http.ResponseWriter, r *http.Request) {
	CloseChan <- 1
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/start", start) // ?prefix=&difficulty=&completionEndpoint=
	r.HandleFunc("/cancel", cancel)
	return r
}

func main() {
	http.ListenAndServe(":8080", Router())
}

func POWWithCallBack(prefix []byte, difficulty int, completionEndpoint string) {
	solutionStr := POW(prefix, difficulty)
	v := url.Values{}
	v.Set("solution", solutionStr)
	_, err := http.Get(completionEndpoint + "?" + v.Encode())
	if err != nil {
		fmt.Println(err)
	}
}

func POW(prefix []byte, difficulty int) string {
	solutionChan := make(chan []byte, 1)
	numberOfCPU := runtime.NumCPU()
	numberOfGoRo := numberOfCPU/3 + 1
	for i := 0; i < numberOfGoRo; i++ {
		runPOWInstance(prefix, 3, solutionChan, CloseChan, difficulty)
	}
	solution := <-solutionChan
	return string(solution)
}

//Create a small poof of work instance with x+1 goroutine.
// 1 random gen -> x number of hasher. Each hasher check if it compute the right hash. If the hash compute the right hash, stop all other hasher
func runPOWInstance(prefix []byte, numberOfHasher int, solutionChannel chan []byte, closeChan chan int, numberOfBits int) {
	blockSize := 1024
	size := numberOfHasher * 2

	unprocessIndex := make(chan int, size)
	processIndex := make(chan int, size)
	offset := len(prefix)

	blocks := make([][][]byte, size)
	for idx := range blocks {
		unprocessIndex <- idx
		blocks[idx] = make([][]byte, blockSize)
		for i := 0; i < blockSize; i++ {
			blocks[idx][i] = make([]byte, 20)
			blocks[idx][i] = append(prefix, blocks[idx][i]...)
		}
	}

	go func() {
		seed := uint64(time.Now().Local().UnixNano())
		for {
			select {
			case blockIndex := <-unprocessIndex:
				for idx := range blocks[blockIndex] {
					seed = RandomString(blocks[blockIndex][idx], offset, seed)
				}
				processIndex <- blockIndex
			case _ = <-closeChan:
				closeChan <- 1
				return
			}
		}
	}()

	for i := 0; i < numberOfHasher; i++ {
		index := len(hashes)
		hashes = append(hashes, 0)
		go func(index int, hashIndex int) {
			var hash bool
			for {
				select {
				case blockIndex := <-processIndex:
					hashes[hashIndex] += blockSize
					for _, random := range blocks[blockIndex] {
						hash = Hash(random, numberOfBits)
						if hash {
							solutionChannel <- random
							return
						}
					}
					unprocessIndex <- blockIndex
				case _ = <-closeChan:
					closeChan <- 1
					return
				}
			}
		}(i, index)
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
