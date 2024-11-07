package ioops

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

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

func PrintCatsTable(fileName string, cats []*cat.Cat) {

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	dashes := strings.Repeat("-", utf8.RuneCountInString(cats[0].String()))

	header := fmt.Sprintf("|%-6s|%-15s|%-6s|%-64s|", "Amžius", "Vardas", "Svoris", "Hash")
	file.WriteString(dashes + "\n")
	file.WriteString(header + "\n")
	file.WriteString(dashes + "\n")

	for _, cat := range cats {
		file.WriteString(cat.String() + "\n")
	}
	file.WriteString(dashes + "\n")
	file.WriteString("\n")
}
