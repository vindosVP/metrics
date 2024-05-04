package main

import (
	"fmt"
	"os"
)

func test() {
	fmt.Println("Hello world!")
	os.Exit(1)
}

func main() {
	fmt.Println("Hello world!")
	os.Exit(1) // want "os.Exit call in main function"
}
