package main

import (
	"fmt"
	"github.com/kooixh/bolehlang/interpreter"
	"github.com/kooixh/bolehlang/lexer"
	"github.com/kooixh/bolehlang/parser"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	arg := os.Args[1]

	if arg == "--version" || arg == "-v" {
		fmt.Printf("bolehlang v%s\n", version)
		return
	}

	if arg == "--help" || arg == "-h" {
		printUsage()
		return
	}

	source, err := os.ReadFile(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Alamak! Cannot read file '%s': %s\n", arg, err)
		os.Exit(1)
	}

	if err := run(string(source), arg); err != nil {
		fmt.Fprintf(os.Stderr, "Mampus! %s\n", err)
		os.Exit(1)
	}
}

func run(source string, filePath string) error {
	// Lexing
	lex := lexer.New(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		return fmt.Errorf("Lexer error: %s", err)
	}

	// Parsing
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		return fmt.Errorf("Parser error: %s", err)
	}

	// Interpreting
	interp := interpreter.New(filePath)
	if err := interp.Run(program); err != nil {
		return fmt.Errorf("Runtime error: %s", err)
	}

	return nil
}

func printUsage() {
	fmt.Println("🇲🇾 bolehlang — Malaysian English programming language")
	fmt.Println()
	fmt.Println("Usage: bolehlang <file.myl>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help      Show this help")
	fmt.Println("  -v, --version   Show version")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  bolehlang hello.myl")
}
