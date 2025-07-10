package parselet

import (
	"fmt"

	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/values"
)

// Parser 解析器接口（需在别处实现）
type IParser interface {
	Parse() (exprs.Expr, error)
	ExpressionPrec(minPrec int) exprs.Expr
	Match(types ...values.TokenType) bool
	Consume(expected values.TokenType, message string) values.Token
	Advance() values.Token
	Check(tp values.TokenType) bool
	Peek() values.Token
	Previous() values.Token
	IsAtEnd() bool
}

// InfixParselet 中缀解析接口
type InfixParselet interface {
	Parse(parser IParser, lhs exprs.Expr, token values.Token) exprs.Expr
	GetPrecedence() int
}

// PrefixParselet 前缀解析接口
type PrefixParselet interface {
	Parse(parser IParser, token values.Token) exprs.Expr
}

type LoxParseError struct {
	Token   values.Token
	Message string
}

func NewLoxParseError(token values.Token, message string) *LoxParseError {
	return &LoxParseError{Token: token, Message: message}
}

func (e *LoxParseError) Error() string {
	return fmt.Sprintf("[line %d] Parse error at '%s': %s",
		e.Token.Line, e.Token.Lexeme, e.Message)
}
