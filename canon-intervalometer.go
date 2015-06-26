package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/qml.v1"
)

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	engine := qml.NewEngine()

	ci := NewCanonIntervalometer(engine)
	ci.Initialize()
	return ci.Run()
}

type CanonIntervalometer struct {
	engine        *qml.Engine
	window        *qml.Window
	cameraManager *CameraManager
	timerProcess  *TimerProcess
}

func NewCanonIntervalometer(e *qml.Engine) *CanonIntervalometer {
	return &CanonIntervalometer{engine: e}
}

func (c *CanonIntervalometer) Initialize() error {
	context := c.engine.Context()
	context.SetVar("intervalometer", c)

	c.cameraManager = NewCameraManager()
	c.timerProcess = NewTimerProcess(c.cameraManager.GetCommandChannel())
	return nil
}

func (c *CanonIntervalometer) Run() error {
	defer c.cameraManager.Stop()
	defer c.timerProcess.Stop()

	controls, err := c.engine.LoadFile("main.qml")
	if err != nil {
		return err
	}

	c.window = controls.CreateWindow(nil)

	c.window.Show()
	c.window.Wait()
	log.Println("Exiting")
	return nil
}

func (c *CanonIntervalometer) StartClicked() {
	interval := c.window.ObjectByName("interval").Int("value")
	log.Printf("Start clicked, interval=%d\n", interval)

	c.window.ObjectByName("startButton").Set("enabled", false)
	c.window.ObjectByName("stopButton").Set("enabled", true)
	c.window.ObjectByName("interval").Set("enabled", false)

	go c.timerProcess.Run(interval)
}

func (c *CanonIntervalometer) StopClicked() {
	log.Println("Stop clicked")
	c.window.ObjectByName("startButton").Set("enabled", true)
	c.window.ObjectByName("stopButton").Set("enabled", false)
	c.window.ObjectByName("interval").Set("enabled", true)

	c.timerProcess.Stop()
}

func (c *CanonIntervalometer) LiveViewToggled() {
	if c.window.ObjectByName("liveView").Bool("checked") {
		log.Println("Start Live View")
		c.cameraManager.GetCommandChannel() <- LIVE_VIEW_ENABLED
	} else {
		log.Println("Stop Live View")
		c.cameraManager.GetCommandChannel() <- LIVE_VIEW_DISABLED
	}
}
