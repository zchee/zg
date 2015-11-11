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
