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
	IHDRchunk    IHCD
	IDATchunks   []IDAT
	IENDchunk    IEND
	idatDecoded  []byte
	positionIdat int
}

func (png *PngImage) Width() int {
	return int(png.IHDRchunk.chunkData.widthImg)
}

func (png *PngImage) TakePixet() ([][]byte, error) {

	switch png.IHDRchunk.chunkData.colortype {

	case 2:
		var data [][]byte
		data = make([][]byte, len(png.idatDecoded)/3+1)

		for i := 0; i < len(png.idatDecoded); i += 3 {
			data[i/3] = png.idatDecoded[i : i+3]
		}

		fmt.Println(data)
		return data, nil

	}

	return nil, errors.New("Can't take the pixel, color image not implemanted!")

}
