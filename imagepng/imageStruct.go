package imagepng

import (
	"errors"
	"fmt"
)

type IHCD struct {
	dimention string
	chunktype string
	chunkData IHCDdata
	crc       string
}

type IHCDdata struct {
	widthImg          int64
	heightImg         int64
	bitDepth          int64
	colortype         int64
	compressionMethod string
	filterMethod      string
	interlaceMethod   string
}

type IDAT struct {
	dimention int64
	chunktype string
	chunkData []byte
	crc       string
}

type IEND struct {
	dimention int64
	chunktype string
	chunkData []byte
	crc       string
}

type PngImage struct {
	chunkIHDR    IHCD
	chunksIDAT   []IDAT
	chunkIEND    IEND
	idatDecoded  []byte
	positionIdat int
}

const (
	bDepthInvalid = 0
	bDepthOne     = 1
	bDepthTwo     = 2
	bDepthFour    = 4
	bDepthEight   = 8
	bDepthSixteen = 16
)

const (
	bColorGray      = 0
	bColorRGB       = 2
	bColorPLET      = 3
	bColorGrayAlpha = 4
	bColorRGBAlpha  = 6
)

func (png *PngImage) Width() int {
	return int(png.chunkIHDR.chunkData.widthImg)
}

func (png *PngImage) TakePixet() ([][]byte, error) {

	fmt.Println(png.Width())

	switch png.chunkIHDR.chunkData.colortype {
	case bColorGray:
		break

	case bColorRGB:
		return rgbSample(png.idatDecoded, int(png.chunkIHDR.chunkData.bitDepth)), nil

	}

	return nil, errors.New("Can't take the pixel, color image not implemanted!")

}

func rgbSample(png []byte, depth int) [][]byte {

	var data [][]byte
	data = make([][]byte, len(png)/3+1)

	for i := 0; i < len(png); i += 3 {
		data[i/3] = png[i : i+3]
	}

	return data
}
