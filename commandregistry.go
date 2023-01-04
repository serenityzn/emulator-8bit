package main

import (
	"time"

	tm "github.com/buger/goterm"
)

type CommandRegistry struct {
	value         byte
	outputEnabled int
	inputEnabled  int
}

func (regcomm *CommandRegistry) readFromBus() {
	regcomm.value = bus >> 4
	//log.Printf("read to command registry %b", regcomm.value)
}

func (regcomm *CommandRegistry) writeToBus() {
	//log.Printf("READING ROM on addrees %v, Value is %b", regm.value, ROM[regm.value])
	bus = regcomm.value
	//log.Printf("BUS is %b", bus)
}

var regCom CommandRegistry

var regComTicker *time.Ticker

func regcom_init() {
	tm.Println("Command Registry Initializing...")
	regCom.value = 0x00
	regCom.inputEnabled = 0
	regCom.outputEnabled = 0
}

func regComlogic() {
	if regCom.inputEnabled == 1 {
		regCom.readFromBus()
	}
	if regCom.outputEnabled == 1 {
		regCom.writeToBus()
	}
}

func registerComRoutine() {
	regComlogic()
}
