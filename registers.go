package main

import (
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

type AddressRegistry struct {
	value          byte
	value16        byte
	outputEnabled  int
	inputEnabled   int
	valueEnabled   int
	value16Enabled int
}

func (regm *AddressRegistry) readFromAdressROM() {
	address := addressROM.address[pmCounter.value[0]] * addressROM.offset[pmCounter.value[0]]
	regm.value16 = ROM[address]
}

func (regm *AddressRegistry) readFromBus() {
	regm.value = bus & 0xf
}

func (regm *AddressRegistry) writeToBus() error {

	if regm.value16Enabled == 1 {
		bus = regm.value16
		return nil
	}

	if regm.valueEnabled == 1 {
		bus = ROM[regm.value+0xf]
		return nil
	}

	bus = ROM[regm.value]
	return nil
}

type CommandRegistry struct {
	value         byte
	outputEnabled int
	inputEnabled  int
}

func (regcomm *CommandRegistry) readFromBus() {
	regcomm.value = bus >> 4
}

func (regcomm *CommandRegistry) writeToBus() {
	bus = regcomm.value
}

var regA RegisterA
var regB RegisterB
var regM AddressRegistry
var regCom CommandRegistry

func rega_init() {
	tm.Println("RegA Initializing...")
	regA.value = 0x00
	regA.inputEnabled = 0
	regA.outputEnabled = 0
}

func regb_init() {
	tm.Println("RegB Initializing...")
	regB.value = 0x00
	regB.inputEnabled = 0
	regB.outputEnabled = 0

}

func regm_init() {
	tm.Println("Memory Registry Initializing...")
	regM.value = 0x00
	regM.inputEnabled = 0
	regM.outputEnabled = 0
	regM.valueEnabled = 0
}

func regcom_init() {
	tm.Println("Command Registry Initializing...")
	regCom.value = 0x00
	regCom.inputEnabled = 0
	regCom.outputEnabled = 0
}

func regAlogic() {
	if regA.inputEnabled == 1 {
		regA.readFromBus()
	}
	if regA.outputEnabled == 1 {
		regA.writeToBus()
	}
}

func regBlogic() {
	if regB.inputEnabled == 1 {
		regB.readFromBus()
	}
	if regB.outputEnabled == 1 {
		regB.writeToBus()
	}
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

func regComlogic() {
	if regCom.inputEnabled == 1 {
		regCom.readFromBus()
	}
	if regCom.outputEnabled == 1 {
		regCom.writeToBus()
	}
}
