package main

import "fmt"

func main() {
	var input string
	fmt.Print("Enter your name: ")
	fmt.Scanln(&input)
	fmt.Printf("Hello, %s!", input)
	fmt.Println()
}
