package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mattemello/asciiImage/assertError"
	"github.com/mattemello/asciiImage/imagepng"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		asserterror.Assert(true, "Not enought arguments", errors.New(""))
	} else {
		//TODO: controll of the args
	}

	try, err := imagepng.Image(args[len(args)-1])
	asserterror.Assert(err != nil, "Can't take the image binary", err)

	data, err := try.TakePixet()
	asserterror.Assert(err != nil, "Can't take the pixel", err)

	if data == nil {

	}

	file, err := os.Create("./imageAsci.txt")
	asserterror.Assert(err != nil, "can't create the file", err)

	for _, i := range data {
		for _, j := range i {
			if j == 0 {
				fmt.Fprintf(file, ".")
			} else {
				fmt.Fprintf(file, "&")
			}
		}
		fmt.Fprintln(file)

	}

}
