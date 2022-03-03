package main

import (
	"image"
	"log"

	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
)

func main() {
	// Open a test image.
	srcImage, err := imaging.Open("cleanrealimage.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	rdstImage800 := imaging.Resize(srcImage, 800, 0, imaging.Lanczos)

	// Save the resulting image as JPEG.
	err = imaging.Save(rdstImage800, "resize.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	blur := imaging.Blur(srcImage, 0.5)

	// Save the resulting image as JPEG.
	err = imaging.Save(blur, "blur.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	edgeDetector := gift.New(gift.Sobel())
	dstImage := image.NewRGBA(edgeDetector.Bounds(srcImage.Bounds()))
	edgeDetector.Draw(dstImage, srcImage)

	// Save the resulting image as JPEG.
	err = imaging.Save(dstImage, "edges.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	unsharp := gift.New(gift.UnsharpMask(1, 1, 0))
	dstImage = image.NewRGBA(unsharp.Bounds(srcImage.Bounds()))
	unsharp.Draw(dstImage, srcImage)

	// Save the resulting image as JPEG.
	err = imaging.Save(dstImage, "unsharp.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

}
