package main

import (
	"fmt"
	"os"

	// "github.com/mattemello/asciiImage/assertError"
	"github.com/mattemello/asciiImage/imagepng"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("ERROR - Not enought arguments")
		os.Exit(1)
	} else {
		//TODO: controll of the args
	}

	try := imagepng.Image(args[len(args)-1])

	try.TryConvert()
}
