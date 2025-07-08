package parser_test

import (
	"testing"

	"github.com/simonwater/gopression/ir/exprs"
	"github.com/simonwater/gopression/parser"
	"github.com/stretchr/testify/assert"
)

func parseExpr(src string) exprs.Expr {
	p := parser.NewParser(src)
	return p.Parse()
}

func TestParseNumberLiteral(t *testing.T) {
	expr := parseExpr("123")
	literalExpr, ok := expr.(*exprs.LiteralExpr)
	assert.True(t, ok)
	assert.Equal(t, "123", literalExpr.String())
}

func TestParseStringLiteral(t *testing.T) {
	expr := parseExpr(`"abc"`)
	literalExpr, ok := expr.(*exprs.LiteralExpr)
	assert.True(t, ok)
	assert.Equal(t, "abc", literalExpr.String())
}

func TestParseIdentifier(t *testing.T) {
	expr := parseExpr("foo")
	idExpr, ok := expr.(*exprs.IdExpr)
	assert.True(t, ok)
	assert.Equal(t, "foo", idExpr.String())
}

func TestParseUnaryExpression(t *testing.T) {
	expr := parseExpr("-1")
	unary, ok := expr.(*exprs.UnaryExpr)
	assert.True(t, ok)
	_, ok = unary.Right.(*exprs.LiteralExpr)
	assert.True(t, ok)
}

func TestParseBinaryExpression(t *testing.T) {
	expr := parseExpr("1 + 2")
	binary, ok := expr.(*exprs.BinaryExpr)
	assert.True(t, ok)
	assert.Equal(t, "+", binary.Operator.Lexeme)
}

func TestOperatorPrecedence(t *testing.T) {
	expr := parseExpr("1 + 2 * 3")
	binary, ok := expr.(*exprs.BinaryExpr)
	assert.True(t, ok)
	_, ok = binary.Right.(*exprs.BinaryExpr)
	assert.True(t, ok)
}

func TestParseGroupedExpression(t *testing.T) {
	expr := parseExpr("3 * (1 + 2)")
	binary, ok := expr.(*exprs.BinaryExpr)
	assert.True(t, ok)
	right, ok := binary.Right.(*exprs.BinaryExpr)
	assert.True(t, ok)
	literal1, _ := right.Left.(*exprs.LiteralExpr)
	literal2, _ := right.Right.(*exprs.LiteralExpr)
	assert.Equal(t, "1", literal1.String())
	assert.Equal(t, "2", literal2.String())
}

func TestParseInvalidInput(t *testing.T) {
	assert.Panics(t, func() {
		parseExpr("@")
	})
}
