package parser

import "testing"

func TestNoFilePath(t *testing.T) {
	parser := Parser{}
	if err := parser.OpenFile(); err.Error() != "Filepath not set" {
		t.Errorf("No error was thrown when no filepath was given")
	}
}

func TestFilePathDoesntExist(t *testing.T) {
	fakeFile := "fakeFilePath.asm"
	parser := Parser{filepath: fakeFile}
	if err := parser.OpenFile(); err.Error() != "open "+fakeFile+": no such file or directory" {
		t.Errorf("Incorrect error was thrown when filepath didn't exist\nGiven error: %s", err.Error())
	}
}

func TestFileHasASMEnding(t *testing.T) {
	file := "badEnding.go"
	parser := Parser{filepath: file}
	if err := parser.OpenFile(); err.Error() != "File must have a .asm ending" {
		t.Errorf("Incorrect error was thrown when file ending was not .asm\nGiven error: %s", err.Error())
	}
}

func TestIsLineCommand(t *testing.T) {
	tests := map[string]bool{
		"     ": false,
		" \t ":  false,
		"   // This is a comment": false,
		"    @22":                 true,
		"   (YO)  ":               true,
		"   M=D":                  true}
	for input, expected := range tests {
		actual := isLineCommand(input)
		if actual != expected {
			t.Errorf("Input: %s\nExpected: %t\nActual: %t", input, expected, actual)
		}
	}
}

func TestCommandType(t *testing.T) {
	tests := map[string]int{
		"  @22":                    ACommand,
		" @R0":                     ACommand,
		"(LOOP)":                   LCommand,
		"   (TEST)   ":             LCommand,
		" @R3 // test":             ACommand,
		" (TEST) // test comment":  LCommand,
		" D=M ":                    CCommand,
		" D=M // test comment":     CCommand,
		" M=D;JGE ":                CCommand,
		" M=D;JMP // test comment": CCommand,
		" M=!M":                    CCommand,
		"D;JGT   // if D>0 (first is greater) goto output_first": CCommand,
		" (ball.new)":          LCommand,
		" (sys.wait$if_true0)": LCommand,
	}

	for input, expected := range tests {
		parser := Parser{currentCommand: input}
		actual := parser.CommandType()
		if actual != expected || parser.currentCommandType != expected {
			t.Errorf("Input: %s\nExpected: %d\nActual: %d", input, expected, actual)
		}
	}
}

func TestCommandTypeSet(t *testing.T) {
	tests := map[string]string{
		"  @22":                                                  "@22",
		" @R0":                                                   "@R0",
		"(LOOP)":                                                 "(LOOP)",
		"   (TEST)   ":                                           "(TEST)",
		" @R3 // test":                                           "@R3",
		" (TEST) // test comment":                                "(TEST)",
		" D=M ":                                                  "D=M",
		" D=M // test comment":                                   "D=M",
		" M=D;JGE ":                                              "M=D;JGE",
		" M=D;JMP // test comment":                               "M=D;JMP",
		" M=D+1;JMP":                                             "M=D+1;JMP",
		" D+1;JMP // comment":                                    "D+1;JMP",
		" //comment":                                             "",
		" M=D-1;JMP":                                             "M=D-1;JMP",
		" D-1":                                                   "D-1",
		" M=-1":                                                  "M=-1",
		" M=!M":                                                  "M=!M",
		" M=D&M":                                                 "M=D&M",
		" M=D|M":                                                 "M=D|M",
		" (OUTPUT_TEST) // testcomment":                          "(OUTPUT_TEST)",
		"D;JGT   // if D>0 (first is greater) goto output_first": "D;JGT",
		" (sys.wait$if_true0)":                                   "(sys.wait$if_true0)",
		" (ball.new)":                                            "(ball.new)",
	}

	for input, expected := range tests {
		parser := Parser{currentCommand: input}
		parser.CommandType()
		if parser.currentCommand != expected {
			t.Errorf("Input: %s\nExpected: %s\nActual: %s\n", input, expected, parser.currentCommand)
		}
	}
}

func TestAGetSymbol(t *testing.T) {
	aTests := map[string]string{
		"@22":                "22",
		"@R0":                "R0",
		"@test":              "test",
		"@OUTPUT_TEST":       "OUTPUT_TEST",
		"@sys.wait$if_true0": "sys.wait$if_true0",
		"@ball.new":          "ball.new",
	}

	for input, expected := range aTests {
		parser := Parser{currentCommand: input, currentCommandType: ACommand}
		actual := parser.GetSymbol()
		if actual != expected {
			t.Errorf("Input: %s\nExpected: %s\nActual: %s\n", input, expected, actual)
		}
	}
}

func TestLGetSymbol(t *testing.T) {
	lTests := map[string]string{
		"(LOOP)":              "LOOP",
		"(TEST)":              "TEST",
		"(OUTPUT_TEST)":       "OUTPUT_TEST",
		"(sys.wait$if_true0)": "sys.wait$if_true0",
		"(ball.new)":          "ball.new",
	}

	for input, expected := range lTests {
		parser := Parser{currentCommand: input, currentCommandType: LCommand}
		actual := parser.GetSymbol()
		if actual != expected {
			t.Errorf("Input: %s\nExpected: %s\nActual: %s\n", input, expected, actual)
		}
	}
}

func TestCGetSymbol(t *testing.T) {
	cTests := []string{
		"M=D",
		"D=M;JMP",
	}

	for index := range cTests {
		input := cTests[index]
		expected := ""
		parser := Parser{currentCommand: input, currentCommandType: CCommand}
		actual := parser.GetSymbol()
		if actual != expected {
			t.Errorf("Input: %s\nExpected: %s\nActual: %s\n", input, expected, actual)
		}
	}
}

func TestGetDestination(t *testing.T) {
	tests := map[string]string{
		"D=M":       "D",
		"D=M;JMP":   "D",
		"M=A;JGE":   "M",
		";JMP":      "",
		"D=M+1;JMP": "D",
		"D+1;JMP":   "",
		"M+1":       "",
		"M=-1":      "M",
		"M=!M":      "M",
		"M=D&M":     "M",
		"M=D|M":     "M",
	}

	for input, expected := range tests {
		parser := Parser{currentCommand: input, currentCommandType: CCommand}
		actual := parser.GetDestination()
		if actual != expected {
			t.Errorf("Input: %s\nExpected: %s\nActual: %s\n", input, expected, actual)
		}
	}
}

func TestGetDestinationWrongCommand(t *testing.T) {
	aTests := []string{
		"@R2",
		"@test",
		"@88",
	}

	lTests := []string{
		"(LOOP)",
		"(TEST)",
	}

	for index := range aTests {
		input := aTests[index]
		parser := Parser{currentCommand: input, currentCommandType: ACommand}
		actual := parser.GetDestination()
		if actual != "" {
			t.Errorf("Input: %s\nExpected a blank string\nActual: %s\n", input, actual)
		}
	}

	for index := range lTests {
		input := lTests[index]
		parser := Parser{currentCommand: input, currentCommandType: LCommand}
		actual := parser.GetDestination()
		if actual != "" {
			t.Errorf("Input: %s\nExpected a blank string\nActual: %s\n", input, actual)
		}
	}
}

func TestGetComp(t *testing.T) {
	tests := map[string]string{
		"D=M+1;JMP": "M+1",
		"D=M+1":     "M+1",
		"M+1;JMP":   "M+1",
		"M+1":       "M+1",
		"D=M":       "M",
		"M=-1":      "-1",
		"M=!M":      "!M",
		"M=!D":      "!D",
		"M=D&M":     "D&M",
		"M=D|M":     "D|M",
	}

	for input, expected := range tests {
		parser := Parser{currentCommand: input, currentCommandType: CCommand}
		actual := parser.GetComp()
		if actual != expected {
			t.Errorf("Input: %s\nExpected: %s\nActual: %s\n", input, expected, actual)
		}
	}
}

func TestGetJump(t *testing.T) {
	tests := map[string]string{
		"D=M+1;JMP": "JMP",
		"D=M+1":     "",
		"M+1;JGE":   "JGE",
		"M+1":       "",
		"D=M":       "",
		"M=-1":      "",
	}

	for input, expected := range tests {
		parser := Parser{currentCommand: input, currentCommandType: CCommand}
		actual := parser.GetJump()
		if actual != expected {
			t.Errorf("Input: %s\nExpected: %s\nActual: %s\n", input, expected, actual)
		}
	}
}

func TestVerifyPredifinedSymbos(t *testing.T) {
	tests := map[string]int{
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15,
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KBD":    24576,
	}

	for input, expected := range tests {
		parser := New("test")
		actual := parser.GetAddress(input)
		if !parser.ContainsSymbol(input) || actual != expected {
			t.Errorf("Input: %s\nExpected: %d\nActual: %d\n", input, expected, actual)
		}
	}
}
