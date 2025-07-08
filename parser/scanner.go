package parser

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/simonwater/gopression/values"
)

type Scanner struct {
	source  string
	tokens  []values.Token
	start   int
	current int
	line    int
	runes   []rune
}

var keywords = map[string]values.TokenType{
	"class":  values.CLASS,
	"else":   values.ELSE,
	"false":  values.FALSE,
	"for":    values.FOR,
	"fun":    values.FUN,
	"if":     values.IF,
	"null":   values.NULL,
	"print":  values.PRINT,
	"return": values.RETURN,
	"super":  values.SUPER,
	"this":   values.THIS,
	"true":   values.TRUE,
	"var":    values.VAR,
	"while":  values.WHILE,
}

func NewScanner(source string) *Scanner {
	runes := []rune(source)
	return &Scanner{
		source:  source,
		tokens:  []values.Token{},
		start:   0,
		current: 0,
		line:    1,
		runes:   runes,
	}
}

func (s *Scanner) ScanTokens() []values.Token {
	for !s.isEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, *values.NewToken(values.EOF, "", values.NewNullValue(), s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(values.LEFT_PAREN, values.NewNullValue())
	case ')':
		s.addToken(values.RIGHT_PAREN, values.NewNullValue())
	case '{':
		s.addToken(values.LEFT_BRACE, values.NewNullValue())
	case '}':
		s.addToken(values.RIGHT_BRACE, values.NewNullValue())
	case ',':
		s.addToken(values.COMMA, values.NewNullValue())
	case '.':
		s.addToken(values.DOT, values.NewNullValue())
	case '-':
		s.addToken(values.MINUS, values.NewNullValue())
	case '+':
		s.addToken(values.PLUS, values.NewNullValue())
	case '*':
		if s.match('*') {
			s.addToken(values.STARSTAR, values.NewNullValue())
		} else {
			s.addToken(values.STAR, values.NewNullValue())
		}
	case ';':
		s.addToken(values.SEMICOLON, values.NewNullValue())
	case '%':
		s.addToken(values.PERCENT, values.NewNullValue())
	case '!':
		if s.match('=') {
			s.addToken(values.BANG_EQUAL, values.NewNullValue())
		} else {
			s.addToken(values.BANG, values.NewNullValue())
		}
	case '=':
		if s.match('=') {
			s.addToken(values.EQUAL_EQUAL, values.NewNullValue())
		} else {
			s.addToken(values.EQUAL, values.NewNullValue())
		}
	case '>':
		if s.match('=') {
			s.addToken(values.GREATER_EQUAL, values.NewNullValue())
		} else {
			s.addToken(values.GREATER, values.NewNullValue())
		}
	case '<':
		if s.match('=') {
			s.addToken(values.LESS_EQUAL, values.NewNullValue())
		} else {
			s.addToken(values.LESS, values.NewNullValue())
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isEnd() {
				s.advance()
			}
		} else {
			s.addToken(values.SLASH, values.NewNullValue())
		}
	case '|':
		if s.match('|') {
			s.addToken(values.OR, values.NewNullValue())
		} else {
			panic(fmt.Sprintf("line %d: unknown character: %c", s.line, c))
		}
	case '&':
		if s.match('&') {
			s.addToken(values.AND, values.NewNullValue())
		} else {
			panic(fmt.Sprintf("line %d: unknown character: %c", s.line, c))
		}
	case ' ', '\t', '\r':
		// ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.stringToken()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identity()
		} else {
			panic(fmt.Sprintf("line %d: unknown character: %c", s.line, c))
		}
	}
}

func (s *Scanner) stringToken() {
	for s.peek() != '"' && !s.isEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isEnd() {
		panic(fmt.Sprintf("line %d: Unterminated string.", s.line))
	}
	s.advance()
	str := string(s.runes[s.start+1 : s.current-1])
	s.addToken(values.STRING, values.NewStringValue(str))
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}
	isDouble := false
	if s.peek() == '.' && isDigit(s.peekNext()) {
		isDouble = true
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	var v values.Value
	numStr := string(s.runes[s.start:s.current])
	if isDouble {
		num, _ := strconv.ParseFloat(numStr, 64)
		v = values.NewDoubleValue(num)
	} else {
		num, _ := strconv.Atoi(numStr)
		v = values.NewIntValue(int32(num))
	}
	s.addToken(values.NUMBER, v)
}

func (s *Scanner) identity() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := string(s.runes[s.start:s.current])
	typ, ok := keywords[text]
	if !ok {
		typ = values.IDENTIFIER
	}
	s.addToken(typ, values.NewNullValue())
}

func (s *Scanner) isEnd() bool {
	return s.current >= len(s.runes)
}

func (s *Scanner) advance() rune {
	c := s.runes[s.current]
	s.current++
	return c
}

func (s *Scanner) match(expected rune) bool {
	if s.isEnd() {
		return false
	}
	if s.runes[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isEnd() {
		return 0
	}
	return s.runes[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.runes) {
		return 0
	}
	return s.runes[s.current+1]
}

func (s *Scanner) addToken(typ values.TokenType, literal values.Value) {
	text := string(s.runes[s.start:s.current])
	s.tokens = append(s.tokens, *values.NewToken(typ, text, literal, s.line))
}

// 工具函数
func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return unicode.IsLetter(c) || c == '_' || isChineseCharacter(c)
}

func isAlphaNumeric(c rune) bool {
	return isDigit(c) || isAlpha(c)
}

func isChineseCharacter(c rune) bool {
	return (c >= '\u4E00' && c <= '\u9FFF') || (c >= '\u3400' && c <= '\u4DBF')
}
