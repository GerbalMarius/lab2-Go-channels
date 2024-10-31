package main

import (
	"fmt"

	"marius.org/hasher"
	jsonreader "marius.org/json_reader"
)

const FILE_NAME string = "cats.json"
const WORKERS int = 3

func main() {
	cats := jsonreader.ReadCatsFromJson(FILE_NAME)

	fmt.Println(cats[0].Serialize(), hasher.HashSha256(cats[0]))
}
