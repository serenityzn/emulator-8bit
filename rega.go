package main

import (
	"time"

	tm "github.com/buger/goterm"
)

type RegisterA struct {
	value         byte
	outputEnabled int
	inputEnabled  int
}

func (rega *RegisterA) readFromBus() {
	rega.value = bus
}

func (rega *RegisterA) writeToBus() {
	bus = rega.value
}

var regA RegisterA

var regATicker *time.Ticker

func rega_init() {
	//tm.MoveCursor(1,1)
	tm.Println("RegA Initializing...")
	//tm.Flush()
	regA.value = 0x00
	regA.inputEnabled = 0
	regA.outputEnabled = 0
}

func regAlogic() {
	if regA.inputEnabled == 1 {
		regA.readFromBus()
	}
	if regA.outputEnabled == 1 {
		regA.writeToBus()
	}
}

func registerARoutine() {
	regAlogic()
}
