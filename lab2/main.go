package main

import (
	"fmt"
	"os"

	"marius.org/actions"
	"marius.org/cat"
	ioops "marius.org/io_ops"

	"marius.org/requests"
)

const FILE_NAME string = "cats.json"
const RESULTS string = "result.txt"
const WORKERS int = 3

func main() {
	if _, err := os.Stat(RESULTS); err == nil {
		os.Remove(RESULTS)
	}

	cats := ioops.ReadCatsFromJson(FILE_NAME)
	catsToProcess := 0
	fmt.Println("Before:", catsToProcess)
	ioops.PrintCatsTable(RESULTS, cats)

	adderChan := make(chan requests.DataRequest)
	removerChan := make(chan requests.DataRequest)
	resultChan := make(chan requests.ResultRequest)
	done := make(chan struct{}) //channel for done requests

	workerDone := make(chan struct{}) // channel for worker finishes

	finalResults := make(chan []*cat.Cat) //channel for processed results

	go actions.ProcessDataThread(adderChan, removerChan, done)
	go actions.ProcessResultThread(resultChan, done, finalResults)

	for i := 0; i < WORKERS; i++ {
		go actions.ProcessData(&catsToProcess, removerChan, resultChan, workerDone)
	}
	for _, c := range cats {
		req := requests.DataRequest{Action: "add", Cat: c, Response: make(chan *cat.Cat)}
		adderChan <- req
		<-req.Response

	}

	for i := 0; i < WORKERS; i++ {
		<-workerDone
	}
	close(done)
	fmt.Println("After:", catsToProcess)
	results := <-finalResults
	ioops.PrintCatsTable(RESULTS, results)

}
