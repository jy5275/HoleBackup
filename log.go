package main

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("hole.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	logger = log.New(io.MultiWriter(os.Stdout, file), "", log.Lshortfile|log.LstdFlags)
}
