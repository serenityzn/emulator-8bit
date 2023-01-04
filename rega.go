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

func rega_init(speed time.Duration, duration string) {
	//tm.MoveCursor(1,1)
	tm.Println("RegA Initializing...")
	//tm.Flush() 
	regA.value = 0x00
	regA.inputEnabled = 0
	regA.outputEnabled = 0
	switch duration {
	case "mili":
		regATicker = time.NewTicker(speed * time.Millisecond)
	case "micro":
		regATicker = time.NewTicker(speed * time.Microsecond)
	case "nano":
		regATicker = time.NewTicker(speed * time.Nanosecond)
	default:
		regATicker = time.NewTicker(speed * time.Millisecond)
	}
}

func regAlogic() {
	if regA.inputEnabled == 1 {
		regA.readFromBus()
	}
	if regA.outputEnabled == 1 {
		regA.writeToBus()
	}
}

func registerARoutine(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			//log.Printf("From reg routing %v", mainClock)
			//log.Println("")
			if mainClock == 1 {
				//	log.Printf("Main clock pules is %v. BUS is %v regA is %v", mainClock, bus, regA)
				regAlogic()
			}
		}
	}
}
