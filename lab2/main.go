package main

import (
	"fmt"

	jsonreader "marius.org/json_reader"
)

const FILE_NAME string = "cats.json"
const WORKERS int = 3

func main() {
	cats := jsonreader.ReadCatsFromJson(FILE_NAME)
	for _, cat := range cats {
		fmt.Println(cat)
	}

}
