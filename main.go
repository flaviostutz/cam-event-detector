package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	flag.Parse()

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

	opt := handlers.Options{
		MaxZoomLevel: *maxZoomLevel,
	}

	if opt.WFSURL == "" {
		logrus.Errorf("'--wfs-url' is required")
		os.Exit(1)
	}

	logrus.Infof("====Starting CAM-EVENT-DETECTOR====")
	h := handlers.NewHTTPServer(opt)
	if err != nil {
		logrus.Errorf("Error starting server. err=%s", err)
		os.Exit(1)
	}

	webcam, err := gocv.OpenVideoCapture(os.Getenv("RTSPLINK"))
	if err != nil {
		panic("Error in opening webcam: " + err.Error())
	}

}
