package main

import (
	"fmt"
	"os"
	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\topencvgo [imgfile]")
		return
	}

	filename := os.Args[1]
	img := gocv.IMRead(filename, gocv.IMReadColor)
	window := gocv.NewWindow("Hello")

	if img.Empty() {
		fmt.Printf("Error reading image from: %v\n", filename)
		return
	}
	grayimg := gocv.NewMat()
	defer grayimg.Close()

	gocv.CvtColor(img, &grayimg, gocv.ColorBGRToGray)
	gocv.MedianBlur(grayimg, &grayimg, 5)

	maskedges := gocv.NewMat()
	defer maskedges.Close()

	gocv.AdaptiveThreshold(grayimg, &maskedges, 255, gocv.AdaptiveThresholdMean, gocv.ThresholdBinary, 9, 9)

	imgcolor := gocv.NewMat()
	defer imgcolor.Close()

	gocv.BilateralFilter(img, &imgcolor, 9, 10, 10)

	imgsticker := gocv.NewMat()
	defer imgsticker.Close()
	gocv.BitwiseAndWithMask(imgcolor, imgcolor, &imgsticker, maskedges)

	for {
		window.IMShow(imgsticker)
		if window.WaitKey(1) >= 0 {
			break
		}
	}

}


