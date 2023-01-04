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


func regb_init(speed time.Duration, duration string) {
	//tm.MoveCursor(1,1)
	tm.Println("RegB Initializing...")
	//tm.Flush()

	regB.value = 0x00
	regB.inputEnabled = 0
	regB.outputEnabled = 0
	switch duration {
	case "mili":
		regbTicker = time.NewTicker(speed * time.Millisecond)
	case "micro":
		regbTicker = time.NewTicker(speed * time.Microsecond)
	case "nano":
		regbTicker = time.NewTicker(speed * time.Nanosecond)
	default:
		regbTicker = time.NewTicker(speed * time.Millisecond)
	}
}

func regblogic() {
	if regB.inputEnabled == 1 {
		regB.readFromBus()
	}
	if regB.outputEnabled == 1 {
		regB.writeToBus()
	}
}

func RegisterBRoutine(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			//log.Printf("From reg routing %v", mainClock)
			//log.Println("")
			if mainClock == 1 {
				//log.Printf("Main clock pules is %v. BUS is %v regb is %v", mainClock, bus, regB)
				regblogic()
			}
		}
	}
}