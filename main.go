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

func clock(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			mainClock = invert(mainClock)
		}
	}
}

func invert(i int) int {
	if i == 0 {
		return 1
	}
	return 0
}

func setSpeed(value time.Duration, duration string) {
	rega_init(value, duration)
	regb_init(value, duration)
	regm_init(value, duration)
	regcom_init(value, duration)
	pmc_init(value, duration)
	//clockTicker = time.NewTicker(value * 20 * time.Millisecond)
}

func main() {
	log_enabled = true

	if log_enabled {
		f, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		log.SetOutput(f)
		log.Println("This is a test log entry")
	}
	/* err := readRomFromFile("rom.bin")
	if err != nil {
		log.Panicln(err)
	} */
	printRom()
	tm.Clear()
	tm.MoveCursor(1, 5)

	//setSpeed(10, "mili")
	rega_init(100, "micro")
	regb_init(100, "micro")
	regm_init(100, "micro")
	regcom_init(100, "micro")
	//pmc_init(250, "mili")
	clockTicker = time.NewTicker(1000 * time.Microsecond)

	go clock(clockTicker)
	go registerARoutine(regATicker)
	go RegisterBRoutine(regbTicker)
	go registerMRoutine(regMTicker)
	go regComRoutine(regComTicker)
	//go pmCounterRoutine(pmCounterTicker)
	go cpu(clockTicker)

	for {
		time.Sleep(100 * time.Microsecond)
	}
	clockTicker.Stop()

}
