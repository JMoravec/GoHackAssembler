package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	// ACommand enum for an a command
	ACommand = iota
	// CCommand enum for a c command
	CCommand = iota
	// LCommand enum for a l pseudo command
	LCommand = iota
)

// Parser parses an hack asm file into its respective parts
type Parser struct {
	filepath           string
	fileScanner        *bufio.Scanner
	currentCommand     string
	currentCommandType int
	LineNumber         int
	symbols            map[string]int
	nextAddress        int
}

// New creates a new parser with the given filepath
func New(filepath string) *Parser {
	newParse := &Parser{filepath: filepath, LineNumber: 0, nextAddress: 16}
	newParse.symbols = map[string]int{
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KBD":    24576,
	}
	for i := 0; i < 16; i++ {
		newParse.addSymbol("R"+strconv.Itoa(i), i)
	}
	return newParse
}

// New creates a new parser with the given filepath and symbol table
func NewWithSymbol(filepath string, OldSymbols map[string]int) *Parser {
	return &Parser{filepath: filepath, LineNumber: 0, symbols: OldSymbols, nextAddress: 16}
}

// OpenFile opens a file and sets the Parser's scanner
func (p *Parser) OpenFile() error {
	if p.filepath == "" {
		return errors.New("Filepath not set")
	}

	if p.filepath[len(p.filepath)-4:] != ".asm" {
		return errors.New("File must have a .asm ending")
	}

	file, err := os.Open(p.filepath)
	if err != nil {
		return err
	}
	p.fileScanner = bufio.NewScanner(file)
	return nil
}

// Advance moves the parser to the next command
func (p *Parser) Advance() error {
	p.fileScanner.Scan()
	for {
		if err := p.fileScanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading file: ", err)
			return err
		}

		if isLineCommand(p.fileScanner.Text()) {
			p.currentCommand = p.fileScanner.Text()
			return nil
		}

		if !p.fileScanner.Scan() && p.fileScanner.Err() == nil {
			return errors.New("EOF")
		}
	}
}

// CommandType returns the current command type (A, C, or L)
// Has the side affect of setting the current command with no whitespace
// or line comments
func (p *Parser) CommandType() int {
	strippedCommand := strings.Replace(p.currentCommand, " ", "", -1)
	if strippedCommand[0:2] == `//` {
		p.currentCommand = ""
		p.currentCommandType = -1
		return -1
	}

	aRe := regexp.MustCompile(`(@[^/]+)`)
	aTest := aRe.FindStringSubmatch(strippedCommand)
	if aTest != nil {
		p.currentCommand = aTest[0]
		p.currentCommandType = ACommand
		p.LineNumber++
		return ACommand
	}

	if strippedCommand[0:1] == "(" {
		lRe := regexp.MustCompile(`\(.+\)`)
		lTest := lRe.FindStringSubmatch(strippedCommand)
		if lTest != nil {
			p.currentCommand = lTest[0]
			p.currentCommandType = LCommand
			return LCommand
		}
	}

	cRe := regexp.MustCompile(`((\w+=)?([-!])?\w([+&|-]\w)?(;\w+)?)`)
	cTest := cRe.FindStringSubmatch(strippedCommand)
	if cTest != nil {
		p.currentCommand = cTest[0]
		p.currentCommandType = CCommand
		p.LineNumber++
		return CCommand
	}

	p.currentCommand = ""
	p.currentCommandType = -1
	return -1
}

// GetSymbol gets the symbol or number for an A command (R0 if @R0)
// or the name of the goto symbol (TEST if (TEST))
func (p *Parser) GetSymbol() string {
	if p.currentCommandType == ACommand {
		re := regexp.MustCompile(`@(.+)`)
		return re.FindStringSubmatch(p.currentCommand)[1]
	}

	if p.currentCommandType == LCommand {
		re := regexp.MustCompile(`\((.+)\)`)
		return re.FindStringSubmatch(p.currentCommand)[1]
	}
	return ""
}

// GetDestination returns the destination for a c command,
// empty string if none
func (p *Parser) GetDestination() string {
	if p.currentCommandType == CCommand {
		re := regexp.MustCompile(`(\w+)=.*`)
		match := re.FindStringSubmatch(p.currentCommand)
		if match != nil {
			return match[1]
		}
	}
	return ""
}

// GetComp returns the comp of a c command, empty string if none
func (p *Parser) GetComp() string {
	if p.currentCommandType == CCommand {
		re := regexp.MustCompile(`(.*=)?([^;]*)(;.*)?`)
		match := re.FindAllStringSubmatch(p.currentCommand, -1)
		if match != nil {
			return match[0][2]
		}
	}
	return ""
}

// GetJump returns the jump of a c command, empty string if none
func (p *Parser) GetJump() string {
	if p.currentCommandType == CCommand {
		re := regexp.MustCompile(`(\w+=)?([^;]+)(;)?(\w+)?`)
		match := re.FindStringSubmatch(p.currentCommand)
		if match != nil {
			return match[4]
		}
	}
	return ""
}

func isLineCommand(line string) bool {
	whiteSpace := strings.Replace(line, " ", "", -1)
	whiteSpace = strings.Replace(whiteSpace, "\t", "", -1)
	return whiteSpace != "" && whiteSpace[0:2] != "//"
}

// AddSymbolLineNumber adds the symbol with the address of the current line number
func (p *Parser) AddSymbolLineNumber(symbol string) {
	p.addSymbol(symbol, p.LineNumber)
}

func (p *Parser) addSymbol(symbol string, address int) {
	if !p.ContainsSymbol(symbol) {
		p.symbols[symbol] = address
	}
}

// AddRamSymbol adds the symbol with the next RAM address available
func (p *Parser) AddRamSymbol(symbol string) {
	p.addSymbol(symbol, p.nextAddress)
	p.nextAddress++
}

// ContainsSymbol returns true if the given symbol is in the dict
func (p *Parser) ContainsSymbol(symbol string) bool {
	_, contains := p.symbols[symbol]
	return contains
}

// GetAddress returns the address of a symbol
func (p *Parser) GetAddress(symbol string) int {
	return p.symbols[symbol]
}

// GetSymbolTable returns the symbol table
func (p *Parser) GetSymbolTable() map[string]int {
	return p.symbols
}

func (p *Parser) GetCurrentCommand() string {
	return p.currentCommand
}
