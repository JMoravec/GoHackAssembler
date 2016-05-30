package main

import (
	"fmt"
	"hack_assembler/code"
	"hack_assembler/parser"
	"os"
	"regexp"
	"strconv"
)

func main() {
	argFilepath := os.Args[1]
	if _, err := os.Stat(argFilepath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	symbolParser := parser.New(argFilepath)
	if err := symbolParser.OpenFile(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err := symbolParser.Advance()
	for err == nil {
		commandType := symbolParser.CommandType()
		if commandType == parser.LCommand {
			symbolParser.AddSymbolLineNumber(symbolParser.GetSymbol())
		}
		err = symbolParser.Advance()
	}

	newFilePath := os.Args[2]
	newFile, err := os.Create(newFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	defer newFile.Close()

	mainParser := parser.NewWithSymbol(argFilepath, symbolParser.GetSymbolTable())
	if err := mainParser.OpenFile(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err = mainParser.Advance()
	for err == nil {
		newLine := ""
		commandType := mainParser.CommandType()

		if commandType == parser.ACommand {
			newLine = "0"

			//assumes symbol is number already
			reg := regexp.MustCompile(`[A-z]`)
			currentCommand := mainParser.GetSymbol()
			regTest := reg.FindStringSubmatch(currentCommand)

			var symbolInt int64

			if regTest != nil {
				if !mainParser.ContainsSymbol(currentCommand) {
					mainParser.AddRamSymbol(currentCommand)
				}
				symbolInt = int64(mainParser.GetAddress(currentCommand))

			} else {
				symbolInt, err = strconv.ParseInt(currentCommand, 10, 64)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			symbol := strconv.FormatInt(symbolInt, 2)
			for len(symbol) < 15 {
				symbol = "0" + symbol
			}

			newLine += symbol
		}

		if commandType == parser.CCommand {
			newLine = "111"

			destination := mainParser.GetDestination()
			comp := mainParser.GetComp()
			jump := mainParser.GetJump()

			newLine += code.CompToBinary(comp)
			newLine += code.DestToBinary(destination)
			newLine += code.JumpToBinary(jump)

		}

		if commandType != parser.LCommand {
			newLine += "\n"
		}

		if _, err := newFile.WriteString(newLine); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		err = mainParser.Advance()
	}
}
