package main

import (
	"errors"
	"log"

	//"log"
	//"os"

	tm "github.com/buger/goterm"
)

// 1   1    1    1       1   1    1        1
// rAi rAo  rBi  rBo   rMi  rMo   commI    commO
var cpuCommands [5]byte = [5]byte{
	0x00,
	0x80, //LDA flags  1000 0000
	0x20, //LDB flags  0010 0000
	0x60, //CPAB flags 0110 0000
	0x02, // HLT flags  0000 0010
	//0x??, // JMP
}

var main_count int = 0
var step int = 0
var halt int = 0

func getBit(a byte, bitNumber int) (int, error) {
	if bitNumber > 8 {
		return 0, errors.New("wrong bit number")
	}

	if bitNumber < 1 {
		return 0, errors.New("wrong bit number")
	}

	bitNumber--
	res := a & (1 << bitNumber) >> bitNumber

	return int(res), nil
}

func execCommand() error {
	var err error
	x := cpuCommands[regCom.value]
	regA.inputEnabled, err = getBit(x, 8)
	if err != nil {
		return err
	}

	regA.outputEnabled, err = getBit(x, 7)
	if err != nil {
		return err
	}

	regB.inputEnabled, err = getBit(x, 6)
	if err != nil {
		return err
	}

	regB.outputEnabled, err = getBit(x, 5)
	if err != nil {
		return err
	}

	regM.inputEnabled, err = getBit(x, 4)
	if err != nil {
		return err
	}

	regM.outputEnabled, err = getBit(x, 3)
	if err != nil {
		return err
	}
	halt, err = getBit(x, 2)
	if err != nil {
		return err
	}

	regM.valueEnabled = 0
	regCom.inputEnabled = 0
	regCom.outputEnabled = 0
	pmCounter.inputEnabled = 0
	pmCounter.outputEnabled = 0
	pmCounter.countEnabled = 0
	regM.value16Enabled = 0
	return nil
}

func setFlags(flags [13]int) {
	regA.inputEnabled = flags[0]        // 1
	regA.outputEnabled = flags[1]       // 2
	regB.inputEnabled = flags[2]        // 3
	regB.outputEnabled = flags[3]       // 4
	regM.inputEnabled = flags[4]        // 5
	regM.outputEnabled = flags[5]       // 6
	regM.valueEnabled = flags[6]        // 7
	regCom.inputEnabled = flags[7]      //8
	regCom.outputEnabled = flags[8]     //9
	pmCounter.inputEnabled = flags[9]   //10
	pmCounter.outputEnabled = flags[10] // 11
	pmCounter.countEnabled = flags[11]  // 12
	regM.value16Enabled = flags[12]     //13
}

func steps(stepCounter byte, programCounter byte) {
	switch stepCounter {
	case 0x00: // Set address from program counter to bus
		//               ai ao bi bo mi mo mv ci co pi po pc m16
		setFlags([13]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0})
	case 0x01: // load address from bus to memory registry ( based on program counter)
		//               ai ao bi bo mi mo mv ci co pi po pc m16
		setFlags([13]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0})
	case 0x02: // load commad from ROM to bus based on value in memory registry
		//               ai ao bi bo mi mo mv ci co pi po pc
		setFlags([13]int{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0})
	case 0x03: // load command from bus (last 4 bits) to command registry
		//               ai ao bi bo mi mo mv ci co pi po pc
		setFlags([13]int{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0})
	case 0x04: // load address from bus ( first 4 bits) to address registry ( for value with offset)
		//               ai ao bi bo mi mo mv ci co pi po pc m16
		setFlags([13]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0})
	case 0x05: // load value from ROM based on value in memory registry
		//              ai ao bi bo mi mo mv ci co pi po pc
		setFlags([13]int{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1})
	case 0x06:
		execCommand()
	default:
		setFlags([13]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	}
}

func setBusFromRom(address byte) error {
	bus = ROM[address]
	return nil
}

func setBustFromProgrammCounter(counter byte) error {
	bus = counter
	return nil
}

func cpuLogging(step int, count byte) {
	if !log_enabled {
		tm.MoveCursor(1, 1)
		tm.Println(tm.Background(tm.Color(tm.Bold("ROM [  0x00  |  0x01  |  0x02  |  0x03  |  0x04  |  0x05  |  0x06  |  0x07  |  0x08  |  0x09  |  0x0A  |  0x0B  |  0x0C  |  0x0D  |  0x0E  |  0x0F  |  0x10  |  0x11  |  0x12  |  0x013  ]              "), tm.BLACK), tm.YELLOW))
		tm.MoveCursor(1, 2)
		tm.Printf("    [%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b]", ROM[0], ROM[1], ROM[2], ROM[3], ROM[4], ROM[5], ROM[6], ROM[8], ROM[9], ROM[10], ROM[11], ROM[12], ROM[13], ROM[14], ROM[15], ROM[16], ROM[17], ROM[18], ROM[19], ROM[20])
		tm.MoveCursor(1, 3)
		tm.Println(tm.Background(tm.Color(tm.Bold(" NPC   PC    SC   REGA  REGB  CMDREG     ADDRREG   ADDR16REG   BUS"), tm.BLACK), tm.YELLOW))
		tm.MoveCursor(1, 4)
		tm.Printf("[ %x ] [ %x ] [ %v ] [ %x ] [ %x ] [ %04b ]  [ %04b ] [ %08b ] [ %08b ]", pmCounter.value[0], count, step, regA.value, regB.value, regCom.value, regM.value, regM.value16, bus)
		tm.MoveCursor(1, 5)
		tm.Println(tm.Background(tm.Color(tm.Bold(" AIE AOE BIE BOE ADIE ADOE ADVE AD16 CIE COE PMI PMO PMC HLT"), tm.BLACK), tm.BLUE))
		tm.MoveCursor(1, 6)
		tm.Printf(" [%v] [%v] [%v] [%v] [%v]  [%v]  [%v]  [%v]  [%v] [%v] [%v] [%v] [%v] [%v]", regA.inputEnabled, regA.outputEnabled, regB.inputEnabled, regB.outputEnabled, regM.inputEnabled, regM.outputEnabled, regM.valueEnabled, regM.value16Enabled, regCom.inputEnabled, regCom.outputEnabled, pmCounter.inputEnabled, pmCounter.outputEnabled, pmCounter.countEnabled, halt)
		tm.Clear()
		tm.Flush()
	} else {
		log.Printf("[%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b]", ROM[0], ROM[1], ROM[2], ROM[3], ROM[4], ROM[5], ROM[6], ROM[8], ROM[9], ROM[10], ROM[11], ROM[12], ROM[13], ROM[14], ROM[15], ROM[16], ROM[17], ROM[18], ROM[19], ROM[20])
		log.Println(" NPC   PC    SC   REGA   REGB    CMDREG      ADDRREG   ADDR16REG       BUS")
		log.Printf("[ %x ] [ %x ] [ %v ] [ %x ] [ %x ] [ %08b ]  [ %08b ]  [ %08b ] [ %08b ]", pmCounter.value[0], count, step, regA.value, regB.value, regCom.value, regM.value, regM.value16, bus)
		log.Println(" AIE AOE BIE BOE ADIE ADOE ADVE AD16 CIE COE PMI PMO PMC HLT")
		log.Printf(" [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v]", regA.inputEnabled, regA.outputEnabled, regB.inputEnabled, regB.outputEnabled, regM.inputEnabled, regM.outputEnabled, regM.valueEnabled, regM.value16Enabled, regCom.inputEnabled, regCom.outputEnabled, pmCounter.inputEnabled, pmCounter.outputEnabled, pmCounter.countEnabled, halt)
		log.Println("---------------------------------------------------------------------------")
	}
}

func cpu() error {
	if halt == 1 {
		tm.Println("HALT !!!")
		cpuLogging(step, pmCounter.value[0])
		return nil
	}
	steps(byte(step), pmCounter.value[0])
	cpuLogging(step, pmCounter.value[0])
	step, _ = stepCounter(step)
	if step == 0 {
		//main_count++
		pmCounter.countEnabled = 1
		setFlags([13]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0})
	}
	return nil
}

func stepCounter(c int) (int, error) {
	if c == 7 {
		return 0, nil
	}
	c++
	return c, nil
}

// add regA to regB and store result to regA
func add() {
	regA.value = regA.value + regB.value
}

func sub() {
	regA.value = regA.value - regB.value
}
