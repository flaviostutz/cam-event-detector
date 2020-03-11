package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/flaviostutz/sort"
	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

// var tracking sort.SORT
var window *gocv.Window
var window2 *gocv.Window
var tracking sort.SORT
var frameCounter = 0

var curImg gocv.Mat
var prevImg gocv.Mat
var flow gocv.Mat

func initTracker() {
	window = gocv.NewWindow("Tracker")
	window.ResizeWindow(600, 600)
	tracking = sort.NewSORT(2, 2, 0.2)
	logrus.SetLevel(logrus.InfoLevel)
	// curImg = gocv.NewMat()
	// prevImg = gocv.NewMat()
	// window2 = gocv.NewWindow("Tracker 2")
	// window2.ResizeWindow(600, 600)
}

func trackFrame(img *gocv.Mat, bboxes [][]float64) {

	// show debugging
	for _, bbox := range bboxes {
		gocv.Rectangle(img, bboxToRect(bbox), color.RGBA{254, 0, 254, 254}, 1)
	}

	tracking.Update(bboxes)

	// gocv.CvtColor(*img, &curImg, gocv.ColorBGRToGray)
	// if frameCounter > 0 {
	// 	opticalFlowDense(window2, curImg, prevImg)
	// }
	// curImg.CopyTo(&prevImg)

	//show debugging
	gocv.PutText(img, "trk-CurrentPrediction()", image.Point{int(5), int(30)}, gocv.FontHersheyPlain, 1.3, color.RGBA{254, 0, 0, 254}, 2)
	gocv.PutText(img, "trk-CurrentState()", image.Point{int(5), int(45)}, gocv.FontHersheyPlain, 1.3, color.RGBA{0, 254, 0, 254}, 3)
	gocv.PutText(img, "trk-LastBBox", image.Point{int(5), int(60)}, gocv.FontHersheyPlain, 1.3, color.RGBA{0, 0, 180, 254}, 2)
	gocv.PutText(img, "detection", image.Point{int(5), int(75)}, gocv.FontHersheyPlain, 1.3, color.RGBA{180, 0, 180, 254}, 2)
	gocv.PutText(img, "trk-LastBBoxIOU", image.Point{int(5), int(90)}, gocv.FontHersheyPlain, 1.3, color.RGBA{254, 254, 0, 254}, 2)
	for _, trk := range tracking.Trackers {
		if trk.Updates > 0 {
			currentPrediction := trk.CurrentPrediction()
			currentState := trk.CurrentState()
			lastBBox := trk.LastBBox
			lastBBox = []float64{lastBBox[0] + 2, lastBBox[1] + 2, lastBBox[2] + 2, lastBBox[3] + 2}
			currentState = []float64{currentState[0] + 4, currentState[1] + 4, currentState[2] + 4, currentState[3] + 4}
			lastBBoxIOU := trk.LastBBoxIOU
			gocv.PutText(img, fmt.Sprintf("%d upt:%d uwp:%d psu:%d", trk.ID, trk.Updates, trk.UpdatesWithoutPredict, trk.PredictsSinceUpdate), image.Point{int(lastBBox[0]), int(lastBBox[1])}, gocv.FontHersheyPlain, 1.2, color.RGBA{254, 0, 0, 254}, 1)
			gocv.Rectangle(img, bboxToRect(currentPrediction), color.RGBA{254, 0, 0, 254}, 1)
			gocv.Rectangle(img, bboxToRect(currentState), color.RGBA{0, 254, 0, 254}, 1)
			gocv.Rectangle(img, bboxToRect(lastBBox), color.RGBA{0, 0, 254, 254}, 2)
			if len(lastBBoxIOU) > 0 {
				lastBBoxIOU = []float64{lastBBoxIOU[0] + 4, lastBBoxIOU[1] + 4, lastBBoxIOU[2] + 4, lastBBoxIOU[3] + 4}
				gocv.Rectangle(img, bboxToRect(lastBBoxIOU), color.RGBA{254, 254, 0, 254}, 2)
			}
		}
	}

	window.IMShow(*img)
	window.WaitKey(1000 / 5)
	fmt.Printf("%d ", frameCounter)
	frameCounter = frameCounter + 1
}

func bboxToRect(bbox []float64) image.Rectangle {
	return gocv.BoundingRect([]image.Point{image.Point{int(bbox[0]), int(bbox[1])}, image.Point{int(bbox[2]), int(bbox[1])}, image.Point{int(bbox[2]), int(bbox[1] - (bbox[3] - bbox[1]))}, image.Point{int(bbox[0]), int(bbox[1] - (bbox[3] - bbox[1]))}})
}
