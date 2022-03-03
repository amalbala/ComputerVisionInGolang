package main

import (
	"fmt"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

func main() {
	img, err := imgio.Open("cleanrealimage.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}

	resized := transform.Resize(img, 800, 800, transform.Linear)

	if err := imgio.Save("resized.png", resized, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}

	blur := blur.Gaussian(img, 3.0)

	if err := imgio.Save("blur.png", blur, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}

	edgeDetec := effect.EdgeDetection(img, 1.0)

	if err := imgio.Save("edge_detection.png", edgeDetec, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}

	unsharp := effect.UnsharpMask(img, 0.6, 1.2)

	if err := imgio.Save("unsharp_mask.jpg", unsharp, imgio.JPEGEncoder(95)); err != nil {
		fmt.Println(err)
		return
	}
}
