package main

import (
	//"errors"
	"fmt"
	"os"
)

var RAM [256]byte
//var ROM []byte
/* 1. -> count = 0x00 -> bus
2.bus(0x00) -> memreg
3. via mem reg func read from ROM[0x00] -> bus
4 bus(100101) -(4 bits) -> com reg [10]
5. bus (100101) -(4bits) -> mem reg[0101] */

 var ROM []byte = []byte{
	0x25, //  LDB 0xfa
	0x11,// LDA 0x0e
	0x30, // CPAB
	0x10, // 0001 0000 LDA 0x0e  LDA 0x00
	0x15, // 0001 0101 LDA 0x4a  LDA 0x03
	0x25, // 0010 0101 LDB 0x4a  LDB 0x03
	0x16, //0001 0110 LDA 0x8f   LDA 0x06
	0x30, //0011 0000 CPAB       CPAB
	0x11, //0001 LDA 0x12
	0x30, //0011 CPAB 0x12 -> B
	0x40, // HLT
	0x00,
	0x00,
	0x00,
	0x00,
	0xaa, // Data address space
	0x0e,
	0x12,
	0xff,
	0x03,
	0xfa,
	0xfb,
	0xfc,
	0xfe,
	0xff,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
} 

func ramInit() {
	for i := 0x00; i < 0xff; i++ {
		RAM[i] = 0x00
	}
}

func printRam() {
	fmt.Println(RAM)
}

func printRom() {
	for i := 0; i < len(ROM); i++ {
		fmt.Printf("%x ", ROM[i])
	}
}

func readRomFromFile(fileName string) error {
	data, err := os.ReadFile(fileName)
	check(err)
	for i := 0; i < len(data); i++ {
		ROM = append(ROM, data[i])
	}
	fmt.Printf("ROM size is %v", len(ROM))

	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
