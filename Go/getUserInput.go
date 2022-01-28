package main

import "fmt"

func getUserInput() string {
	var archiveInput string
	fmt.Println("Want archive files? (yes or y/any key to not)")
	fmt.Scanln(&archiveInput)
	return archiveInput
}
