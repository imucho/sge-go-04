package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	lsbDecode()
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
