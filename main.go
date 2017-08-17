package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kisom/randfault/rand"
)

var minTime = 5 * time.Minute
var maxTime = 24 * time.Hour

var response = struct {
	code    int
	message string
}{}

const version = "1.0.0"

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	response.code = http.StatusOK
	response.message = fmt.Sprintf("randfail v%s started at %d on %s",
		version, time.Now().Unix(), hostname)
}

var index = defaultIndex

func defaultIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("request for %s by %s", r.URL.Path, r.RemoteAddr)
	w.WriteHeader(response.code)
	w.Write([]byte(response.message))
}

func main() {
	var flAddress, failMode string
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	var address string
	var failProb float64

	if port != "" {
		address = host + ":" + port
	}

	flag.StringVar(&flAddress, "a", ":8080", "address to listen on")
	flag.StringVar(&failMode, "m", "timed", "failure mode to test")
	flag.DurationVar(&minTime, "min", minTime, "lower bound for timed failure")
	flag.DurationVar(&maxTime, "max", maxTime, "upper bound for a timed failure")
	flag.Float64Var(&failProb, "p", 0.5, "probability service will fail")
	flag.Parse()

	if address == "" {
		address = flAddress
	}

	willFail := rand.Coin(failProb)

	// see failure.go for this.
	failer, ok := failureModes[failMode]
	if !ok {
		log.Fatal("invalid failure mode ", failMode)
	}

	http.HandleFunc("/", index)

	if willFail {
		log.Println("service will fail")
		go failer()
	}

	log.Println("starting listener on", address)
	http.ListenAndServe(address, nil)
}
