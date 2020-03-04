package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/flaviostutz/sort"
	"gocv.io/x/gocv"
)

// var tracking sort.SORT
var window *gocv.Window
var tracking sort.SORT
var frameCounter = 0

func initTracker() {
	window = gocv.NewWindow("Cam")
	tracking = sort.NewSORT(20, 3, 0.3)
}

func trackFrame(img *gocv.Mat, bboxes [][]float64) {

	// tracking.Update()

	//show debugging
	for _, bbox := range bboxes {
		gocv.Rectangle(img, bboxToRect(bbox), color.RGBA{0, 0, 254, 254}, 1)
		// gocv.DrawContours(img, [][]image.Point{rrect.Contour}, 0, color.RGBA{254, 0, 0, 254}, 2)
	}
	window.IMShow(*img)
	window.WaitKey(350)
	fmt.Printf("%d ", frameCounter)
	frameCounter = frameCounter + 1
}

func bboxToRect(bbox []float64) image.Rectangle {
	return gocv.BoundingRect([]image.Point{image.Point{int(bbox[0]), int(bbox[1])}, image.Point{int(bbox[2]), int(bbox[1])}, image.Point{int(bbox[2]), int(bbox[1] - (bbox[3] - bbox[1]))}, image.Point{int(bbox[0]), int(bbox[1] - (bbox[3] - bbox[1]))}})
}
