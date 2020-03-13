package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/flaviostutz/signalutils"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

const minimumDimension = 5

const awarenessNormal = "Awareness-RegularMovement"
const awarenessHigh = "Awareness-HighMovement"
const awarenessLow = "Awareness-LowMovement"

const movementIdle = "Movement-None"
const movementDetected = "Movement-Detected"

type stateData struct {
	level     float64
	diff      float64
	imageData *gocv.Mat
	imageTime time.Time
}

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
	// imgHSL := gocv.NewMat()
	imgOne := gocv.NewMat()
	imgDelta := gocv.NewMat()
	imgThresh := gocv.NewMat()
	// mog2 := gocv.NewBackgroundSubtractorKNNWithParams(4, 6, false)
	mog2 := gocv.NewBackgroundSubtractorMOG2WithParams(4, 6, false)

	//AWARENESS EVENTS
	dynamicSchimttTrigger, err := signalutils.NewDynamicSchmittTriggerTimeWindow(30*time.Second, 60, 10, 2.0, 0.5, false)
	if err != nil {
		return err
	}
	diffAverage := signalutils.NewMovingAverageTimeWindow(1*time.Second, 10)
	awarenessState := signalutils.NewStateTracker("", 4, onAwarenessChanged, 10*time.Second, onAwarenessUnchanged)

	//MOVEMENT LEVEL EVENTS
	// movementAveragerBase := signalutils.NewMovingAverageTimeWindow(30*time.Second, 60)
	movementAverager := signalutils.NewMovingAverageTimeWindow(10*time.Second, 20)
	movementState := signalutils.NewStateTracker("", 3, onMovementChanged, 10*time.Second, onMovementUnchanged)

	// initTracker()

	scale := 1.0

	window1 := gocv.NewWindow("Detector")
	window1.ResizeWindow(600, 600)

	window2 := gocv.NewWindow("Cam")
	window2.ResizeWindow(600, 600)

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
		// gocv.CvtColor(img, &imgHSL, gocv.ColorBGRToHLS)
		gocv.CvtColor(img, &imgOne, gocv.ColorBGRToGray)

		// imgOne = gocv.Split(imgHSL)[2]

		// gocv.Blur(imgOne, &imgOne, image.Pt(10, 10))

		//+33% - 61% - grey       - //+37% - 65% - color
		mog2.Apply(imgOne, &imgDelta)
		// differentialCollins(imgs, &imgDelta)
		// window.WaitKey(1000000)

		//+3% - 64% - grey
		gocv.Threshold(imgDelta, &imgThresh, 10, 255, gocv.ThresholdBinary)

		//+1% - 65% - grey
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Erode(imgThresh, &imgThresh, kernel)
		gocv.Dilate(imgThresh, &imgThresh, kernel)
		gocv.Dilate(imgThresh, &imgThresh, kernel)
		gocv.Dilate(imgThresh, &imgThresh, kernel)
		gocv.Dilate(imgThresh, &imgThresh, kernel)
		// gocv.Erode(imgThresh, &imgThresh, kernel)
		gocv.Erode(imgThresh, &imgThresh, kernel)
		gocv.Erode(imgThresh, &imgThresh, kernel)
		gocv.Erode(imgThresh, &imgThresh, kernel)
		gocv.Erode(imgThresh, &imgThresh, kernel)
		window1.IMShow(imgThresh)

		// now find contours and corresponding bbox
		//+18% - 83% - grey
		bboxes := make([][]float64, 0)
		totalContoursArea := 0.0
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for _, c := range contours {
			// area := gocv.ContourArea(c)
			// if area < minimumArea {
			// 	continue
			// }

			totalContoursArea = totalContoursArea + gocv.ContourArea(c)

			rect := gocv.BoundingRect(c)
			if rect.Dx() < minimumDimension || rect.Dy() < minimumDimension {
				continue
			}
			// rotatedRect := gocv.MinAreaRect(c)

			//scale bounding rect to input image sizes
			rscale := 1 / scale
			rx := float64(rect.Bounds().Min.X) * rscale
			ry := float64(rect.Bounds().Max.Y) * rscale
			rw := float64(rect.Dx()) * rscale
			rh := float64(rect.Dy()) * rscale
			bboxes = append(bboxes, []float64{rx, ry, rx + rw, ry + rh})
		}

		//MOVEMENT AVERAGER
		ok := movementAverager.AddSampleIfNearAverage(totalContoursArea, 3.0)
		if !ok {
			logrus.Debugf("Sample skipped1")
		}
		// ok := movementAveragerBase.AddSampleIfNearAverage(totalContoursArea, 3.0)
		// if !ok {
		// 	logrus.Debugf("Sample skipped2")
		// }

		// gocv.PutText(&imgThresh, fmt.Sprintf("%d", int(movementAverager.Average())), image.Point{int(20), int(100)}, gocv.FontHersheyPlain, 1.3, color.RGBA{254, 254, 254, 254}, 2)
		// window1.IMShow(imgThresh)

		//SCHMITT TRIGGER OVER MOVIMENT
		ok, diff0 := dynamicSchimttTrigger.SetCurrentValue(totalContoursArea)
		if !ok {
			logrus.Debugf("Sample skipped3")
		}
		// fmt.Printf("CD %.1f\n", diff0)
		diffAverage.AddSample(diff0)
		diff := diffAverage.Average()
		movementLevel := movementAverager.Average()
		// diff := movementAverager.Average() - movementAveragerBase.Average()
		ci := img.Clone()
		data := &stateData{
			level:     movementLevel,
			diff:      diff,
			imageData: &ci,
			imageTime: time.Now(),
		}
		if diff > 5000.0 {
			awarenessState.SetTransientStateWithData(awarenessHigh, diff, data)
		} else if diff < -5000.0 {
			awarenessState.SetTransientStateWithData(awarenessLow, diff, data)
		} else {
			awarenessState.SetTransientStateWithData(awarenessNormal, diff, data)
		}

		if movementLevel > 10 {
			movementState.SetTransientStateWithData(movementDetected, movementLevel, data)
		} else {
			movementState.SetTransientStateWithData(movementIdle, movementLevel, data)
		}

		status := fmt.Sprintf("  Normal level=%.1f", diff)
		statusColor := color.RGBA{0, 150, 0, 254}
		if awarenessState.CurrentState.Name == awarenessHigh {
			status = fmt.Sprintf("Attention HIGH level=%.1f", diff)
			statusColor = color.RGBA{254, 0, 0, 254}
		} else if awarenessState.CurrentState.Name == awarenessLow {
			status = fmt.Sprintf("Attention LOW level=%.1f", diff)
			statusColor = color.RGBA{0, 0, 254, 254}
		}
		gocv.PutText(&img, status, image.Point{int(20), int(50)}, gocv.FontHersheyPlain, 1.3, statusColor, 2)

		l, u := dynamicSchimttTrigger.GetLowerUpperLimits()
		gocv.PutText(&img, fmt.Sprintf("%d", int64(u)), image.Point{int(20), int(100)}, gocv.FontHersheyPlain, 1.0, color.RGBA{100, 0, 0, 254}, 1)
		gocv.PutText(&img, fmt.Sprintf("%d", int64(totalContoursArea)), image.Point{int(20), int(112)}, gocv.FontHersheyPlain, 1.0, color.RGBA{0, 100, 0, 254}, 1)
		gocv.PutText(&img, fmt.Sprintf("%d", int64(l)), image.Point{int(20), int(124)}, gocv.FontHersheyPlain, 1.0, color.RGBA{0, 0, 100, 254}, 1)

		level := 0
		if movementState.CurrentState.Data != nil {
			ed := movementState.CurrentState.Data.(*stateData)
			level = int(ed.level)
		}
		gocv.PutText(&img, fmt.Sprintf("%s %d", movementState.CurrentState.Name, level), image.Point{int(20), int(140)}, gocv.FontHersheyPlain, 1.3, color.RGBA{100, 0, 100, 254}, 1)

		window2.IMShow(img)
		// logrus.Debugf(status)

		// logrus.Debugf("Movement level: %f", movementAverager.Average())

		window1.WaitKey(1000 / 4)
		window2.WaitKey(1000 / 4)
		//83% - grey
		//89% - color

		//+16% - 109% - grey
		// trackFrame(&img, bboxes)
	}
	return nil
}

func onAwarenessChanged(prevState *signalutils.State, state *signalutils.State) {
	if state.Data == nil {
		return
	}
	// prevData := (prevState.Data).(*eventData)
	fmt.Printf("AWARENESS CHANGED: prevState=%v\n", prevState.Name)
	// diff := d[0].(float64)
	// img := d[2].(*gocv.Mat)
	if prevState != nil {
		ev, imgBytes := newEventReport(prevState)
		enqueueEventReport(ev, &imgBytes)
	} else {
		ev, imgBytes := newEventReport(state)
		enqueueEventReport(ev, &imgBytes)
	}
}
func onAwarenessUnchanged(state *signalutils.State) {
	if state.Data == nil {
		return
	}
	// d := (state.Data).(*eventData)
	fmt.Printf("AWARENESS UNCHANGED: %s\n", state.Name)
	// diff := d[0].(float64)
	// img := d[2].(*gocv.Mat)
	ev, imgBytes := newEventReport(state)
	enqueueEventReport(ev, &imgBytes)
}

func onMovementChanged(prevState *signalutils.State, state *signalutils.State) {
	if state.Data == nil {
		return
	}
	// d := (state.Data).(*eventData)
	fmt.Printf("MOVEMENT CHANGED: %s\n", state.Name)
	// level := d[0].(float64)
	// img := d[1].(*gocv.Mat)
	ev, imgBytes := newEventReport(state)
	enqueueEventReport(ev, &imgBytes)
}
func onMovementUnchanged(state *signalutils.State) {
	if state.Data == nil {
		return
	}
	// d := (state.Data).(*eventData)
	fmt.Printf("MOVEMENT UNCHANGED: %s\n", state.Name)
	// level := d[0].(float64)
	// img := d[1].(*gocv.Mat)
	ev, imgBytes := newEventReport(state)
	enqueueEventReport(ev, &imgBytes)
}

func newEventReport(state *signalutils.State) (*eventReport, []byte) {
	if state.Data == nil {
		return nil, []byte{}
	}

	evtData := state.Data.(*stateData)
	imageData := evtData.imageData
	imageTime := evtData.imageTime
	level := state.Level
	if state.HighestData != nil {
		highestData := state.HighestData.(*stateData)
		imageData = highestData.imageData
		imageTime = highestData.imageTime
		level = state.HighestLevel
	}

	imgBytes, err := gocv.IMEncode(gocv.JPEGFileExt, *imageData)
	if err != nil {
		logrus.Warn("Error encoding image for event. err=%s", err)
	}
	return &eventReport{
		UUID:      uuid.NewV4().String(),
		EventType: state.Name,
		CamID:     opt.camID,
		Timestamp: time.Now(),
		Level:     level,
		Start:     &state.Start,
		Stop:      state.Stop,
		Active:    (state.Stop == nil),
		ImageTime: &imageTime,
	}, imgBytes

}

// func differentialCollins(imgs []gocv.Mat, imgDelta *gocv.Mat) {
// 	hd1 := gocv.NewMatWithSize(imgs[0].Rows(), imgs[0].Cols(), imgs[0].Type())
// 	hd2 := hd1.Clone()
// 	gocv.Subtract(imgs[0], imgs[2], &hd1)
// 	gocv.Subtract(imgs[1], imgs[2], &hd2)
// 	gocv.BitwiseAnd(hd1, hd2, imgDelta)
// }
