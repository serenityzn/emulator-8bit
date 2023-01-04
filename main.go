package main

import (
	"log"
	"os"
	"time"

	tm "github.com/buger/goterm"
)

var mainClock int = 0
var bus byte
var clockTicker *time.Ticker
var log_enabled bool

func clock(ticker *time.Ticker, socket chan int) {
	for {
		select {
		case <-ticker.C:
			mainClock = invert(mainClock)
			socket <- mainClock
		}
	}
}

func invert(i int) int {
	if i == 0 {
		return 1
	}
	return 0
}

func main() {
	bus = 0xff
	log_enabled = false
	if log_enabled {
		f, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		log.SetOutput(f)
		log.Println("This is a test log entry")
	}

	//printRom()
	tm.Clear()
	tm.MoveCursor(1, 5)
	rega_init()
	regb_init()
	regm_init()
	regcom_init()
	pmc_init()

	clockTicker = time.NewTicker(50 * time.Millisecond)
	socket := make(chan int)
	go clock(clockTicker, socket)

	tm.Clear()
	for {
		switch <-socket {
		case 1:
			registerARoutine()
			registerBRoutine()
			registerMRoutine()
			registerComRoutine()
			pmCounterRoutine()
		case 0:
			cpu()
		case 22:
			clockTicker.Stop()
			os.Exit(0)
		}
	}

}
