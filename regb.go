package main

import (
	"time"

	tm "github.com/buger/goterm"
)

type RegisterB struct {
	value         byte
	outputEnabled int
	inputEnabled  int
}

func (regb *RegisterB) readFromBus() {
	regb.value = bus
}

func (regb *RegisterB) writeToBus() {
	bus = regb.value
}

var regB RegisterB
var regbTicker *time.Ticker

func regb_init() {
	//tm.MoveCursor(1,1)
	tm.Println("RegB Initializing...")
	//tm.Flush()

	regB.value = 0x00
	regB.inputEnabled = 0
	regB.outputEnabled = 0

}

func regblogic() {
	if regB.inputEnabled == 1 {
		regB.readFromBus()
	}
	if regB.outputEnabled == 1 {
		regB.writeToBus()
	}
}

func registerBRoutine() {
	regblogic()
}
