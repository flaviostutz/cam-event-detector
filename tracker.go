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
	tracking = sort.NewSORT(2, 2, 0.2)
}

func trackFrame(img *gocv.Mat, bboxes [][]float64) {

	//show debugging
	for _, bbox := range bboxes {
		gocv.Rectangle(img, bboxToRect(bbox), color.RGBA{0, 0, 254, 254}, 2)
	}

	tracking.Update(bboxes)

	//show debugging
	for _, trk := range tracking.Trackers {
		bbox := trk.LastBBox
		bbox1 := trk.CurrentState()
		bbox2 := trk.CurrentPrediction()
		gocv.PutText(img, fmt.Sprintf("%d upt:%d psu:%d uwp:%d", trk.ID, trk.Updates, trk.PredictsSinceUpdate, trk.UpdatesWithoutPredict), image.Point{int(bbox[0]), int(bbox[1])}, gocv.FontHersheyPlain, 1.2, color.RGBA{254, 0, 0, 254}, 1)
		gocv.Rectangle(img, bboxToRect(bbox), color.RGBA{0, 254, 0, 254}, 1)
		gocv.Rectangle(img, bboxToRect(bbox1), color.RGBA{0, 0, 254, 254}, 1)
		gocv.Rectangle(img, bboxToRect(bbox2), color.RGBA{254, 0, 0, 254}, 1)
	}

	window.IMShow(*img)
	window.WaitKey(1000 / 2)
	fmt.Printf("%d ", frameCounter)
	frameCounter = frameCounter + 1
}

func bboxToRect(bbox []float64) image.Rectangle {
	return gocv.BoundingRect([]image.Point{image.Point{int(bbox[0]), int(bbox[1])}, image.Point{int(bbox[2]), int(bbox[1])}, image.Point{int(bbox[2]), int(bbox[1] - (bbox[3] - bbox[1]))}, image.Point{int(bbox[0]), int(bbox[1] - (bbox[3] - bbox[1]))}})
}
