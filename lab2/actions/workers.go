package actions

import (
	"slices"
	"strings"

	"marius.org/cat"
	"marius.org/hasher"
	"marius.org/requests"
)

func ProcessData(dataChannel chan<- requests.DataRequest, resultChannel chan<- requests.ResultRequest, workerDone chan<- struct{}) {
	req := requests.DataRequest{Action: "remove", Response: make(chan *cat.Cat)}
	dataChannel <- req
	cat := <-req.Response

	if cat != nil {
		if cat.Weight >= 5 {
			hash := hasher.HashSha256(cat)
			cat.UpdateHash(hash)
			resReq := requests.ResultRequest{Cat: cat, Request: make(chan bool)}
			resultChannel <- resReq
			<-resReq.Request //send request for results
		}
	}
	workerDone <- struct{}{} //send that worker is done
}

func ProcessDataThread(dataChannel <-chan requests.DataRequest, done <-chan struct{}) {
	capacity := 10
	cats := make([]*cat.Cat, 0, capacity)
	for {
		select {
		case req := <-dataChannel:
			if req.Action == "add" && len(cats) < capacity {
				cats = append(cats, req.Cat)
				req.Response <- nil
			}
			if req.Action == "remove" && len(cats) > 0 {
				req.Response <- cats[len(cats)-1] //removing the last element of the array
				cats = cats[:len(cats)-1]
			} else {
				req.Response <- nil // array is empty

			}
		case <-done:
			return
		}
	}
}
func sortedInsert(cats []*cat.Cat, item *cat.Cat) []*cat.Cat {
	if len(cats) == 0 {
		return append(cats, item)
	}
	idx, _ := slices.BinarySearchFunc(cats, item, func(a, b *cat.Cat) int {
		return strings.Compare(a.Name, b.Name)
	})
	cats = append(cats[:idx], append([]*cat.Cat{item}, cats[idx:]...)...)
	return cats

}
func ProcessResultThread(resultChannel <-chan requests.ResultRequest, done <-chan struct{}, mainChan chan<- *cat.Cat) {
	results := make([]*cat.Cat, 0, 30)
	idx := 0

	for {
		select {
		case req := <-resultChannel:
			results = sortedInsert(results, req.Cat)
			req.Request <- true
			mainChan <- results[idx]
			idx++

		case <-done:
			return
		}

	}
}
