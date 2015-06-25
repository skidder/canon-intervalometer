package main

import (
	"fmt"
	"gopkg.in/qml.v1"
	"os"
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
	engine *qml.Engine
	window *qml.Window
}

func NewCanonIntervalometer(e *qml.Engine) *CanonIntervalometer {
	return &CanonIntervalometer{engine: e}
}

func (c *CanonIntervalometer) Initialize() {
	context := c.engine.Context()
	context.SetVar("intervalometer", c)
}

func (c *CanonIntervalometer) Run() error {
	controls, err := c.engine.LoadFile("main.qml")
	if err != nil {
		return err
	}

	c.window = controls.CreateWindow(nil)

	c.window.Show()
	c.window.Wait()
	return nil
}

func (c *CanonIntervalometer) StartClicked() {
	fmt.Println("Start clicked")
	startButton := c.window.ObjectByName("startButton")
	startButton.Set("enabled", false)
	stopButton := c.window.ObjectByName("stopButton")
	stopButton.Set("enabled", true)

	fmt.Printf("Interval: %d\n", c.window.ObjectByName("interval").Int("value"))
}

func (c *CanonIntervalometer) StopClicked() {
	fmt.Println("Stop clicked")
	startButton := c.window.ObjectByName("startButton")
	startButton.Set("enabled", true)
	stopButton := c.window.ObjectByName("stopButton")
	stopButton.Set("enabled", false)
}
