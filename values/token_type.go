package values

import "fmt"

type TokenType int

const (
	// Single-character tokens.
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	PERCENT

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	STAR
	STARSTAR
	AND
	OR

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NULL
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	ERROR
	EOF
)

// 标题映射表
var tokenTitles = map[TokenType]string{
	// 单字符标记
	LEFT_PAREN:  "LEFT_PAREN",
	RIGHT_PAREN: "RIGHT_PAREN",
	LEFT_BRACE:  "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",
	COMMA:       "COMMA",
	DOT:         "DOT",
	MINUS:       "MINUS",
	PLUS:        "PLUS",
	SEMICOLON:   "SEMICOLON",
	SLASH:       "SLASH",
	PERCENT:     "PERCENT",

	// 多字符标记
	BANG:          "BANG",
	BANG_EQUAL:    "BANG_EQUAL",
	EQUAL:         "EQUAL",
	EQUAL_EQUAL:   "EQUAL_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",
	STAR:          "STAR",
	STARSTAR:      "STARSTAR",
	AND:           "AND",
	OR:            "OR",

	// 字面量
	IDENTIFIER: "IDENTIFIER",
	STRING:     "STRING",
	NUMBER:     "NUMBER",

	// 关键字
	CLASS:  "CLASS",
	ELSE:   "ELSE",
	FALSE:  "FALSE",
	FUN:    "FUN",
	FOR:    "FOR",
	IF:     "IF",
	NULL:   "NULL",
	PRINT:  "PRINT",
	RETURN: "RETURN",
	SUPER:  "SUPER",
	THIS:   "THIS",
	TRUE:   "TRUE",
	VAR:    "VAR",
	WHILE:  "WHILE",

	// 特殊标记
	ERROR: "ERROR",
	EOF:   "EOF",
}

// Title 返回枚举的标题字符串
func (tt TokenType) Title() string {
	if title, ok := tokenTitles[tt]; ok {
		return title
	}
	panic(fmt.Sprintf("Unknown(%d)", tt))
}

// String 实现Stringer接口，调用Title
func (tt TokenType) String() string {
	return tt.Title()
}
