package imagepng

import (
	// "encodinry"
	"fmt"
	"os"
	"strconv"

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

	var newImage PngImage

	lenght := fmt.Sprintf("%x", bite[8:12])
	chunkType := fmt.Sprintf("%x", bite[12:16])

	if chunkType == "49484452" {

		newImage.IHDRchunk.dimention = lenght
		newImage.IHDRchunk.chunktype = chunkType

		newImage.IHDRchunk.chunkData.widthImg = fmt.Sprintf("%x", bite[16:20])
		newImage.IHDRchunk.chunkData.heightImg = fmt.Sprintf("%x", bite[20:24])
		newImage.IHDRchunk.chunkData.bitDepth = fmt.Sprintf("%x", bite[24])
		newImage.IHDRchunk.chunkData.colortype = fmt.Sprintf("%x", bite[25])
		newImage.IHDRchunk.chunkData.compressionMethod = fmt.Sprintf("%x", bite[26])
		newImage.IHDRchunk.chunkData.filterMethod = fmt.Sprintf("%x", bite[27])
		newImage.IHDRchunk.chunkData.interlaceMethod = fmt.Sprintf("%x", bite[28])
	}

	fmt.Println(strconv.ParseInt(newImage.IHDRchunk.chunkData.widthImg, 16, 64))
	fmt.Println(newImage)

}
