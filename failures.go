package main

import (
	"os"
	"time"

	"github.com/kisom/randfault/rand"
)

var failureModes = map[string]func(){
	"timed": timedFailure,
}

func init() {
	rand.Seed()
}

func timedFailure() {
	failAfter := rand.Between(int64(minTime), int64(maxTime))
	time.Sleep(time.Duration(failAfter))
	os.Exit(1)
}
