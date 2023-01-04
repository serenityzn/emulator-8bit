package main

import (
	"time"

	tm "github.com/buger/goterm"
)

type MemoryAddressRegistry struct {
	value         byte
	outputEnabled int
	inputEnabled  int
	valueEnabled int
}

func (regm *MemoryAddressRegistry) readFromBus() {
	regm.value = bus & 0xf
}

func (regm *MemoryAddressRegistry) writeToBus()  error{
	if regm.valueEnabled == 1 {
		//log.Printf("REading value from ROM reg value is [%v] value with offset is [%v] Value is [%x]",regm.value,regm.value+0xf,ROM[regm.value+0xf])
		bus = ROM[regm.value+0xf]
		return nil
	}
	bus = ROM[regm.value]
	return nil
}

var regM MemoryAddressRegistry

var regMTicker *time.Ticker

func regm_init(speed time.Duration, duration string) {
	tm.Println("Memory Registry Initializing...")
	regM.value = 0x00
	regM.inputEnabled = 0
	regM.outputEnabled = 0
	regM.valueEnabled = 0
	regMTicker = time.NewTicker(speed * time.Millisecond)
	switch duration {
	case "mili":
		regMTicker = time.NewTicker(speed * time.Millisecond)
	case "micro":
		regMTicker = time.NewTicker(speed * time.Microsecond)
	case "nano":
		regMTicker = time.NewTicker(speed * time.Nanosecond)
	default:
		regMTicker = time.NewTicker(speed * time.Millisecond)
	}
}

func regMlogic() {
	if regM.inputEnabled == 1 {
		regM.readFromBus()
	}
	if regM.outputEnabled == 1 {
		regM.writeToBus()
	}
}

func registerMRoutine(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			if mainClock == 1 {
				regMlogic()
			}
		}
	}
}
