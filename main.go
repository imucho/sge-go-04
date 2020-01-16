package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	stringEncode()
	stringDecode()
}

func loadPNG(name string) (image.Image, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func savePNG(name string, img *image.RGBA) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	png.Encode(file, img)
	return nil
}

func alphaEncode() {
	coverImg, err := loadPNG("cover.png")
	if err != nil {
		fmt.Println(err)
	}
	sourceImg, err := loadPNG("alpha_source.png")
	if err != nil {
		fmt.Println(err)
	}

	rect := coverImg.Bounds()
	if rect != sourceImg.Bounds() {
		fmt.Println("画像サイズが一致しません")
	}

	newImg := image.NewRGBA(image.Rectangle{rect.Min, rect.Max})

	for y := 0; y < rect.Dx(); y++ {
		for x := 0; x < rect.Dy(); x++ {
			r1, g1, b1, _ := coverImg.At(x, y).RGBA()
			r2, g2, b2, _ := sourceImg.At(x, y).RGBA()
			a := 255
			if uint8(r2) != 255 || uint8(g2) != 255 || uint8(b2) != 255 {
				a = 254
			}
			newImg.SetRGBA(x, y, color.RGBA{uint8(r1), uint8(g1), uint8(b1), uint8(a)})
		}
	}
	savePNG("alpha_encode.png", newImg)
}

func alphaDecode() {
	img, err := loadPNG("alpha_encode.png")
	if err != nil {
		fmt.Println(err)
	}

	rect := img.Bounds()
	newImg := image.NewRGBA(image.Rectangle{rect.Min, rect.Max})

	for y := 0; y < rect.Dx(); y++ {
		for x := 0; x < rect.Dy(); x++ {
			r1, g1, b1, a1 := img.At(x, y).RGBA()
			r2, g2, b2, a2 := uint8(r1), uint8(g1), uint8(b1), uint8(a1)
			if a2 < 255 {
				r2, g2, b2 = 255, 255, 255
			}
			newImg.SetRGBA(x, y, color.RGBA{r2, g2, b2, 255})
		}
	}
	savePNG("alpha_decode.png", newImg)
}

func lsbEncode() {
	coverImg, err := loadPNG("cover.png")
	if err != nil {
		fmt.Println(err)
	}
	sourceImg, err := loadPNG("lsb_source.png")
	if err != nil {
		fmt.Println(err)
	}

	rect := coverImg.Bounds()
	if rect != sourceImg.Bounds() {
		fmt.Println("画像サイズが一致しません")
	}

	newImg := image.NewRGBA(image.Rectangle{rect.Min, rect.Max})

	for y := 0; y < rect.Dx(); y++ {
		for x := 0; x < rect.Dy(); x++ {
			r1, g1, b1, _ := coverImg.At(x, y).RGBA()
			r2, g2, b2, _ := sourceImg.At(x, y).RGBA()
			r3 := uint8(r1)&0xfc + uint8(r2)>>6
			g3 := uint8(g1)&0xfc + uint8(g2)>>6
			b3 := uint8(b1)&0xfc + uint8(b2)>>6
			newImg.SetRGBA(x, y, color.RGBA{r3, g3, b3, 255})
		}
	}
	savePNG("lsb_encode.png", newImg)
}

func lsbDecode() {
	img, err := loadPNG("lsb_encode.png")
	if err != nil {
		fmt.Println(err)
	}

	rect := img.Bounds()
	newImg := image.NewRGBA(image.Rectangle{rect.Min, rect.Max})

	for y := 0; y < rect.Dx(); y++ {
		for x := 0; x < rect.Dy(); x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r2 := uint8(r1) << 6
			g2 := uint8(g1) << 6
			b2 := uint8(b1) << 6
			newImg.SetRGBA(x, y, color.RGBA{r2, g2, b2, 255})
		}
	}
	savePNG("lsb_decode.png", newImg)
}

func stringEncode() {
	data := []byte("Hello, World!🌍🤗🍣🍖🍕🍺")
	e := base64.StdEncoding.EncodeToString(data)
	encBytes := []byte(e)
	bitLength := uint32(len(encBytes) * 8)
	bits := make([]byte, 32+bitLength)

	// 最初の32bitにデータ長を保存する
	for i := 0; i < 32; i++ {
		bit := byte(bitLength << i >> 31)
		bits[i] = bit
	}
	// データを保存
	for i := 0; i < len(encBytes); i++ {
		byteData := encBytes[i]
		for j := 0; j < 8; j++ {
			bit := byteData << j >> 7
			bits[32+i*8+j] = bit
		}
	}

	coverImg, err := loadPNG("cover.png")
	if err != nil {
		fmt.Println(err)
	}
	rect := coverImg.Bounds()
	if (rect.Dx() * rect.Dy() * 4) < (32 + len(bits)) {
		fmt.Println("保存可能な容量を超えます")
	}

	newImg := image.NewRGBA(image.Rectangle{rect.Min, rect.Max})
	i := 0
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			r, g, b, a := coverImg.At(x, y).RGBA()
			newRed, newGreen, newBlue := uint8(r), uint8(g), uint8(b)
			if i < len(bits) {
				newRed = newRed&0xfe + bits[i]
				i++
			}
			if i < len(bits) {
				newGreen = newGreen&0xfe + bits[i]
				i++
			}
			if i < len(bits) {
				newBlue = newBlue&0xfe + bits[i]
				i++
			}
			newImg.SetRGBA(x, y, color.RGBA{newRed, newGreen, newBlue, uint8(a)})
		}
	}
	savePNG("string_encode.png", newImg)
}

func stringDecode() {
	img, err := loadPNG("string_encode.png")
	if err != nil {
		fmt.Println(err)
	}

	rect := img.Bounds()
	bits := make([]byte, rect.Dx()*rect.Dy()*4)
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			bits[y*rect.Dx()+x*3] = byte(r) & 0x01
			bits[y*rect.Dx()+x*3+1] = byte(g) & 0x01
			bits[y*rect.Dx()+x*3+2] = byte(b) & 0x01
		}
	}
	lengthBits := bits[:32]
	var length uint32
	for i := 0; i < len(lengthBits); i++ {
		length += uint32(lengthBits[i]) << (31 - i)
	}
	dataBits := bits[32 : length+32]
	dataBytes := make([]byte, length/8)
	for i := 0; i < len(dataBytes); i++ {
		var dataByte uint8
		for j := 0; j < 8; j++ {
			dataByte += dataBits[i*8+j] << (7 - j)
		}
		dataBytes[i] = dataByte
	}
	decBytes, err := base64.StdEncoding.DecodeString(string(dataBytes))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(decBytes))
}
