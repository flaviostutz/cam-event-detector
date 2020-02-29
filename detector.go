package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

const minimumArea = 200

func runDetector() error {
	logrus.Infof("Opening source video feed...")
	feed, err := gocv.OpenVideoCapture(opt.videoSourceURL)
	if err != nil {
		return fmt.Errorf("Error opening stream. source=%s. err=%s", opt.videoSourceURL, err)
	}
	defer feed.Close()
	logrus.Debugf("Feed opened")

	window := gocv.NewWindow("Cam")
	defer window.Close()

	// window1 := gocv.NewWindow("Debug1")
	// defer window1.Close()
	// window2 := gocv.NewWindow("Debug2")
	// defer window2.Close()

	img := gocv.NewMat()
	defer img.Close()

	imgSmall := gocv.NewMat()
	defer imgSmall.Close()

	imgGrey := gocv.NewMat()
	defer imgGrey.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()

	status := "Ready"
	statusColor := color.RGBA{0, 255, 0, 0}

	scale := 1.0

	// time.Sleep(5 * time.Second)
	logrus.Infof("Starting detections...")
	for {
		//+20% - 22%
		// logrus.Debugf("read")
		if ok := feed.Read(&img); !ok {
			break
		}
		if img.Empty() {
			time.Sleep(400 * time.Millisecond)
			continue
		}

		// time.Sleep(50 * time.Millisecond)

		//+3% - 50% less CPU after
		gocv.Resize(img, &imgSmall, image.Pt(0, 0), scale, scale, gocv.InterpolationNearestNeighbor)

		//+4% - 27%
		gocv.CvtColor(imgSmall, &imgGrey, gocv.ColorBGRAToGray)

		//+33% - 61% - grey       - //+37% - 65% - color
		mog2.Apply(imgGrey, &imgDelta)
		// differentialCollins(imgs, &imgDelta)

		//+3% - 64% - grey
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)
		// window1.IMShow(imgDelta)

		//+1% - 65% - grey
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Dilate(imgThresh, &imgThresh, kernel)
		// window2.IMShow(imgThresh)

		status = "Ready"
		statusColor = color.RGBA{0, 255, 0, 0}

		// // now find contours
		//+18% - 83% - grey
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for _, c := range contours {
			area := gocv.ContourArea(c)
			if area < minimumArea {
				continue
			}

			status = "Motion detected"
			statusColor = color.RGBA{255, 0, 0, 0}
			// gocv.DrawContours(&img, contours, i, statusColor, 2)

			rect := gocv.BoundingRect(c)

			//scale bounding rect to input image sizes
			rscale := 1 / scale
			rx := int(float64(rect.Bounds().Min.X) * rscale)
			ry := int(float64(rect.Bounds().Max.Y) * rscale)
			rw := int(float64(rect.Dx()) * rscale)
			rh := int(float64(rect.Dy()) * rscale)
			rects := gocv.BoundingRect([]image.Point{image.Point{rx, ry}, image.Point{rx + rw, ry}, image.Point{rx + rw, ry - rh}, image.Point{rx, ry - rh}})

			//scale bounds to input image
			gocv.Rectangle(&img, rects, color.RGBA{0, 0, 255, 0}, 2)
		}

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		//83% - grey
		//89% - color

		//+16% - 109% - grey
		window.IMShow(img)
		// window.IMShow(imgDelta)

		if window.WaitKey(20) == 27 {
			break
		}
	}

	return nil
}

// func differentialCollins(imgs []gocv.Mat, imgDelta *gocv.Mat) {
// 	hd1 := gocv.NewMatWithSize(imgs[0].Rows(), imgs[0].Cols(), imgs[0].Type())
// 	hd2 := hd1.Clone()
// 	gocv.Subtract(imgs[0], imgs[2], &hd1)
// 	gocv.Subtract(imgs[1], imgs[2], &hd2)
// 	gocv.BitwiseAnd(hd1, hd2, imgDelta)
// }
