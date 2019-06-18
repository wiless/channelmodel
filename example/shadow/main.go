package main

import (
	"github.com/wiless/channelmodel"
)

func main() {
	var sg CM.ShadowGrid
	sg.SetEnv(CM.RMA)
	sg.SetGridSize(20, 500)
	sg.Create(CM.LOS)

}
