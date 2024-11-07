package requests

import "marius.org/cat"

//used for add and remove requests
type DataRequest struct {
	Cat      *cat.Cat
	Response chan *cat.Cat //result to output if the cat was removed from array, nil otherwise
}
