package main

import (
	tm "github.com/buger/goterm"
)

type MemoryAddressRegistry struct {
	value          byte
	value16        byte
	outputEnabled  int
	inputEnabled   int
	valueEnabled   int
	value16Enabled int
}

func (regm *MemoryAddressRegistry) readFromAdressROM() {
	address := addressROM.address[pmCounter.value[0]] * addressROM.offset[pmCounter.value[0]]
	regm.value16 = ROM[address]
}

func (regm *MemoryAddressRegistry) readFromBus() {
	regm.value = bus & 0xf
}

func (regm *MemoryAddressRegistry) writeToBus() error {

	if regm.value16Enabled == 1 {
		bus = regm.value16
		return nil
	}

	if regm.valueEnabled == 1 {
		//log.Printf("REading value from ROM reg value is [%v] value with offset is [%v] Value is [%x]",regm.value,regm.value+0xf,ROM[regm.value+0xf])
		bus = ROM[regm.value+0xf]
		return nil
	}

	bus = ROM[regm.value]
	return nil
}

var regM MemoryAddressRegistry

func regm_init() {
	tm.Println("Memory Registry Initializing...")
	regM.value = 0x00
	regM.inputEnabled = 0
	regM.outputEnabled = 0
	regM.valueEnabled = 0
}

func regMlogic() {
	if regM.inputEnabled == 1 {
		regM.readFromBus()
		regM.readFromAdressROM()
	}
	if regM.outputEnabled == 1 {
		regM.writeToBus()
	}
}

func registerMRoutine() {
	regMlogic()
}
