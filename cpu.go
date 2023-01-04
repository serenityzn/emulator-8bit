package main

import (
	"errors"
	"log"

	//"log"
	//"os"
	"time"

	tm "github.com/buger/goterm"
)

// 1   1    1    1       1   1    1        1
// rAi rAo  rBi  rBo   rMi  rMo   commI    commO
var cpuCommands [5]byte = [5]byte{
	0x00,
	0x80, //LDA flags  1000 0000
	0x20, //LDB flags  0010 0000
	0x60, //CPAB flags 0110 0000
	0x00, // HLT flags  0000 0000
	//0x??, // JMP
}

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
	regM.valueEnabled = 0
	regCom.inputEnabled = 0
	regCom.outputEnabled = 0
	pmCounter.inputEnabled = 0
	pmCounter.outputEnabled = 0
	pmCounter.countEnabled = 0
	return nil
}

func setFlags(flags [12]int) {
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
}

func steps(stepCounter byte, programCounter byte) {
	switch stepCounter {
	case 0x00: // Set address from program counter to bus
		//               ai ao bi bo mi mo mv ci co pi po pc
		setFlags([12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0})
		setBustFromProgrammCounter(programCounter)
	case 0x01: // load address from bus to memory registry ( based on program counter)
		//               ai ao bi bo mi mo mv ci co pi po pc
		setFlags([12]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0})
	case 0x02: // load commad from ROM to bus based on value in memory registry
		//               ai ao bi bo mi mo mv ci co pi po pc
		setFlags([12]int{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0})
	case 0x03: // load command from bus (last 4 bits) to command registry
		//               ai ao bi bo mi mo mv ci co pi po pc
		setFlags([12]int{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0})
	case 0x04: // load address from bus ( first 4 bits) to address registry ( for value with offset)
		//               ai ao bi bo mi mo mv ci co pi po pc
		setFlags([12]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0})
	case 0x05: // load value from ROM based on value in memory registry
		//              ai ao bi bo mi mo mv ci co pi po pc
		setFlags([12]int{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0})
	case 0x06:
		execCommand()
	default:
		setFlags([12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
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

func cpu(ticker *time.Ticker) {

	tm.Println("Initializing CPU...")

	var count int = 0
	var step int = 0
	for {
		select {
		case <-ticker.C:
			if mainClock == 1 {

				//steps(byte(step), pmCounter.value[0])
				steps(byte(step), byte(count))
				if !log_enabled {
					//log.Printf("CPU tick. PR Count is %v, Step count is %v, bus [ %x ], regA %x regB %x regCom %b", count, step, bus, regA.value, regB.value, regCom.value)
					tm.MoveCursor(1, 10)
					//tm.Printf("CPU tick. [ PR Counter is %v ], [ Step counter  %v ], bus [ %x ], regA %x regB %x regCom %b\n", count, step, bus, regA.value, regB.value, regCom.value)
					tm.Println(tm.Background(tm.Color(tm.Bold(" NPC   PC    SC   REGA  REGB  CMDREG MEMREG BUS"), tm.BLACK), tm.YELLOW))
					tm.MoveCursor(1, 11)
					tm.Printf("[ %x ] [ %x ] [ %v ] [ %x ] [ %x ] [ %b ]  [ %b ] [ %b ]", pmCounter.value[0], count, step, regA.value, regB.value, regCom.value, regM.value, bus)
					tm.MoveCursor(1, 12)
					tm.Println(tm.Background(tm.Color(tm.Bold(" AIE AOE BIE BOE MIE MOE MVE CIE COE PMI PMO PMC"), tm.BLACK), tm.BLUE))
					tm.MoveCursor(1, 13)
					tm.Printf(" [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v]", regA.inputEnabled, regA.outputEnabled, regB.inputEnabled, regB.outputEnabled, regM.inputEnabled, regM.outputEnabled, regM.valueEnabled, regCom.inputEnabled, regCom.outputEnabled, pmCounter.inputEnabled, pmCounter.outputEnabled, pmCounter.countEnabled)
					tm.Clear()
					tm.Flush()
				} else {
					log.Printf("[%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b]", ROM[0], ROM[1], ROM[2], ROM[3], ROM[4], ROM[5], ROM[6], ROM[8], ROM[9], ROM[10], ROM[11], ROM[12], ROM[13], ROM[14], ROM[15], ROM[16], ROM[17], ROM[18], ROM[19], ROM[20])
					log.Println(" NPC   PC    SC   REGA  REGB  CMDREG MEMREG BUS")
					log.Printf("[ %x ] [ %x ] [ %v ] [ %x ] [ %x ] [ %b ]  [ %b ] [ %b ]", pmCounter.value[0], count, step, regA.value, regB.value, regCom.value, regM.value, bus)
					log.Println(" AIE AOE BIE BOE MIE MOE MVE CIE COE PMI PMO PMC")
					log.Printf(" [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v]", regA.inputEnabled, regA.outputEnabled, regB.inputEnabled, regB.outputEnabled, regM.inputEnabled, regM.outputEnabled, regM.valueEnabled, regCom.inputEnabled, regCom.outputEnabled, pmCounter.inputEnabled, pmCounter.outputEnabled, pmCounter.countEnabled)
					log.Println("---------------------------------------------------------------------------")
				}

				step, _ = stepCounter(step)
				if step == 0 {
					count++
					setFlags([12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
					/* 	if pmCounter.value[0] == 0xff {
						os.Exit(0)
					} */
				}
			}
		}
	}
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
