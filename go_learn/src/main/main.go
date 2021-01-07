package main

import (
	"io/ioutil"
	"log"
	"os"
)
func main()  {
	filename := "mr-*-0"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()
	println(string(content))
}