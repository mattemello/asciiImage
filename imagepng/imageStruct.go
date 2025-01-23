package imagepng

type IHCD struct {
	dimention string
	chunktype string
	chunkData IHCDdata
	crc       string
}

type IHCDdata struct {
	widthImg          string
	heightImg         string
	bitDepth          string
	colortype         string
	compressionMethod string
	filterMethod      string
	interlaceMethod   string
}

type PngImage struct {
	IHDRchunk IHCD
}
