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

	for {
		char := l.advance()
		switch char {

		case eof, ' ', '\t', '\n':
			return sb.String(), false, nil

		case '\'':
			token, err := l.handleSingleQuote()
			if err != nil {
				return "", false, err
			}
			sb.WriteString(token)
		case '"':
			token, err := l.handleDoubleQuote()
			if err != nil {
				return "", false, err
			}
			sb.WriteString(token)
		case '\\':
			token := l.handleEscapeBackSlash()
			sb.WriteRune(token)

		default:
			sb.WriteRune(char)
		}
	}
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

func (l *Lexer) handleSingleQuote() (string, error) {
	var sb strings.Builder
	for {
		char := l.advance()
		if char == eof {
			return "", fmt.Errorf("error: unclosed double quote")
		}
		if char == '\'' {
			return sb.String(), nil
		}
		sb.WriteRune(char)
	}
}

func (l *Lexer) handleEscapeBackSlash() rune {
	char := l.advance()
	return char
}

func (l *Lexer) handleDoubleQuote() (string, error) {
	var sb strings.Builder
	for {
		char := l.advance()
		switch char {
		case eof:
			return "", fmt.Errorf("error: unclosed single quote")
		case '"':
			return sb.String(), nil
		case '\\':
			next := l.peek()
			if next == '"' || next == '\\' || next == '$' || next == '`' {
				char = l.handleEscapeBackSlash()
			}
		}
		sb.WriteRune(char)
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
