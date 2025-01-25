package imagepng

import "fmt"

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

func (png *PngImage) TryConvert() {
	if png.IHDRchunk.chunkData.colortype != 2 {
		return
	}

	var data [][]byte
	data = make([][]byte, len(png.idatDecoded))

	var i int

	for i = 0; i < len(png.idatDecoded); i += 3 {

		data[i/3] = png.idatDecoded[i : i+3]
	}

	fmt.Println(i, png.IHDRchunk.chunkData.widthImg*png.IHDRchunk.chunkData.heightImg*3)

	// for _, rgb := range data {
	// 	fmt.Printf("%d", ((int(rgb[0]) * 65536) + (int(rgb[1]) * 256) + int(rgb[2])))
	// }

	fmt.Println()
}
