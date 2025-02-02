package imagepng

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

var plteIn = false

func Image(filePath string) (PngImage, error) {
	fileImage, err := os.Open(filePath)
	if err != nil {
		return PngImage{}, err
	}
	defer fileImage.Close()

	state, err := fileImage.Stat()
	if err != nil {
		return PngImage{}, err
	}

	var bite = make([]byte, state.Size())

	fileImage.Read(bite)

	pngfile := fmt.Sprintf("%x", bite[:8])

	if pngfile != "89504e470d0a1a0a" {
		return PngImage{}, errors.New("Invalid type of the file, it's not a png file")
	}

	var newImage PngImage

	lenght := fmt.Sprintf("%x", bite[8:12])
	chunkType := string(bite[12:16])

	if chunkType != "IHDR" {
		return PngImage{}, errors.New("NOT A IHDR AT FIRST - in the first position chunck there isn't the IHDR")
	}

	newImage.chunkIHDR.dimention = lenght
	newImage.chunkIHDR.chunktype = chunkType

	newImage.chunkIHDR.chunkData.widthImg, err = strconv.ParseInt(fmt.Sprintf("%x", bite[16:20]), 16, 64)
	if err != nil {
		return PngImage{}, err
	}

	newImage.chunkIHDR.chunkData.heightImg, err = strconv.ParseInt(fmt.Sprintf("%x", bite[20:24]), 16, 64)
	if err != nil {
		return PngImage{}, err
	}

	bitDepth, _ := strconv.ParseInt(fmt.Sprintf("%x", bite[24]), 16, 64)
	colortype, _ := strconv.ParseInt(fmt.Sprintf("%x", bite[25]), 16, 64)

	bitDepth, err = controllDepthandColor(bitDepth, colortype)
	if err != nil {
		return PngImage{}, err
	}

	newImage.chunkIHDR.chunkData.bitDepth = bitDepth
	newImage.chunkIHDR.chunkData.colortype = colortype

	compressionMethod := fmt.Sprintf("%x", bite[26])
	filterMethod := fmt.Sprintf("%x", bite[27])
	interlaceMethod := fmt.Sprintf("%x", bite[28])

	fmt.Println(filterMethod)

	newImage.chunkIHDR.chunkData.compressionMethod = compressionMethod
	newImage.chunkIHDR.chunkData.filterMethod = filterMethod
	newImage.chunkIHDR.chunkData.interlaceMethod = interlaceMethod
	newImage.chunkIHDR.crc = fmt.Sprintf("%x", bite[28:32])

	dim := 33

	dimensionIDAT := dimensionIDAT(bite, dim)
	newImage.chunksIDAT, newImage.chunkIEND, newImage.idatDecoded, err = IDATake(bite, dim, dimensionIDAT)
	if err != nil {
		return PngImage{}, err
	}

	if plteIn {
		newImage.chunkPLTE, err = controllPLTE(bite, dim)
		if err != nil {
			return PngImage{}, err
		}
	}

	newImage.positionIdat = 0

	return newImage, nil
}

func controllPLTE(bite []byte, i int) (PLTE, error) {
	var plte PLTE

	for {
		plte.dimention, _ = strconv.ParseInt(fmt.Sprintf("%x", bite[i:i+4]), 16, 64)
		plte.chunktype = string(bite[i+4 : i+8])
		plte.chunkData = bite[i+8 : i+8+int(plte.dimention)]
		plte.crc = fmt.Sprintf("%x", bite[i+8+int(plte.dimention):i+12+int(plte.dimention)])

		if plte.chunktype == "PLTE" {
			fmt.Println("enter")
			break
		}

		i += int(plte.dimention) + 4 + 8
	}

	if (plte.dimention % 3) != 0 {
		return PLTE{}, errors.New("The PLTE chunk is not good")

	}

	return plte, nil

}

func controllDepthandColor(bitDepth, colorType int64) (int64, error) {
	switch colorType {

	case bColorGray:
		if bitDepth != bDepthOne && bitDepth != bDepthTwo && bitDepth != bDepthFour && bitDepth != bDepthEight && bitDepth != bDepthSixteen {
			return 0, errors.New(fmt.Sprintf("Invalid colore and depth bit. Color: %d  Depth: %d", colorType, bitDepth))
		} else {
			fmt.Println("gray scale sample")
		}
	case bColorRGB:
		if bitDepth != bDepthEight && bitDepth != bDepthSixteen {
			return 0, errors.New(fmt.Sprintf("Invalid colore and depth bit. Color: %d  Depth: %d", colorType, bitDepth))
		} else {
			fmt.Println("each pixel is an RGB")
		}
	case bColorPLET:
		if bitDepth != bDepthOne && bitDepth != bDepthTwo && bitDepth != bDepthFour && bitDepth != bDepthEight {
			return 0, errors.New(fmt.Sprintf("Invalid colore and depth bit. Color: %d  Depth: %d", colorType, bitDepth))
		} else {
			// bitDepth alwais 8
			bitDepth = bDepthEight
			fmt.Println("each pixel is palete index; PLTE need check")
		}
	case bColorGrayAlpha:
		if bitDepth != bDepthEight && bitDepth != bDepthSixteen {
			return 0, errors.New(fmt.Sprintf("Invalid colore and depth bit. Color: %d  Depth: %d", colorType, bitDepth))
		} else {
			fmt.Println("gray scale and alpha")
		}
	case bColorRGBAlpha:
		if bitDepth != bDepthEight && bitDepth != bDepthSixteen {
			return 0, errors.New(fmt.Sprintf("Invalid colore and depth bit. Color: %d  Depth: %d", colorType, bitDepth))
		} else {
			fmt.Println("RGB and alpha")
		}

	default:
		return 0, errors.New(fmt.Sprintf("Invalid color type -> Color: %d", colorType))
	}

	return bitDepth, nil

}

func dimensionIDAT(bite []byte, dim int) int {
	dimensionIdat := 0
	var lenght int64
	var chunkType string
	for {
		lenght, _ = strconv.ParseInt(fmt.Sprintf("%x", bite[dim:dim+4]), 16, 64)
		chunkType = string(bite[dim+4 : dim+8])
		data := bite[dim+8 : dim+8+int(lenght)]
		if chunkType == "IDAT" {
			dimensionIdat++
		} else if chunkType == "IEND" {
			break
		} else {
			fmt.Println(chunkType, data)
			if chunkType == "PLTE" {
				plteIn = true
			}

		}

		dim += int(lenght) + 4 + 8
	}

	return dimensionIdat
}

func decodeIDAT(idatChunk []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(idatChunk))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var decoded []byte
	var buf = make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			decoded = append(decoded, buf[:n]...)
			break
		}
		if err != nil {
			return nil, err
		}

		decoded = append(decoded, buf[:n]...)
	}

	return decoded, nil
}

func IDATake(bite []byte, dim, dimensionIdat int) ([]IDAT, IEND, []byte, error) {
	var lenght int64
	var chunkType string
	var crc string
	var idatChunks = make([]IDAT, dimensionIdat)
	var allData []byte
	var iend IEND
	i := 0

	for {
		lenght, _ = strconv.ParseInt(fmt.Sprintf("%x", bite[dim:dim+4]), 16, 64)
		chunkType = string(bite[dim+4 : dim+8])
		data := bite[dim+8 : dim+8+int(lenght)]
		crc = fmt.Sprintf("%x", bite[dim+8+int(lenght):dim+12+int(lenght)])

		if chunkType == "IDAT" {

			allData = append(allData, data...)
			idatChunks[i].dimention = lenght
			idatChunks[i].chunktype = chunkType
			idatChunks[i].chunkData = data
			idatChunks[i].crc = crc

			i++
		} else if chunkType == "IEND" {
			iend.dimention = lenght
			iend.chunktype = chunkType
			iend.chunkData = data
			iend.crc = crc

			decodedImage, err := decodeIDAT(allData)
			if err != nil {
				return nil, IEND{}, nil, err
			}

			return idatChunks, iend, decodedImage, nil
		}

		dim += int(lenght) + 4 + 8
	}
}
