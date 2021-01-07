package main

import (
	"errors"
	"fmt"
)

type Dictionary map[string] string
var ErrNotFound = errors.New("could not find the word you were looking for")


func (d Dictionary) Search (word string) (string, error)  {
	definition, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}

	return definition, nil
}

func main()  {
	d := Dictionary {"124": "234"}
	key := "124"
	v, ok := d["d"]
	if !ok {
		fmt.Println(ok)
	} else {
		fmt.Printf(v)
	}
	d["2"] = "sd"
	delete(d, "2")

	value, err := d.Search(key)
	if err != nil {
		fmt.Printf("Not key: %s\n", key)
	} else {
		fmt.Printf("Value for key %s is %s", key, value)
	}
	fmt.Printf("123434%d\n", 1)
}
