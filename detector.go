package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

func runDetector() error {
	logrus.Infof("Opening source video feed...")
	feed, err := gocv.OpenVideoCapture(opt.videoSourceURL)
	if err != nil {
		return fmt.Errorf("Error opening stream. source=%s. err=%s", opt.videoSourceURL, err)
	}

	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	logrus.Infof("Starting detections...")
	for {
		feed.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}
}
