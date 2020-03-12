package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/chiefnoah/goalpost"
	"github.com/sirupsen/logrus"
)

//Define a type that implements the goalpost.Worker interface
type worker struct {
	id string
}

func (w *worker) ID() string {
	return w.id
}

func (w *worker) DoWork(ctx context.Context, job *goalpost.Job) error {
	logrus.Debugf("Cloud Upload Worker. jobID=%d", job.ID)

	ej := eventJob{}
	err := json.Unmarshal(job.Data, &ej)
	if err != nil {
		logrus.Warnf("Error unmarshaling job data. Moving job to error queue. err=%s", err)
		errorQueue.PushBytes(job.Data)
		return nil
	}

	logrus.Debugf(">>>Uploading image data...")
	reader := bytes.NewReader(ej.image)
	request, err := http.NewRequest("POST", opt.imagePostEndpoint, reader)
	if err != nil {
		logrus.Warnf("Could not prepare post to %s. retries=%d. err=%s", opt.imagePostEndpoint, job.RetryCount, err)
		return goalpost.NewRecoverableWorkerError("Error on prepare HTTP POST")
	}
	request.Header.Add("Content-Type", "image/jpeg")
	client := &http.Client{}

	logrus.Debugf("Sending HTTP POST for job %d to %s", job.ID, opt.imagePostEndpoint)
	resp, err1 := client.Do(request)
	if err1 != nil {
		logrus.Infof("Could not post to %s. failures=%d. err=%s", opt.imagePostEndpoint, job.RetryCount, err)
		return goalpost.NewRecoverableWorkerError("Error on execute HTTP POST")
	}
	if resp.StatusCode != http.StatusCreated {
		logrus.Debugf("Server returned an error. statusCode=%d", resp.StatusCode)
		return goalpost.NewRecoverableWorkerError("Server returned error")
	}
	imageLocation0, ok := resp.Header["Location"]
	if !ok || len(imageLocation0) != 1 || imageLocation0[0] == "" {
		logrus.Debugf("Server didn't return a valid Location URL for the uploaded image. location=%v", imageLocation0)
		return goalpost.NewRecoverableWorkerError("Server with no Location header")
	}
	imageLocation := imageLocation0[0]
	logrus.Debugf("Image sent to target server successfully")

	logrus.Debugf(">>>Uploading event data...")
	ej.evt.ImageLocationURL = imageLocation
	eventBytes, err1 := json.Marshal(ej.evt)
	if err1 != nil {
		logrus.Warnf("Error unmarshaling event data. Moving job to error queue. err=%s", err)
		errorQueue.PushBytes(job.Data)
		return nil
	}
	reader = bytes.NewReader(eventBytes)
	request, err = http.NewRequest("POST", opt.eventPostEndpoint, reader)
	if err != nil {
		logrus.Warnf("Could not prepare post to %s. retries=%d. err=%s", opt.eventPostEndpoint, job.RetryCount, err)
		return goalpost.NewRecoverableWorkerError("Error on prepare HTTP POST")
	}
	request.Header.Add("Content-Type", "application/json")
	client = &http.Client{}

	logrus.Debugf("Sending HTTP POST for job %d to %s", job.ID, opt.eventPostEndpoint)
	resp, err1 = client.Do(request)
	if err1 != nil {
		logrus.Infof("Could not post to %s. failures=%d. err=%s", opt.eventPostEndpoint, job.RetryCount, err)
		return goalpost.NewRecoverableWorkerError("Error on execute HTTP POST")
	}
	if resp.StatusCode != http.StatusCreated {
		logrus.Debugf("Server returned an error. statusCode=%d", resp.StatusCode)
		return goalpost.NewRecoverableWorkerError("Server returned error")
	}
	logrus.Debugf("Event sent to target server successfully")

	return nil
}

var pushQueue *goalpost.Queue
var errorQueue *goalpost.Queue

func initDispatcher() error {
	logrus.Infof("Initializing dispatcher...")

	//create events queue
	os.MkdirAll("/data/queue", os.ModePerm)
	wqueue0, err := goalpost.Init("/data/queue")
	if err != nil {
		return err
	}
	pushQueue = wqueue0

	os.MkdirAll("/data/error_queue", os.ModePerm)
	equeue0, err1 := goalpost.Init("/data/error_queue")
	if err1 != nil {
		return err1
	}
	errorQueue = equeue0

	logrus.Debugf("Registering worker...")
	w1 := &worker{
		id: "worker1",
	}
	pushQueue.RegisterWorker(w1)

	return nil
}

type eventJob struct {
	evt   event
	image []byte
}

func enqueueEvent(ev event, img *[]byte) {
	eventJob := eventJob{
		evt:   ev,
		image: *img,
	}
	ej, err := json.Marshal(eventJob)
	if err != nil {
		logrus.Warnf("Error generating event JSON. err=%s", err)
	}
	pushQueue.PushBytes(ej)
}
