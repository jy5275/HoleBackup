package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	INTERVAL_MID    = 5 * time.Second
	INTERVAL_RANGE  = 1 * time.Second
	LOG_FILE        = "hole.log"
	BACKTRACK_PAGES = 2
)

var logger *log.Logger

func Run() {
	hs := NewHoleStorage()

	for {
		timeToWait := time.Duration(rand.Int63n(int64(2*INTERVAL_RANGE))) + INTERVAL_MID - INTERVAL_RANGE
		logger.Printf("================== Sleep for %v ==================\n", timeToWait)
		time.Sleep(timeToWait)

		posts, err := GetLists(BACKTRACK_PAGES)
		if err != nil {
			logger.Println("failed to get list: ", err)
			continue
		}

		hs.InsertAndCheck(posts)
		logger.Printf("Deleted list: %+v\n", hs.GetAllDeleted())
	}
}

func main() {
	file, err := os.OpenFile("hole.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	logger = log.New(io.MultiWriter(os.Stdout, file), "", log.Lshortfile|log.LstdFlags)
	Run()
}
