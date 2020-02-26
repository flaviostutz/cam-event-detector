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
	logrus.Debugf("DoWork. jobID=%d. event=%s", job.ID, string(job.Data))

	reader := bytes.NewReader(job.Data)
	request, err := http.NewRequest("POST", opt.eventPostEndpoint, reader)
	if err != nil {
		logrus.Warnf("Could not prepare post to %s. retries=%d. err=%s", opt.eventPostEndpoint, job.RetryCount, err)
		return goalpost.NewRecoverableWorkerError("Error on prepare HTTP POST")
	}
	client := &http.Client{}

	logrus.Debugf("Sending HTTP POST for job %d to %s", job.ID, opt.eventPostEndpoint)
	resp, err1 := client.Do(request)
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

var wqueue *goalpost.Queue

func initDispatcher() error {
	logrus.Infof("Initializing dispatcher...")

	//create events queue
	os.MkdirAll("/data/queue", os.ModePerm)
	wqueue0, err := goalpost.Init("/data/queue")
	if err != nil {
		return err
	}
	wqueue = wqueue0

	logrus.Debugf("Registering worker...")
	w1 := &worker{
		id: "worker1",
	}
	wqueue.RegisterWorker(w1)

	return nil
}

func enqueueEvent(ev event) {
	eb, err := json.Marshal(ev)
	if err != nil {
		logrus.Warnf("Error generating event JSON. err=%s", err)
	}
	wqueue.PushBytes(eb)
}
