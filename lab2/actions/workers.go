package actions

import (
	"slices"
	"strings"

	"marius.org/cat"
	"marius.org/hasher"
	"marius.org/requests"
)

func ProcessData(itemsCount *int, dataChannel chan<- requests.DataRequest, resultChannel chan<- requests.ResultRequest, workerDone chan<- struct{}) {
	for {
		req := requests.DataRequest{Action: "remove", Response: make(chan *cat.Cat)}
		dataChannel <- req
		cat := <-req.Response
		*itemsCount = *itemsCount + 1
		if *itemsCount >= 25 {
			workerDone <- struct{}{}
			return
		}

		// Process the cat if it meets the weight condition
		if cat != nil && cat.Weight > 6 {
			hash := hasher.HashSha256(cat)
			cat.UpdateHash(hash)
			resReq := requests.ResultRequest{Cat: cat, Request: make(chan bool)}
			resultChannel <- resReq
			<-resReq.Request // Wait for confirmation of result processing
		}

	}
}

func ProcessDataThread(adderChan <-chan requests.DataRequest, removerChan <-chan requests.DataRequest, done <-chan struct{}) {
	cats := make([]*cat.Cat, 0, 10)

	for {
		select {
		case req := <-adderChan:
			// Add cat to slice if not full
			if len(cats) < cap(cats) {
				cats = append(cats, req.Cat)
				req.Response <- nil

			} else {
				req.Response <- nil
			}

		case req := <-removerChan:
			// Remove cat from slice and return it
			if len(cats) > 0 {
				req.Response <- cats[len(cats)-1]
				cats = cats[:len(cats)-1]

			} else {
				req.Response <- nil
			}

		case <-done:
			// Terminate when done channel is closed
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
func ProcessResultThread(resultChannel chan requests.ResultRequest, done <-chan struct{}, mainChan chan []*cat.Cat) {
	results := make([]*cat.Cat, 0, 30)

	for {
		select {
		case req := <-resultChannel:
			results = sortedInsert(results, req.Cat)
			req.Request <- true // received results

		case <-done:
			mainChan <- results
			return
		}

	}
}
