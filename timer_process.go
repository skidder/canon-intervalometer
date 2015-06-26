package main

import (
	"log"
	"time"
)

type TimerProcess struct {
	cameraCommandChannel chan CameraCommand
	running              bool
	stopChannel          chan bool
}

func NewTimerProcess(cameraCommandChannel chan CameraCommand) *TimerProcess {
	return &TimerProcess{cameraCommandChannel: cameraCommandChannel, stopChannel: make(chan bool)}
}

func (t *TimerProcess) Run(sleepInterval int) {
	t.running = true

	t.cameraCommandChannel <- TAKE_PICTURE
	for {
		select {
		case t.running = <-t.stopChannel:
			log.Println("Timer Process returning")
			return
		case <-time.After(time.Duration(sleepInterval) * time.Second):
			t.cameraCommandChannel <- TAKE_PICTURE
		}
	}
}

func (t *TimerProcess) Stop() {
	t.stopChannel <- true
}

func (t *TimerProcess) IsRunning() bool {
	return t.running
}
