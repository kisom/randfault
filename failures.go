package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kisom/randfault/rand"
)

var failureModes = map[string]func(){
	"timedeath":  timedFailure,
	"timesick":   timedSick,
	"hangindex":  hangIndex,
	"hanghealth": hangHealthCheck,
}

func init() {
	rand.Seed()
}

func timedFailure() {
	failAfter := rand.Between(int64(minTime), int64(maxTime))
	log.Println("fault: timed failure")
	time.Sleep(time.Duration(failAfter))
	os.Exit(1)
}

func timedSick() {
	failAfter := rand.Between(int64(minTime), int64(maxTime))
	time.Sleep(time.Duration(failAfter))
	log.Println("fault: timed sick")
	health.code = http.StatusInternalServerError
	health.message = "health check failed"
}

var emptyChan = make(chan []byte, 0)

func hungEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("hanging endpoint")
	w.WriteHeader(http.StatusOK)
	w.Write(<-emptyChan)
}

func hangIndex() {
	failAfter := rand.N(int64(maxTime))
	time.Sleep(time.Duration(failAfter))
	log.Println("fault: hanging index")
	index = hungEndpoint
}

func hangHealthCheck() {
	failAfter := rand.N(int64(maxTime))
	time.Sleep(time.Duration(failAfter))
	log.Println("fault: hanging health check")
	healthCheck = hungEndpoint
}
