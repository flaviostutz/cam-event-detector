package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

func runDetector() error {
	logrus.Infof("Opening source video feed...")
	feed, err := gocv.OpenVideoCapture(opt.videoSourceURL)
	if err != nil {
		return fmt.Errorf("Error opening stream. source=%s. err=%s", opt.videoSourceURL, err)
	}
	logrus.Debugf("Feed opened")

	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	time.Sleep(5 * time.Second)
	logrus.Infof("Starting detections...")
	for {
		feed.Read(&img)
		window.IMShow(img)

		//testing
		// evt := event{
		// 	uuid: "test",
		// }
		// enqueueEvent(evt)
		// time.Sleep(10 * time.Second)
	}

}
