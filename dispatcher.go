package main

import (
	"context"
	"fmt"
	"os"
	"time"

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
	fmt.Printf("Hello, %s\n", job.Data)
	if job.RetryCount < 9 { //try 10 times
		return goalpost.NewRecoverableWorkerError("Something broke, try again")
	}

	//Something *really* broke, don't retry
	//return errors.New("Something broke, badly")
	return nil
}

func runDispatcher() {
	logrus.Infof("Preparing events queue...")

	//create events queue
	os.MkdirAll("/data/queue", os.ModePerm)
	wqueue, err := goalpost.Init("/data/queue")
	if err != nil {
		logrus.Errorf("Error initializing queue engine. err=%s", err)
		os.Exit(1)
	}
	defer wqueue.Close()

	//register worker
	w1 := &worker{
		id: "worker1",
	}
	wqueue.RegisterWorker(w1)

	wqueue.PushBytes([]byte("World1"))
	time.Sleep(5 * time.Second)
	wqueue.PushBytes([]byte("World2"))
	time.Sleep(5 * time.Second)
	wqueue.PushBytes([]byte("World3"))
	time.Sleep(5 * time.Second)
	wqueue.PushBytes([]byte("World4"))
	time.Sleep(5 * time.Second)
}
