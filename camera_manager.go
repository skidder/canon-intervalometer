package main

import (
	"errors"
	"log"

	"github.com/urlgrey/canon-eos-go/eos"
)

type CameraCommand int

const (
	TAKE_PICTURE CameraCommand = iota
	LIVE_VIEW_ENABLED
	LIVE_VIEW_DISABLED
)

type CameraManager struct {
	commandChannel chan CameraCommand
	stopChannel    chan bool
	stopped        bool
	eosClient      *eos.EOSClient
	camera         *eos.CameraModel
}

func NewCameraManager() *CameraManager {
	return &CameraManager{commandChannel: make(chan CameraCommand), stopChannel: make(chan bool)}
}

func (c *CameraManager) Initialize() error {
	c.eosClient = eos.NewEOSClient()
	if err := c.eosClient.Initialize(); err != nil {
		return err
	}
	models, _ := c.eosClient.GetCameraModels()
	if len(models) == 0 {
		return errors.New("No cameras connected")
	}
	c.camera = &models[0]
	c.camera.OpenSession()
	return nil
}

func (c *CameraManager) GetCommandChannel() chan CameraCommand {
	return c.commandChannel
}

func (c *CameraManager) Run() {
	defer c.eosClient.Release()
	defer c.camera.Release()
	defer c.camera.CloseSession()
	c.camera.SetLiveViewOutputDevice(eos.TFT)

	var cmd CameraCommand
	for {
		select {
		case c.stopped = <-c.stopChannel:
			log.Println("Camera Manager returning")
			return
		case cmd = <-c.commandChannel:
			switch cmd {
			case TAKE_PICTURE:
				c.camera.TakePicture()
			case LIVE_VIEW_ENABLED:
				c.camera.SetLiveViewOutputDevice(eos.TFT)
				c.camera.StartLiveView()
			case LIVE_VIEW_DISABLED:
				c.camera.StopLiveView()
			default:
				log.Printf("Unrecognized command: %d\n", cmd)
			}
		}
	}
}

func (c *CameraManager) IsStopped() bool {
	return c.stopped
}

func (c *CameraManager) Stop() {
	c.stopChannel <- true
}
