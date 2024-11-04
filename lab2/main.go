package main

import (
	"fmt"

	"marius.org/cat"
	jsonreader "marius.org/json_reader"
	"marius.org/requests"
)

const FILE_NAME string = "cats.json"
const WORKERS int = 3

func main() {
	cats := jsonreader.ReadCatsFromJson(FILE_NAME)
	for _, cat := range cats {
		fmt.Println(cat)
	}
	dataChan := make(chan requests.DataRequest)
	resultChan := make(chan requests.ResultRequest)
	finished := make(chan struct{}) //end of operations
	workerDone := make(chan struct{})

	finalResults := make(chan *[]cat.Cat)

}
