package main

import (
	"errors"
	"image/png"
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

	_, err = try.TakePixet()
	asserterror.Assert(err != nil, "Can't take the pixel", err)

	png

}
