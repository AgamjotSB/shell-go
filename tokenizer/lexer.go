package tokenizer

import (
	"fmt"
	"strings"
)

const eof = rune(0)

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		pos:   0,
	}
}

func (l *Lexer) NextToken() (token string, isEOF bool, err error) {
	l.skipWhitespace()

	if l.peek() == eof {
		return "", true, nil
	}

	var sb strings.Builder
	inSingleQuote := false

	for {
		char := l.peek()
		if char == eof {
			if inSingleQuote {
				return "", false, fmt.Errorf("error: unclosed single quote")
			}
		}

		if !inSingleQuote && (char == ' ' || char == '\t' || char == '\n') {
			break
		}

		l.advance()

		// toggle but don't write the quote to the result
		if char == '\'' {
			inSingleQuote = !inSingleQuote
			continue
		}

		sb.WriteRune(char)
	}

	return sb.String(), false, nil
}

func (l *Lexer) skipWhitespace() {
	for {
		char := l.peek()
		if char == ' ' || char == '\t' || char == '\n' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return eof
	}

	return l.input[l.pos]
}

func (l *Lexer) advance() rune {
	char := l.peek()
	if char != eof {
		l.pos++
	}
	return char
}
