package main

import (
	"github.com/dustin/go-humanize"
	"time"
	"fmt"
)

func powSimple(prefix string, bitLength int) {
	start := time.Now()

	totalHashesProcessed := 0
	seed := uint64(time.Now().Local().UnixNano())
	randomBytes := make([]byte, 20)
	randomBytes = append([]byte(prefix), randomBytes...)
	for {
		totalHashesProcessed++
		seed = RandomString(randomBytes, len(prefix), seed)
		if Hash(randomBytes, bitLength) {
			fmt.Println(string(randomBytes))
			break
		}
	}
	end := time.Now()

	fmt.Println("time:", end.Sub(start).Seconds())
	fmt.Println("processed", humanize.Comma(int64(totalHashesProcessed)))
	fmt.Printf("processed/sec: %s\n", humanize.Comma(int64(float64(totalHashesProcessed)/end.Sub(start).Seconds())))

}
