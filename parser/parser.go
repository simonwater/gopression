package parser

import (
	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/parser/parselet"
	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/values"
)

// 注册前缀表达式处理器
var prefixParselets = map[values.TokenType]parselet.PrefixParselet{
	values.NUMBER:     parselet.NewLiteralParselet(),
	values.STRING:     parselet.NewLiteralParselet(),
	values.IDENTIFIER: parselet.NewIdParselet(),
	values.LEFT_PAREN: parselet.NewGroupParselet(),
	values.MINUS:      parselet.NewPreUnaryParselet(PREC_UNARY),
	values.BANG:       parselet.NewPreUnaryParselet(PREC_UNARY),
	values.IF:         parselet.NewIfParselet(),
}

// 注册中缀表达式处理器
var infixParselets = map[values.TokenType]parselet.InfixParselet{
	values.PLUS:          parselet.NewBinaryParselet(PREC_TERM),
	values.MINUS:         parselet.NewBinaryParselet(PREC_TERM),
	values.PERCENT:       parselet.NewBinaryParselet(PREC_FACTOR),
	values.STAR:          parselet.NewBinaryParselet(PREC_FACTOR),
	values.SLASH:         parselet.NewBinaryParselet(PREC_FACTOR),
	values.STARSTAR:      parselet.NewRightCombineBinaryParselet(PREC_POWER),
	values.EQUAL:         parselet.NewAssignParselet(PREC_ASSIGNMENT),
	values.OR:            parselet.NewLogicParselet(PREC_OR),
	values.AND:           parselet.NewLogicParselet(PREC_AND),
	values.EQUAL_EQUAL:   parselet.NewBinaryParselet(PREC_EQUALITY),
	values.BANG_EQUAL:    parselet.NewBinaryParselet(PREC_EQUALITY),
	values.LESS:          parselet.NewBinaryParselet(PREC_COMPARISON),
	values.LESS_EQUAL:    parselet.NewBinaryParselet(PREC_COMPARISON),
	values.GREATER:       parselet.NewBinaryParselet(PREC_COMPARISON),
	values.GREATER_EQUAL: parselet.NewBinaryParselet(PREC_COMPARISON),
	values.LEFT_PAREN:    parselet.NewCallParselet(PREC_CALL),
	values.DOT:           parselet.NewGetParselet(PREC_CALL),
}

// ================== 解析器实现 ==================
type Parser struct {
	tokens  []values.Token
	current int
}

func NewParser(source string) *Parser {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	p := &Parser{
		tokens:  tokens,
		current: 0,
	}
	return p
}

// ================== 接口方法实现 ==================

// Parse 解析整个表达式
func (p *Parser) Parse() (exprs.Expr, error) {
	result, err := util.SafeExecute(func() exprs.Expr {
		return p.ExpressionPrec(PREC_NONE)
	})
	if p.Peek().Type != values.EOF {
		return nil, parselet.NewLoxParseError(p.Peek(), "unknown token: "+p.Peek().Lexeme)
	}
	return result, err
}

// ExpressionPrec 解析操作符优先级大于minPrec的子表达式
func (p *Parser) ExpressionPrec(minPrec int) exprs.Expr {
	token := p.Advance()
	prefixParselet := prefixParselets[token.Type]
	if prefixParselet == nil {
		panic(parselet.NewLoxParseError(token, "unknown token: "+token.Lexeme))
	}

	lhs := prefixParselet.Parse(p, token)

	for !p.IsAtEnd() {
		next := p.Peek()
		infixParselet := infixParselets[next.Type]
		if infixParselet == nil {
			break
		}

		precedence := infixParselet.GetPrecedence()
		if precedence <= minPrec {
			break
		}

		token = p.Advance()
		lhs = infixParselet.Parse(p, lhs, token)
	}

	return lhs
}

// Match 检查当前token是否匹配任意给定的类型
func (p *Parser) Match(types ...values.TokenType) bool {
	for _, t := range types {
		if p.Check(t) {
			p.Advance()
			return true
		}
	}
	return false
}

// Consume 消费指定类型的token
func (p *Parser) Consume1(expected values.TokenType) values.Token {
	return p.Consume(expected, "Expected token "+expected.String()+" and found "+p.Peek().Type.String())
}

// ConsumeWithMessage 消费指定类型的token，提供自定义错误信息
func (p *Parser) Consume(expected values.TokenType, message string) values.Token {
	if p.Check(expected) {
		return p.Advance()
	}
	panic(parselet.NewLoxParseError(p.Peek(), message))
}

// Advance 前进到下一个token
func (p *Parser) Advance() values.Token {
	if !p.IsAtEnd() {
		p.current++
	}
	return p.Previous()
}

// Check 检查当前token类型
func (p *Parser) Check(tp values.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.Peek().Type == tp
}

// Peek 查看当前token
func (p *Parser) Peek() values.Token {
	return p.tokens[p.current]
}

// Previous 获取前一个token
func (p *Parser) Previous() values.Token {
	return p.tokens[p.current-1]
}

// IsAtEnd 是否到达token流末尾
func (p *Parser) IsAtEnd() bool {
	return p.Peek().Type == values.EOF
}
