package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-cd <target-directory>")
		return
	}

	targetDir := os.Args[1]
	// Validate the directory (optional)
	fmt.Print(targetDir) // Print the directory path as output
}
