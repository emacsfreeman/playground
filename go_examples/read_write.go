package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func CreateFile() {
	file, err := os.Create("test.txt") // Truncates if file already exists, be careful!
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close() // Make sure to close the file when you're done

	len, err := file.WriteString("The Go Programming Language, also commonly referred to as Golang, is a general-purpose programming language, developed by a team at Google.")

	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s", file.Name())
}

func ReadFile() {
	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len(data))
	fmt.Printf("\nData: %s", data)
	fmt.Printf("\nError: %v", err)
}

func main() {
	fmt.Printf("########Create a file and Write the content #########\n")
	CreateFile()

	fmt.Printf("\n\n########Read file #########\n")
	ReadFile()
}
