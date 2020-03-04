package main

import (
	"flag"
	"log"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/sirupsen/logrus"
)

type options struct {
	camID                  string
	videoSourceURL         string
	eventPostEndpoint      string
	eventObjectImageEnable bool
	eventSceneImageEnable  bool
	eventMaxKeypoints      int
}

type event struct {
	uuid string
}

var opt options

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	camID := flag.String("cam-id", "", "cam id used in event payloads")
	videoSourceURL := flag.String("video-source-url", "", "video feed url that will be used as source for analysis. Any source supported by OpenCV")
	eventPostEndpoint := flag.String("event-post-endpoint", "", "Target HTTP endpoint that will receive POST requests with events detected by this detector")
	eventObjectImageEnable := flag.Bool("event-object-image-enable", true, "Include detected image crop in event payload?")
	eventSceneImageEnable := flag.Bool("event-scene-image-enable", false, "Include full scene image in event payload?")
	eventMaxKeypoints := flag.Int("event-max-keypoints", -1, "Max number of keypoints in payload. Keypoints may be simplified if too large. defaults to -1 (no limit)")
	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	switch *logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	opt = options{
		camID:                  *camID,
		videoSourceURL:         *videoSourceURL,
		eventPostEndpoint:      *eventPostEndpoint,
		eventObjectImageEnable: *eventObjectImageEnable,
		eventSceneImageEnable:  *eventSceneImageEnable,
		eventMaxKeypoints:      *eventMaxKeypoints,
	}

	if opt.camID == "" {
		logrus.Errorf("'--cam-id' is required")
		os.Exit(1)
	}

	if opt.videoSourceURL == "" {
		logrus.Errorf("'--video-source-url' is required")
		os.Exit(1)
	}

	if opt.eventPostEndpoint == "" {
		logrus.Errorf("'--event-post-endpoint' is required")
		os.Exit(1)
	}

	logrus.Infof("====Starting CAM-EVENT-DETECTOR====")

	// err := initDispatcher()
	// if err != nil {
	// 	logrus.Errorf("Error initializing dispatcher. err=%s", err)
	// 	os.Exit(1)
	// }

	err1 := runDetector()
	if err1 != nil {
		logrus.Errorf("Error running detector. err=%s", err1)
		os.Exit(1)
	}

}
