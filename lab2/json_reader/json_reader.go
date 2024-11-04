package jsonreader

import (
	"encoding/json"
	"os"

	"marius.org/cat"
)

func ReadCatsFromJson(fileName string) []*cat.Cat {
	cats := make([]*cat.Cat, 0, 25)

	bytes, readErr := os.ReadFile(fileName)
	if readErr != nil {
		panic(readErr)
	}

	json.Unmarshal(bytes, &cats)
	return cats
}
