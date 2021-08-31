package expr

import (
	"errors"
	"fmt"
)

// Expr represents a parsed expression
type Expr struct {
	op    string
	left  *Expr
	right *Expr
	ident string
	num   int
}

func (e *Expr) String() string {
	if e == nil {
		return ""
	}
	if e.op == "" {
		if e.left != nil {
			return e.left.String()
		}
		if e.ident != "" {
			return e.ident
		}
		return fmt.Sprint(e.num)
	}

	left := e.left.String()
	right := e.right.String()
	if left == "" {
		return fmt.Sprintf("(%s%s)", e.op, right)
	}

	return fmt.Sprintf("(%s %s %s)", left, e.op, right)
}

const eof = 0

type parser struct {
	s string
}

func (p *parser) next(doSkip bool) byte {
	if doSkip {
		p.skip()
	}
	if p.s == "" {
		return eof
	}
	c := p.s[0]
	p.s = p.s[1:]

	return c
}

func (p *parser) peek(doSkip bool) byte {
	if doSkip {
		p.skip()
	}
	if p.s == "" {
		return eof
	}

	return p.s[0]
}

func (p *parser) skip() {
	for p.s != "" && p.starts(" \t\n\r") {
		p.s = p.s[1:]
	}
}

func (p *parser) starts(set string) bool {
	if len(p.s) < 1 {
		return false
	}
	c := p.s[0]
	for i := 0; i < len(set); i++ {
		if c == set[i] {
			return true
		}
	}

	return false
}

func (p *parser) nextOpLen() int {
	if p.s == "" {
		return 0
	}
	switch p.s[0] {
	case '+', '-', '*', '/', '%', '^':
		return 1
	case '<':
		return p.maybe("=<")
	case '>':
		return p.maybe("=>")
	case '&':
		return p.maybe("&^")
	case '|':
		return p.maybe("|")
	case '!':
		return p.maybe("=")
	case '=': // = is not an operator but == is.
		if len(p.s) >= 2 || p.s[1] == '=' {
			return 2
		}
	}

	return 0
}

// We are at an operator. Does it extend to the second character?
// Return the length of operator corresponding to the present
// character (assumed), plus possibly a character from extra.
func (p *parser) maybe(extra string) int {
	if len(p.s) >= 2 {
		for i := 0; i < len(extra); i++ {
			if p.s[1] == extra[i] {
				return 2
			}
		}
	}

	return 1
}

func (p *parser) op(singles, doubles string) string {
	n := p.nextOpLen()
	op := p.s[:n]
	switch n {
	case 1:
		for i := 0; i < len(singles); i++ {
			if op[0] == singles[i] {
				p.s = p.s[n:]
				return op
			}
		}
	case 2:
		for i := 0; i < len(doubles); i += 2 {
			if op == doubles[i:i+2] {
				p.s = p.s[n:]
				return op
			}
		}
	}

	return ""
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isAlpha(c byte, digitOK bool) bool {
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' {
		return true
	}
	return digitOK && isDigit(c)
}

func recoverer(errp *error) {
	if r := recover(); r != nil {
		var ok bool
		*errp, ok = r.(error)
		if ok {
			return
		}
		panic(errp)
	}
}

func Parse(s string) (expr *Expr, err error) {
	p := &parser{s: s}
	defer recoverer(&err)
	expr = orList(p)
	if p.peek(true) != eof {
		throw("syntax error at ", p.remaining())
	}

	return
}

func throw(s ...interface{}) {
	panic(errors.New(fmt.Sprint(s...)))
}

func (p *parser) remaining() string {
	if p.s != "" {
		return fmt.Sprintf("%q", p.s)
	}

	return "eof"
}

func (p *parser) parse(singles, doubles string, nextLevel func(*parser) *Expr) *Expr {
	e := nextLevel(p)
	for {
		if p.peek(true) == eof {
			return e
		}
		op := p.op(singles, doubles)
		if op == "" {
			return e
		}
		e = &Expr{
			op:    op,
			left:  e,
			right: nextLevel(p),
		}
	}
}

// orlist = andList | andList '||' orList.
func orList(p *parser) *Expr {
	return p.parse("", "||", andList)
}

// andlist = cmpList | cmpList '&&' andList.
func andList(p *parser) *Expr {
	return p.parse("", "&&", cmpList)
}

// cmpList = expr | expr ('>' | '<' | '==' | '!=' | '>=' | '<=') expr.
func cmpList(p *parser) *Expr {
	return p.parse("+-|^!><", "==!=>=<=", expr)
}

// expr = term | term ('+' | '-' | '|' | '^') term.
func expr(p *parser) *Expr {
	return p.parse("+-|^!", "", term)
}

// term = factor | factor ('*' | '/' | '%' | '>>' | '<<' | '&' | '&^') factor
func term(p *parser) *Expr {
	return p.parse("*/%&", ">><<&^", factor)
}

// factor = constant | identifier | '+' factor | '-' factor | '^' factor | '!' factor | '(' orList ')'
func factor(p *parser) *Expr {
	c := p.peek(true)
	switch {
	case c == eof:
		throw("unexpected eof")
	case isDigit(c):
		return &Expr{
			num: p.number(),
		}
	case isAlpha(c, false):
		return &Expr{
			ident: p.identifier(),
		}
	case p.starts("+-^!"):
		op := p.s[:1]
		p.next(false)
		return &Expr{
			op:    op,
			right: factor(p),
		}
	case c == '(':
		p.next(false)
		e := orList(p)
		if p.next(true) != ')' {
			throw("unclosed paren at ", p.remaining())
		}
		return e
	}
	throw("bad expression at ", p.remaining())
	return nil
}

func (p *parser) number() int {
	n := 0
	for {
		c := p.peek(false)
		if !isDigit(c) {
			break
		}
		p.next(false)
		n = 10*n + int(c) - '0'
	}
	return n
}

func (p *parser) identifier() string {
	s := ""
	for {
		c := p.peek(false)
		if !isAlpha(c, s != "") {
			break
		}
		p.next(false)
		s = s + string(c)
	}
	return s
}

// ErrorMode specifies how to handle arithmetic errors such as division by zero or
// undefined variable: Either return an error (ReturnError) or replace the
// erroneous calculation with zero and press on (ReturnZero).
type ErrorMode int

const (
	ReturnError ErrorMode = iota
	ReturnZero
)

func (e ErrorMode) error(s ...interface{}) int {
	switch e {
	case ReturnError:
		throw(s...)
	case ReturnZero:
		return 0
	}
	panic("bad error mode")
}

// Eval evaluates the expression.
func (e *Expr) Eval(vars map[string]int, errMode ErrorMode) (result int, err error) {
	defer recoverer(&err)
	result = e.eval(vars, errMode)
	return
}

func (e *Expr) eval(vars map[string]int, errMode ErrorMode) int {
	if e == nil {
		return 0
	}
	if e.op == "" {
		if e.ident != "" {
			n, ok := vars[e.ident]
			if !ok {
				return errMode.error("undefined variable ", e.ident)
			}
			return n
		}
		return e.num
	}
	// Binary operators.
	if e.left != nil && e.right != nil {
		left := e.left.eval(vars, errMode)
		right := e.right.eval(vars, errMode)
		switch e.op {
		case "+":
			return left + right
		case "-":
			return left - right
		case "*":
			return left * right
		case "/":
			if right == 0 {
				return errMode.error("division by zero")
			}
			return left / right
		case "%":
			if right == 0 {
				return errMode.error("modulo by zero")
			}
			return left % right
		case "&":
			return left & right
		case "|":
			return left | right
		case "^":
			return left ^ right
		case "&^":
			return left &^ right
		case ">>":
			if right < 0 {
				return errMode.error("negative right shift amount")
			}
			return left >> right
		case "<<":
			if right < 0 {
				return errMode.error("negative left shift amount")
			}
			return left << right
		case "==":
			return toInt(left == right)
		case "!=":
			return toInt(left != right)
		case ">=":
			return toInt(left >= right)
		case "<=":
			return toInt(left <= right)
		case "<":
			return toInt(left < right)
		case ">":
			return toInt(left > right)
		case "||":
			return toInt(left != 0 || right != 0)
		case "&&":
			return toInt(left != 0 && right != 0)
		default:
			throw("unknown binary operator (can't happen) ", e.op)
		}
	}
	if e.right != nil {
		right := e.right.eval(vars, errMode)
		switch e.op {
		case "+":
			return right
		case "-":
			return -right
		case "^":
			return ^right
		case "!":
			return toInt(right == 0)
		default:
			throw("unknown unary operator (can't happen) ", e.op)
		}
	}
	throw("unrecognized expression: can't happen")
	panic("not reached")
}

func toInt(t bool) int {
	if t {
		return 1
	}
	return 0
}
