package code

import "testing"

func TestDestToBinary(t *testing.T) {
	tests := map[string]string{
		"":    "000",
		"M":   "001",
		"D":   "010",
		"MD":  "011",
		"A":   "100",
		"AM":  "101",
		"AD":  "110",
		"AMD": "111",
	}

	for input, expected := range tests {
		actual := DestToBinary(input)
		if actual != expected {
			t.Errorf("Input: %s\nEpected: %s\nActual: %s\n", input, expected, actual)
		}
	}
}

func TestJumpToBinary(t *testing.T) {
	tests := map[string]string{
		"":    "000",
		"JGT": "001",
		"JEQ": "010",
		"JGE": "011",
		"JLT": "100",
		"JNE": "101",
		"JLE": "110",
		"JMP": "111",
	}
	for input, expected := range tests {
		actual := JumpToBinary(input)
		if actual != expected {
			t.Errorf("Input: %s\nEpected: %s\nActual: %s\n", input, expected, actual)
		}
	}
}

func TestCompToBinary(t *testing.T) {
	tests := map[string]string{
		"0":   "0101010",
		"D-1": "0001110",
		"D|A": "0010101",
		"M":   "1110000",
		"M-D": "1000111",
		"-1":  "0111010",
	}
	for input, expected := range tests {
		actual := CompToBinary(input)
		if actual != expected {
			t.Errorf("Input: %s\nEpected: %s\nActual: %s\n", input, expected, actual)
		}
	}
}
