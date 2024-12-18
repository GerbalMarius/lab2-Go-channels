package actions

import (
	"fmt"
	"slices"
	"strings"

	"marius.org/cat"
	"marius.org/hasher"
	"marius.org/requests"
)

func ProcessData(adderChan chan requests.DataRequest, dataChannel chan<- requests.DataRequest, resultChannel chan<- requests.ResultRequest,
	finishedRemoving <-chan struct{}, sizeChan chan int) {
	for {
		currSize := <-sizeChan
		if currentSize <= 0 {
			<-adderChan
			continue
		} else {
			select {
			case <-finishedRemoving:
				return
			default:
				req := requests.DataRequest{Response: make(chan *cat.Cat)}
				dataChannel <- req
				cat := <-req.Response

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

	}
}

func ProcessDataThread(elementsToProcess int, adderChan <-chan requests.DataRequest, removerChan <-chan requests.DataRequest,
	finishedRemoving chan struct{}, done chan struct{}, sizeChan chan int) {
	cats := make([]*cat.Cat, 0, 10)
	removed := 0

	for {
		sizeChan <- len(cats)
		select {
		case req := <-adderChan:
			// Add cat to slice if not full
			if len(cats) < cap(cats) {
				cats = append(cats, req.Cat)

			}

		case req := <-removerChan:
			// Remove cat from slice and return it
			if len(cats) > 0 {
				req.Response <- cats[len(cats)-1]
				removed++
				fmt.Println("Removed element:", cats[len(cats)-1], "Current removals:", removed)
				cats = cats[:len(cats)-1]
				if removed == elementsToProcess {
					close(finishedRemoving)
				}

			} else {
				req.Response <- nil
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
