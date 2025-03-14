package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}

type model struct {
	bookmarks map[string]string
	cursor    int
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		bookmarks: map[string]string{
			"home": "C:\\Users\\f8col\\OneDrive\\Desktop\\Projects\\EBM",
			"src":  "C:\\Users\\f8col\\OneDrive\\Desktop\\Projects\\EBM\\src",
			"cmd":  "C:\\Users\\f8col\\OneDrive\\Desktop\\Projects\\EBM\\cmd",
		},
	}
}
