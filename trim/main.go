package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: trim [input.png]")
		fmt.Println("Usage: trim [input.png] [output.png]")
		return
	}

	ex := filepath.Ext(os.Args[1])
	if ex != ".png" {
		fmt.Println("This image is not supported.")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	minX := float64(img.Bounds().Max.X)
	minY := float64(img.Bounds().Max.Y)
	maxX := float64(-1)
	maxY := float64(-1)

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			pixel := img.At(x, y)
			_, _, _, a := pixel.RGBA()

			if a > 0 {
				minX = math.Min(float64(minX), float64(x))
				maxX = math.Max(float64(maxX), float64(x))
				minY = math.Min(float64(minY), float64(y))
				maxY = math.Max(float64(maxY), float64(y))
			}
		}
	}

	if maxX == -1 {
		fmt.Println("This is a fully transparent image.")
		return;
	}

	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	outImg := img.(SubImager).SubImage(image.Rect(int(minX), int(minY), int(maxX) + 1, int(maxY) + 1))
	outImgFileName := os.Args[1]
	if len(os.Args) >= 3 {
		outImgFileName = os.Args[2]
	}
	out, err := os.Create(outImgFileName)
	if err != nil {
		panic(err)
	}

	png.Encode(out, outImg)
}