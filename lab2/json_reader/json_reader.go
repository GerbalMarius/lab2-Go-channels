package jsonreader

import (
	"encoding/json"
	"fmt"
	"os"

	"marius.org/cat"
)

func ReadCatsFromJson(fileName string) []*cat.Cat {
	cats := make([]*cat.Cat, 0, 25)

	bytes, readErr := os.ReadFile(fileName)
	if readErr != nil {
		fmt.Println("ERROR reading file:", fileName, readErr)
	}

	json.Unmarshal(bytes, &cats)
	return cats
}
