package main

import (
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
	var catSize int = len(cats)
	ioops.PrintCatsTable(RESULTS, cats)

	adderChan := make(chan requests.DataRequest)    // channel for add operations
	removerChan := make(chan requests.DataRequest)  // channel for remove operations
	resultChan := make(chan requests.ResultRequest) // channel for processed items
	done := make(chan struct{})                     //channel for done requests
	sizeChan := make(chan int)

	finalResults := make(chan []*cat.Cat)   //channel for processed results
	finishedRemoving := make(chan struct{}) // signal for finished removals

	go actions.ProcessDataThread(catSize, adderChan, removerChan, finishedRemoving, done, sizeChan)
	go actions.ProcessResultThread(resultChan, done, finalResults)

	for i := 0; i < WORKERS; i++ {
		go actions.ProcessData(adderChan, removerChan, resultChan, finishedRemoving, sizeChan)
	}
	for _, c := range cats {
		req := requests.DataRequest{Cat: c, Response: make(chan *cat.Cat)}
		adderChan <- req
	}

	close(done)
	results := <-finalResults
	ioops.PrintCatsTable(RESULTS, results)

}
