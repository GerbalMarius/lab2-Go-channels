package requests

import "marius.org/cat"

//used for storing calculated cat results
type ResultRequest struct {
	Cat     *cat.Cat
	Request chan bool
}
