module github.com/flaviostutz/cam-event-detector

go 1.12

require (
	github.com/chiefnoah/goalpost v0.1.1-0.20191012203108-550534edfdd1
	github.com/flaviostutz/signalutils v0.0.0-20200307152758-a500f59e21b6
	github.com/flaviostutz/sort v0.0.0-20200307103650-6c63bccb15d5
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	gocv.io/x/gocv v0.22.0
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
)

// replace github.com/flaviostutz/sort => ../sort

replace github.com/flaviostutz/signalutils => ../signalutils
