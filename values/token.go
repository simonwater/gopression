package values

import (
	"fmt"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal Value // 或 *values.Value，取决于你的Value定义
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal Value, line int) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %s %v", t.Type, t.Lexeme, t.Literal)
}
