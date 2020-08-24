package main

import (
	"github.com/dustin/go-humanize"
	"runtime"
	"time"
	"fmt"
)

func powWithGoRoutine(prefix string, bitLength int) {
	start := time.Now()
	hash := []int{}
	totalHashesProcessed := 0

	numberOfCPU := runtime.NumCPU()
	closeChan := make(chan int, 1)
	solutionChan := make(chan []byte, 1)
	for idx := 0; idx < numberOfCPU; idx++ {
		hash = append(hash, 0)
		//pass in idx to ensure it stay the same as idx can change value
		go func(hashIndex int) {
			seed := uint64(time.Now().Local().UnixNano())
			randomBytes := make([]byte, 20)
			randomBytes = append([]byte(prefix), randomBytes...)
			for {
				select {
				case <-closeChan:
					closeChan <- 1
					return
				case <-time.After(time.Nanosecond):
					count := 0
					for count < 5000 {
						count++
						seed = RandomString(randomBytes, len(prefix), seed)
						if Hash(randomBytes, bitLength) {
							hash[hashIndex] += count
							solutionChan <- randomBytes
							closeChan <- 1
							return
						}
					}
					hash[hashIndex] += count
				}
			}
		}(idx)
	}
	<-solutionChan
	for _, v := range hash {
		totalHashesProcessed += v
	}
	end := time.Now()
	fmt.Println("time:", end.Sub(start).Seconds())
	fmt.Println("processed", humanize.Comma(int64(totalHashesProcessed)))
	fmt.Printf("processed/sec: %s\n", humanize.Comma(int64(float64(totalHashesProcessed)/end.Sub(start).Seconds())))
}
