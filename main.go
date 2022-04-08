package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage(os.Args[0])
		os.Exit(1)
	}

	var input string

	if os.Args[1] == "run" {
		if len(os.Args) < 3 {
			fmt.Println("Please provide a file to read")
			os.Exit(1)
		}
		file, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			input += scanner.Text()
		}
	}else if os.Args[1] == "str" {
		if len(os.Args) < 3 {
			fmt.Println("Please provide a string to read")
			os.Exit(1)
		}
		input = os.Args[2]
	}else {
		usage(os.Args[0])
		os.Exit(1)
	}

	tokens := lexer(strings.Split(input, ""))
	run(tokens)
}

func usage(name string) {
	fmt.Printf("%s [run|str] [input]\n", name)
	fmt.Printf("  %s run [file]\n", name)
	fmt.Printf("    Runs the brainfuck program from the given file\n")
	fmt.Printf("  %s str [string]\n", name)
	fmt.Printf("    Runs the brainfuck program from the given string\n")
}

type TokenType int
const (
	TK_RIGHT TokenType = iota
	TK_LEFT
	TK_INCREMENT
	TK_DECREMENT
	TK_LOOP_START
	TK_LOOP_END
	TK_OUTPUT
	TK_INPUT
)

type Token struct {
	Type TokenType
}

func lexer(input []string) []Token {
	var res []Token
	for _, v := range input {
		switch v {
		case ">":
			res = append(res, Token{Type: TK_RIGHT})
		case "<":
			res = append(res, Token{Type: TK_LEFT})
		case "+":
			res = append(res, Token{Type: TK_INCREMENT})
		case "-":
			res = append(res, Token{Type: TK_DECREMENT})
		case "[":
			res = append(res, Token{Type: TK_LOOP_START})
		case "]":
			res = append(res, Token{Type: TK_LOOP_END})
		case ".":
			res = append(res, Token{Type: TK_OUTPUT})
		case ",":
			res = append(res, Token{Type: TK_INPUT})
		}
	}
	return res
}

func run(tokens []Token) {
	var tape map[int]uint8 = map[int]uint8{0: 0}
	var pointer int = 0
	var paren int = 0

	for i := 0; i < len(tokens); {
		curr := tokens[i]
		switch curr.Type {
		case TK_RIGHT:
			pointer++
			i++
		case TK_LEFT:
			pointer--
			i++
		case TK_INCREMENT:
			tape[pointer]++
			i++
		case TK_DECREMENT:
			tape[pointer]--
			i++
		case TK_LOOP_START:
			if tape[pointer] == 0 {
				paren++
				for {
					i++
					if i >= len(tokens) {
						fmt.Println("Loop end not found")
						os.Exit(1)
					}
					if tokens[i].Type == TK_LOOP_START {
						paren++
					} else if tokens[i].Type == TK_LOOP_END {
						paren--
					}
					if paren == 0 {
						break
					}
				}
			} else {
				i++
			}
		case TK_LOOP_END:
			if tape[pointer] != 0 {
				paren--
				for {
					i--
					if i < 0 {
						fmt.Println("Loop start not found")
						os.Exit(1)
					}
					if tokens[i].Type == TK_LOOP_START {
						paren++
					} else if tokens[i].Type == TK_LOOP_END {
						paren--
					}
					if paren == 0 {
						break
					}
				}
			} else {
				i++
			}
		case TK_OUTPUT:
			fmt.Printf("%c", tape[pointer])
			i++
		case TK_INPUT:
			var input uint8
			fmt.Scanf("%c", &input)
			tape[pointer] = input
			i++
		}
	}
}
