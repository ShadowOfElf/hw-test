package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	incomingString := "Hello, OTUS!"
	reverseString := reverse.String(incomingString)
	fmt.Println(reverseString)
}
