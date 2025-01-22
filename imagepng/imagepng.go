package imagepng

import (
	// "encodinry"
	"fmt"
	"os"

	asserterror "github.com/mattemello/asciiImage/assertError"
)

func Image(filePath string) {
	fileImage, err := os.Open(filePath)
	asserterror.Assert(err != nil, "Can't open the file!!", err)
	defer fileImage.Close()

	state, err := fileImage.Stat()

	var bite = make([]byte, state.Size())

	asserterror.Assert(err != nil, "Can't take the stat of the file", err)

	fileImage.Read(bite)

	pngfile := fmt.Sprintf("%x", bite[:8])

	if pngfile != "89504e470d0a1a0a" {
		fmt.Println("FILE NOT A PNG - the file provided is not a png file")
		return
	}

	fmt.Println()

	lenght := fmt.Sprintf("%x", bite[8:12])
	chunkType := fmt.Sprintf("%x", bite[12:16])

	if chunkType == "49484452" {
		fmt.Println("it's a IHDR!")
	}

	fmt.Println("png: ", pngfile)
	fmt.Println("lenght: ", lenght)
	fmt.Println("chunk type ex: ", chunkType)
	fmt.Println("lenght bt: ", bite[8:12])
	fmt.Println("chunk type: ", bite[12:16])
}
