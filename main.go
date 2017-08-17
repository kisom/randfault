package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/kisom/randfault/rand"
)

var minTime = 5 * time.Minute
var maxTime = 24 * time.Hour

type responseT struct {
	code    int
	message string
}

var (
	response responseT
	health   responseT
)

const version = "1.0.3"

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	response.code = http.StatusOK
	health.code = http.StatusOK
	response.message = fmt.Sprintf("randfail v%s started at %d on %s",
		version, time.Now().Unix(), hostname)
	health.message = "health check OK"
}

var (
	index       = defaultIndex
	healthCheck = defaultHealthCheck
)

func defaultIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("request for %s by %s", r.URL.Path, r.RemoteAddr)
	w.WriteHeader(response.code)
	w.Write([]byte(response.message))
}

func defaultHealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("health check from %s", r.RemoteAddr)
	w.WriteHeader(health.code)
	w.Write([]byte(health.message))
}

func listModes() {
	var modes []string
	for k := range failureModes {
		modes = append(modes, k)
	}

	sort.Strings(modes)
	fmt.Println("Valid failure modes:")
	for _, mode := range modes {
		fmt.Printf("\t%s\n", mode)
	}
	os.Exit(0)
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
	flag.StringVar(&failMode, "m", "timedeath", "failure mode to test")
	flag.DurationVar(&minTime, "min", minTime, "lower bound for timed failure")
	flag.DurationVar(&maxTime, "max", maxTime, "upper bound for a timed failure")
	flag.Float64Var(&failProb, "p", 0.5, "probability service will fail")
	flag.Parse()

	if failMode == "list" {
		listModes()
	}

	if address == "" {
		address = flAddress
	}

	willFail := rand.Coin(failProb)
	if failMode == "any" {
		var modes []string
		for k := range failureModes {
			modes = append(modes, k)
		}

		selector := rand.Between(0, int64(len(modes)))
		failMode = modes[selector]
		log.Println("selected failure mode", failMode)
	}

	// see failure.go for this.
	failer, ok := failureModes[failMode]
	if !ok {
		log.Fatal("invalid failure mode ", failMode)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index(w, r)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		healthCheck(w, r)
	})

	if willFail {
		log.Println("service will fail")
		go failer()
	}

	log.Println("starting listener on", address)
	http.ListenAndServe(address, nil)
}
