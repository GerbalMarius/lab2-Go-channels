package main

import (
	"fmt"

	"marius.org/actions"
	"marius.org/cat"
	jsonreader "marius.org/json_reader"
	"marius.org/requests"
)

const FILE_NAME string = "cats.json"
const WORKERS int = 3

func main() {
	cats := jsonreader.ReadCatsFromJson(FILE_NAME)

	dataChan := make(chan requests.DataRequest, WORKERS)
	resultChan := make(chan requests.ResultRequest, WORKERS)
	done := make(chan struct{}) //channel for done requests

	workerDone := make(chan struct{}) // channel for worker finishes

	finalResults := make(chan *cat.Cat) //channel for processed results

	go actions.ProcessDataThread(dataChan, done)
	go actions.ProcessResultThread(resultChan, done, finalResults)

	for _, c := range cats {
		req := requests.DataRequest{Action: "add", Cat: c, Response: make(chan *cat.Cat)}
		dataChan <- req
		<-req.Response
	}

	for i := 0; i < WORKERS; i++ {
		go actions.ProcessData(dataChan, resultChan, workerDone)
	}
	for i := 0; i < WORKERS; i++ {
		<-workerDone
	}

	close(done)

	for i := 0; i < len(cats); i++ {
		result := <-finalResults
		fmt.Println(result)
	}
	close(dataChan)
	close(resultChan)
	close(finalResults)

}
