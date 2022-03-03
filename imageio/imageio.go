package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type ImageIO struct {
	filename     string
	extension    string
	imageRGB     *image.RGBA
	alphaChannel *image.Alpha
	embeddedmap  *color.Palette
}

func splitPalettedImage(img image.Image) (image.Image, color.Palette) {
	fmt.Println("embedded image")
	imagePalleted := img.ColorModel()
	imageRGB := image.NewRGBA(img.Bounds())

	bounds := imageRGB.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := imageRGB.At(x, y).RGBA()
			imageRGB.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
		}
	}
	return imageRGB, imagePalleted.(color.Palette)
}

func mergeRGBA(imgRGB image.Image, imgAlpha image.Image) *image.RGBA {

	imageRGBA := image.NewRGBA(imgRGB.Bounds())
	bounds := imgRGB.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := imgRGB.At(x, y).RGBA()
			_, _, _, a := imgAlpha.At(x, y).RGBA()

			imageRGBA.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}

	return imageRGBA
}

func splitRGBA(img image.Image) (*image.RGBA, *image.Alpha) {
	fmt.Println("transparent image")
	imageRGB := image.NewRGBA(img.Bounds())
	imageAlphaChannel := image.NewAlpha(img.Bounds())
	bounds := imageRGB.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			imageRGB.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
			imageAlphaChannel.Set(x, y, color.Alpha{A: uint8(a)})
		}
	}

	return imageRGB, imageAlphaChannel
}

func YCbCrToRGB(img image.Image) *image.RGBA {
	var (
		bounds  = img.Bounds()
		imgRGBA = image.NewRGBA(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			imgRGBA.Set(x, y, color.RGBAModel.Convert(img.At(x, y)))
		}
	}
	return imgRGBA
}

func openImage(filename string) *ImageIO {
	img := ImageIO{extension: filepath.Ext(filename),
		filename: strings.TrimSuffix(filename, filepath.Ext(filename))}

	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	imageData, _, err := image.Decode(reader)

	fmt.Println("Type: ", reflect.TypeOf(imageData))
	fmt.Println("fileame: ", img.filename, " ext: ", img.extension)

	switch imageData.ColorModel() {
	case color.NRGBAModel:
		fmt.Println("is alpha image")
		img.imageRGB, img.alphaChannel = splitRGBA(imageData)
	case color.YCbCrModel:
		img.imageRGB = YCbCrToRGB(imageData)
	}

	return &img
}

func stringColorModel(colormodel color.Model) string {
	switch colormodel {
	case color.RGBAModel:
		return "RGBA"
	}
	return "undefined"
}

func saveImage(filename string, img *ImageIO) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	imgToEncode := img.imageRGB

	if img.alphaChannel != nil {
		imgToEncode = mergeRGBA(img.imageRGB, img.alphaChannel)
	}

	switch img.extension {

	case ".png":
		err = png.Encode(file, imgToEncode)
		if err != nil {
			log.Fatal(err)
		}
	default:
		opt := jpeg.Options{
			Quality: 90,
		}
		err = jpeg.Encode(file, imgToEncode, &opt)
		if err != nil {
			log.Fatal(err)
		}

	}

}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\timageio [imgfile]")
		return
	}

	filename := os.Args[1]
	outfilename := os.Args[2]
	tempImage := openImage(filename)

	fmt.Println("Open: ", tempImage.filename)
	fmt.Println("Type:", stringColorModel(tempImage.imageRGB.ColorModel()))

	saveImage(outfilename, tempImage)

}
