package main

import (
	"gocv.io/x/gocv"
)

// OPTICAL FLOW TESTS
//https://opencv-python-tutroals.readthedocs.io/en/latest/py_tutorials/py_video/py_lucas_kanade/py_lucas_kanade.html
func opticalFlowDense(optWindow *gocv.Window, curImg gocv.Mat, prevImg gocv.Mat) {
	flow := gocv.NewMat()
	// prevImg = curImg.Clone()
	icols := curImg.Cols()
	irows := curImg.Rows()
	gocv.CalcOpticalFlowFarneback(prevImg, curImg, &flow, 0.5, 3, 15, 3, 5, 1.2, 0)
	m0 := gocv.NewMatWithSize(irows, icols, gocv.MatTypeCV32FC1)
	m1 := gocv.NewMatWithSize(irows, icols, gocv.MatTypeCV32FC1)
	for i := 0; i < flow.Rows(); i++ {
		for j := 0; j < flow.Cols(); j++ {
			v := flow.GetVecfAt(i, j)
			m0.SetFloatAt(i, j, v[0])
			m1.SetFloatAt(i, j, v[1])
		}
	}
	mmag := gocv.NewMatWithSize(irows, icols, gocv.MatTypeCV32FC1)
	mang := gocv.NewMatWithSize(irows, icols, gocv.MatTypeCV32FC1)
	gocv.CartToPolar(m0, m1, &mmag, &mang, true)

	// mmagn := gocv.NewMat()
	// gocv.Normalize(mmag, &mmagn, 0, 255, gocv.NormMinMax)

	mh := gocv.NewMatWithSize(irows, icols, gocv.MatTypeCV8UC1)
	ms := gocv.NewMatWithSize(irows, icols, gocv.MatTypeCV8UC1)
	mv := gocv.NewMatWithSize(irows, icols, gocv.MatTypeCV8UC1)
	for i := 0; i < flow.Rows(); i++ {
		for j := 0; j < flow.Cols(); j++ {

			hv := (mang.GetFloatAt(i, j) / 360.0) * 255.0
			mh.SetUCharAt(i, j, uint8(hv))

			ms.SetUCharAt(i, j, 255)

			vv := (mmag.GetFloatAt(i, j) / 30) * 255
			mv.SetUCharAt(i, j, uint8(vv))
		}
	}

	flowImg := gocv.NewMatWithSize(irows, icols, gocv.MatTypeCV8UC3)
	gocv.Merge([]gocv.Mat{mh, ms, mv}, &flowImg)

	flowImg2 := gocv.NewMat()
	gocv.CvtColor(flowImg, &flowImg2, gocv.ColorHSVToRGB)
	optWindow.IMShow(flowImg2)
}
