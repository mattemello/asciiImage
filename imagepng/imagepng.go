package imagepng

import (
	// "encodinry"
	"fmt"
	"os"
	"strconv"

	asserterror "github.com/mattemello/asciiImage/assertError"
)

func Image(filePath string) PngImage {
	fileImage, err := os.Open(filePath)
	asserterror.Assert(err != nil, "Can't open the file!!", err)
	defer fileImage.Close()

	state, err := fileImage.Stat()
	asserterror.Assert(err != nil, "Can't take the stat of the file", err)

	var bite = make([]byte, state.Size())

	fileImage.Read(bite)

	pngfile := fmt.Sprintf("%x", bite[:8])

	if pngfile != "89504e470d0a1a0a" {
		asserterror.AssertUnexpected("FILE NOT A PNG - the file provided is not a png file")
	}

	var newImage PngImage

	lenght := fmt.Sprintf("%x", bite[8:12])
	chunkType := string(bite[12:16])

	if chunkType != "IHDR" {
		asserterror.AssertUnexpected("NOT A IHDR AT FIRST - in the first position chunck there isn't the IHDR")
	}

	newImage.IHDRchunk.dimention = lenght
	newImage.IHDRchunk.chunktype = chunkType

	newImage.IHDRchunk.chunkData.widthImg, err = strconv.ParseInt(fmt.Sprintf("%x", bite[16:20]), 16, 64)
	asserterror.Assert(err != nil, "can't take the width of the image", err)

	newImage.IHDRchunk.chunkData.heightImg, err = strconv.ParseInt(fmt.Sprintf("%x", bite[20:24]), 16, 64)
	asserterror.Assert(err != nil, "can't take the height of the image", err)

	bitDepth, _ := strconv.ParseInt(fmt.Sprintf("%x", bite[24]), 16, 64)
	colortype, _ := strconv.ParseInt(fmt.Sprintf("%x", bite[25]), 16, 64)

	newImage.IHDRchunk.chunkData.bitDepth, newImage.IHDRchunk.chunkData.colortype = controllDepthandColor(bitDepth, colortype)

	compressionMethod := fmt.Sprintf("%x", bite[26])
	filterMethod := fmt.Sprintf("%x", bite[27])
	interlaceMethod := fmt.Sprintf("%x", bite[28])
	newImage.IHDRchunk.chunkData.compressionMethod = compressionMethod
	newImage.IHDRchunk.chunkData.filterMethod = filterMethod
	newImage.IHDRchunk.chunkData.interlaceMethod = interlaceMethod
	newImage.IHDRchunk.crc = fmt.Sprintf("%x", bite[28:32])

	dim := 33
	dimensionIDAT := dimensionIDAT(bite, dim)
	newImage.IDATchunks, newImage.IENDchunk = IDATake(bite, dim, dimensionIDAT)

	fmt.Println(newImage.IHDRchunk)
	fmt.Println(newImage.IENDchunk)

	return newImage

}

func controllDepthandColor(bitDepth, colorType int64) (int64, int64) {
	switch colorType {

	case 0:
		if bitDepth != 1 && bitDepth != 2 && bitDepth != 4 && bitDepth != 8 && bitDepth != 16 {
			asserterror.AssertUnexpected("error in the bitDepth and colorType value")
		}
	case 2:
		if bitDepth != 8 && bitDepth != 16 {
			asserterror.AssertUnexpected("error in the bitDepth and colorType value")
		}
	case 3:
		if bitDepth != 1 && bitDepth != 2 && bitDepth != 4 && bitDepth != 8 {
			asserterror.AssertUnexpected("error in the bitDepth and colorType value")
		}
	case 4:
		if bitDepth != 8 && bitDepth != 16 {
			asserterror.AssertUnexpected("error in the bitDepth and colorType value")
		}
	case 6:
		if bitDepth != 8 && bitDepth != 16 {
			asserterror.AssertUnexpected("error in the bitDepth and colorType value")
		}
	}

	return bitDepth, colorType

}

func dimensionIDAT(bite []byte, dim int) int {
	dimensionIdat := 0
	var lenght int64
	var chunkType string
	for {
		lenght, _ = strconv.ParseInt(fmt.Sprintf("%x", bite[dim:dim+4]), 16, 64)
		chunkType = string(bite[dim+4 : dim+8])
		if chunkType == "IDAT" {
			dimensionIdat++
		} else if chunkType == "IEND" {
			break
		} else {
			asserterror.AssertUnexpected("new type in the bit")
		}

		dim += int(lenght) + 4 + 8
	}

	return dimensionIdat
}

func IDATake(bite []byte, dim, dimensionIdat int) ([]IDAT, IEND) {
	var lenght int64
	var chunkType string
	var crc string
	var idatChunks = make([]IDAT, dimensionIdat)
	var iend IEND
	i := 0

	for {
		lenght, _ = strconv.ParseInt(fmt.Sprintf("%x", bite[dim:dim+4]), 16, 64)
		chunkType = string(bite[dim+4 : dim+8])
		data := bite[dim+8 : dim+8+int(lenght)]
		crc = fmt.Sprintf("%x", bite[dim+8+int(lenght):dim+12+int(lenght)])

		if chunkType == "IDAT" {
			idatChunks[i].dimention = lenght
			idatChunks[i].chunktype = chunkType
			idatChunks[i].chunkData = data
			idatChunks[i].crc = crc

			i++
		} else {
			iend.dimention = lenght
			iend.chunktype = chunkType
			iend.chunkData = data
			iend.crc = crc

			return idatChunks, iend
		}

		dim += int(lenght) + 4 + 8
	}
}
