package main

import (
	"fmt"
	"log"
)

func ask() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	ok := "y"
	no := "n"
	if response == ok {
		return true
	} else if response == no {
		return false
	} else {
		fmt.Println("Please type [y, n] and then press enter:")
		return ask()
	}
}

// // You might want to put the following two functions in a separate utility package.
// // posString returns the first index of element in slice.
// // If slice does not contain element, returns -1.
// func posString(slice []string, element string) int {
// 	for index, elem := range slice {
// 		if elem == element {
// 			return index
// 		}
// 	}
// 	return -1
// }
//
// // containsString returns true iff slice contains element
// func containsString(slice []string, element string) bool {
// 	return !(posString(slice, element) == -1)
// }
