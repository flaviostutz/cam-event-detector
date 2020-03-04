package main

import (
	"fmt"
	"image"
	"time"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

const minimumArea = 200

func runDetector() error {
	logrus.Infof("Opening video feed...")
	feed, err := gocv.OpenVideoCapture(opt.videoSourceURL)
	if err != nil {
		return fmt.Errorf("Error opening stream. source=%s. err=%s", opt.videoSourceURL, err)
	}
	defer feed.Close()
	logrus.Debugf("Feed opened")

	img := gocv.NewMat()
	// imgSmall := gocv.NewMat()
	imgGrey := gocv.NewMat()
	imgDelta := gocv.NewMat()
	imgThresh := gocv.NewMat()
	mog2 := gocv.NewBackgroundSubtractorMOG2()

	initTracker()

	scale := 1.0

	// time.Sleep(5 * time.Second)
	logrus.Infof("Starting detections...")
	for {
		//+20% - 22%
		if ok := feed.Read(&img); !ok {
			break
		}
		if img.Empty() {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		//+3% - 50% less CPU after
		// gocv.Resize(img, &imgSmall, image.Pt(0, 0), scale, scale, gocv.InterpolationNearestNeighbor)

		//+4% - 27%
		gocv.CvtColor(img, &imgGrey, gocv.ColorBGRAToGray)

		//+33% - 61% - grey       - //+37% - 65% - color
		mog2.Apply(imgGrey, &imgDelta)
		// differentialCollins(imgs, &imgDelta)

		//+3% - 64% - grey
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)

		//+1% - 65% - grey
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Dilate(imgThresh, &imgThresh, kernel)

		// now find contours and corresponding bbox
		//+18% - 83% - grey
		bboxes := make([][]float64, 0)
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for _, c := range contours {
			area := gocv.ContourArea(c)
			if area < minimumArea {
				continue
			}

			rect := gocv.BoundingRect(c)
			// rotatedRect := gocv.MinAreaRect(c)

			//scale bounding rect to input image sizes
			rscale := 1 / scale
			rx := float64(rect.Bounds().Min.X) * rscale
			ry := float64(rect.Bounds().Max.Y) * rscale
			rw := float64(rect.Dx()) * rscale
			rh := float64(rect.Dy()) * rscale
			bboxes = append(bboxes, []float64{rx, ry, rx + rw, ry + rh})
		}

		//83% - grey
		//89% - color

		//+16% - 109% - grey
		trackFrame(&img, bboxes)
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
