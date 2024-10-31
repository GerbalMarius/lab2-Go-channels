package cat

import (
	"encoding/json"
	"fmt"
)

type Cat struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
	hash   string
}

func NewCat(name string, age int, weight float64) *Cat {

	return &Cat{name, age, weight, ""}
}

func (c *Cat) UpdateHash(newHash string) {
	c.hash = newHash
}

func (c *Cat) String() string {
	return fmt.Sprintf("|%6d|%-15s|%6.2f|%s|", c.Age, c.Name, c.Weight, c.hash)
}

func (c *Cat) Serialize() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		panic("Couldn't encode cat to json")
	}
	return string(bytes)
}
