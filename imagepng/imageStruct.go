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

type PLTE struct {
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
	chunkPLTE    PLTE
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

func (png *PngImage) Height() int {
	return int(png.chunkIHDR.chunkData.heightImg)
}

func (png *PngImage) TakePixet() ([][]byte, error) {

	fmt.Println(png.Width())

	switch png.chunkIHDR.chunkData.colortype {
	case bColorGray:
		return graySample(png.idatDecoded, int(png.chunkIHDR.chunkData.bitDepth), png.Height(), png.Width()), nil

	case bColorRGB:
		return rgbSample(png.idatDecoded, int(png.chunkIHDR.chunkData.bitDepth)), nil

	}

	return nil, errors.New("Can't take the pixel, color image not implemanted!")

}

func graySample(png []byte, depth, height, width int) [][]byte {
	var image = make([][]byte, len(png))

	switch depth {
	case bDepthOne:
		for i := 0; i < len(png); i += 8 {
			b := png[i/8]
			row := make([]byte, width)
			for j := 0; j < 8 && j+i < len(png); j++ {
				row[j] = (b >> 7) * 0xff
				b <<= 1
			}
			image[i] = row
		}
		break

	case bDepthTwo:
		for i := 0; i < len(png); i += 4 {
			b := png[i/8]
			row := make([]byte, width)
			for j := 0; j < 8 && j+i < len(png); j++ {
				row[j] = (b >> 6) * 0x55
				b <<= 1
			}
			image[i] = row
		}
		break

	case bDepthFour:
		for i := 0; i < len(png); i += 2 {
			b := png[i/8]
			row := make([]byte, width)
			for j := 0; j < 8 && j+i < len(png); j++ {
				row[j] = (b >> 4) * 0x11
				b <<= 1
			}
			image[i] = row
		}
		break

	case bDepthEight:
		for i := 0; i < len(png); i += 1 {
			b := png[i/8]
			row := make([]byte, width)
			for j := 0; j < 8 && j+i < len(png); j++ {
				row[j] = (b >> 4) * 0x11
				b <<= 1
			}
			image[i] = row
		}
		break

	}

	return image

}

func rgbSample(png []byte, depth int) [][]byte {

	switch depth {

	}

	var data [][]byte
	data = make([][]byte, len(png)/3+1)

	for i := 0; i < len(png); i += 3 {
		data[i/3] = png[i : i+3]
	}

	return data
}
