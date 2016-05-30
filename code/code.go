package code

import "strings"

var destinationMap = map[string]string{
	"":    "000",
	"M":   "001",
	"D":   "010",
	"MD":  "011",
	"A":   "100",
	"AM":  "101",
	"AD":  "110",
	"AMD": "111",
}

var jumpMap = map[string]string{
	"":    "000",
	"JGT": "001",
	"JEQ": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}

var mMap = map[string]string{
	"M":   "110000",
	"!M":  "110001",
	"-M":  "110011",
	"M+1": "110111",
	"M-1": "110010",
	"D+M": "000010",
	"D-M": "010011",
	"M-D": "000111",
	"D&M": "000000",
	"D|M": "010101",
}

var compMap = map[string]string{
	"0":   "101010",
	"1":   "111111",
	"-1":  "111010",
	"D":   "001100",
	"A":   "110000",
	"!D":  "001101",
	"!A":  "110001",
	"-D":  "001111",
	"-A":  "110011",
	"D+1": "011111",
	"A+1": "110111",
	"D-1": "001110",
	"A-1": "110010",
	"D+A": "000010",
	"D-A": "010011",
	"A-D": "000111",
	"D&A": "000000",
	"D|A": "010101",
}

// DestToBinary converts the input destination command to the binary representation
func DestToBinary(input string) string {
	return destinationMap[input]
}

// JumpToBinary converts the input jump command to the binary representation
func JumpToBinary(input string) string {
	return jumpMap[input]
}

// CompToBinary converts the input comp command to the binary representation
func CompToBinary(input string) string {
	finalString := ""
	if strings.Contains(input, "M") {
		finalString += "1"
		finalString += mMap[input]
	} else {
		finalString += "0"
		finalString += compMap[input]
	}
	return finalString
}
