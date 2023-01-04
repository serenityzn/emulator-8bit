package main

import (
	"log"
	"time"

	tm "github.com/buger/goterm"
)

type ProgramCounter struct {
	value         [2]byte // [1] 1111 1111    [0] 1111 1111
	outputEnabled int
	inputEnabled  int
	countEnabled  int
}

func (pc *ProgramCounter) readFromBus() {
	pc.value[0] = bus
	pc.value[1] = 0x00
}

func (pc *ProgramCounter) writeToBus() {
	bus = pc.value[0]
}

var pmCounter ProgramCounter

var pmCounterTicker *time.Ticker

func pmc_init() {
	tm.Println("Program Counter Initializing...")
	pmCounter.value[0] = 0x00
	pmCounter.value[1] = 0x00
	pmCounter.inputEnabled = 0
	pmCounter.outputEnabled = 0
	pmCounter.countEnabled = 0
}

func pmCounterlogic() {
	if pmCounter.countEnabled == 1 {
		count()
		pmCounter.countEnabled = 0
	}
	if pmCounter.inputEnabled == 1 {
		pmCounter.readFromBus()
		log.Printf("reading from BUS...............")
	}
	if pmCounter.outputEnabled == 1 {
		pmCounter.writeToBus()
	}
}

func count() {
	if pmCounter.value[0] == 0xff {
		pmCounter.value[1]++
	}
	if pmCounter.value[0] == 0xff && pmCounter.value[1] == 0xff {
		pmCounter.value[0] = 0x00
		pmCounter.value[1] = 0x00
	}
	pmCounter.value[0]++
}

func pmCounterRoutine() {
	pmCounterlogic()
}
