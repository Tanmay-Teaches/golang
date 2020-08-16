package main

/*
Example of using third party package with go module
1) Initialize module
	go mod init <module name>
		example: go mod init example
2) Add package to module cache using go module
	go mod download <url>
		example: go mod download github.com/otiai10/primes
*/
import (
	"fmt"
	"github.com/otiai10/primes"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		println("Usage:", os.Args[0], "<number>")
		os.Exit(1)
	}
	number, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}
	f := primes.Factorize(int64(number))
	fmt.Println("primes:", len(f.Powers()) == 1)
}
