package main

import (
	"fmt"
)

var CHARS_MAPPING_1 = " !\"#$%&'()*+,-./"
var CHARS_MAPPING_2 = "0123456789:;<=>?"
var CHARS_MAPPING_3 = "@ABCDEFGHIJKLMNO"
var CHARS_MAPPING_4 = "PQRSTUVWXYZ[!]^_"
var CHARS_MAPPING_5 = "`abcdefghijklmno"
var CHARS_MAPPING_6 = "pqrstuvwxyz{|}~"

func main() {
	// fmt.Println(getCharByMapping("@"))
	WriteString("i")
}

func CountBits(decimal uint) int {
	binary := []uint{}
	result := 0
	for decimal != 0 {
		binary = append(binary, decimal%2)
		decimal = decimal / 2
	}
	for _, v := range binary {
		if v == 1 {
			result++
		}
	}
	return result
}

func getCharBits(decimal int) [4]bool {
	binary := [4]bool{false, false, false, false}
	counter := 3
	for decimal != 0 {
		delim := decimal % 2
		if delim == 0 {
			binary[counter] = false
		} else {
			binary[counter] = true
		}
		decimal = decimal / 2
		counter--
	}
	return binary
}

func getCharByMapping(charCode int) [2][4]bool {
	if charCode < 32 && charCode > 128 {
		return [2][4]bool{}
	}
	currentGroupVertical := [4]bool{false, false, false, false}
	currentGroupHorizontal := [4]bool{false, false, true, false}
	if charCode >= 32 && charCode <= 47 {
		currentGroupHorizontal = [4]bool{false, false, true, false}
		currentGroupVertical = getCharBits(charCode - 32)
	} else if charCode >= 48 && charCode <= 63 {
		currentGroupHorizontal = [4]bool{false, false, true, true}
		currentGroupVertical = getCharBits(charCode - 48)
	} else if charCode >= 64 && charCode <= 79 {
		currentGroupHorizontal = [4]bool{false, true, false, false}
		currentGroupVertical = getCharBits(charCode - 64)
	} else if charCode >= 80 && charCode <= 95 {
		currentGroupHorizontal = [4]bool{false, true, false, true}
		currentGroupVertical = getCharBits(charCode - 80)
	} else if charCode >= 96 && charCode <= 111 {
		currentGroupHorizontal = [4]bool{false, true, true, false}
		currentGroupVertical = getCharBits(charCode - 96)
	} else if charCode >= 112 && charCode <= 126 {
		currentGroupHorizontal = [4]bool{false, true, true, true}
		currentGroupVertical = getCharBits(charCode - 112)
	}
	return [2][4]bool{currentGroupVertical, currentGroupHorizontal}
}

func WriteString(str string) {
	for _, v := range str {
		responder := getCharByMapping(int(v))
		if len(responder) > 1 {
			fmt.Println(responder)
		}
	}
}
