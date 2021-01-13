package src

import (
	"fmt"
	"regexp"
	"strings"
)

//SourceCharacter ::=  #x0009 | #x000A | #x000D | [#x0020-#xFFFF]
//Name            ::= [_A-Za-z][_0-9A-Za-z]*
//StringCharacter ::= SourceCharacter - '"'
//String          ::= '"' '"' Ignored | '"' StringCharacter '"' Ignored
//Variable        ::= "$" Name Ignored
//Assignment      ::= Variable Ignored "=" Ignored String Ignored
//Print           ::= "print" "(" Ignored Variable Ignored ")" Ignored
//Statement       ::= Print | Assignment
//SourceCode      ::= Statement+

// token const
const (
	TOKEN_EOF         = iota // end of file
	TOKEN_VAR_PREFIX         // $
	TOKEN_LEFT_PAREN         // (
	TOKEN_RIGHT_PAREN        // )
	TOKEN_EQUAL              // =
	TOKEN_QUOTE              // "
	TOKEN_DUOQUOTE           // ""
	TOKEN_NAME               // Name :: [_A-Za-z][_0-9a-zA-Z]
	TOKEN_PRINT              // print
	TOKEN_IGNORED            // Ignored
)

var tokenNameMap = map[int]string{
	TOKEN_EOF:         "EOF",
	TOKEN_VAR_PREFIX:  "$",
	TOKEN_LEFT_PAREN:  "(",
	TOKEN_RIGHT_PAREN: ")",
	TOKEN_EQUAL:       "=",
	TOKEN_QUOTE:       "\"",
	TOKEN_DUOQUOTE:    "\"\"",
	TOKEN_NAME:        "Name",
	TOKEN_PRINT:       "print",
	TOKEN_IGNORED:     "Ignored",
}

// multiple characters k-v map
var keywords = map[string]int{
	"print": TOKEN_PRINT,
}

// regex match patterns
var regexName = regexp.MustCompile(`^[_\d\w]+`)

type Lexer struct {
	sourceCode       string // 源代码
	lineNum          int    // 执行到代码的当前行数
	nextToken        string // 下一个Token的内容
	nextTokenType    int    // 下一个Token的类型
	nextTokenLineNum int    // 下一个Token的行号
}

func NewLexer(sourceCode string) *Lexer {
	// start at line 1 in default.
	return &Lexer{sourceCode, 1, "", 0, 0}
}

func (lexer *Lexer) scan(regexp *regexp.Regexp) string {
	if token := regexp.FindString(lexer.sourceCode); token != "" {
		lexer.skipSourceCode(len(token))
		return token
	}
	panic("unreachable!")
	return ""
}

func (lexer *Lexer) scanBeforeToken(token string) string {
	s := strings.Split(lexer.sourceCode, token)
	if len(s) < 2 {
		panic("unreached!")
		return ""
	}
	lexer.skipSourceCode(len(s[0]))

	return s[0]
}

func (lexer *Lexer) scanName() string {
	return lexer.scan(regexName)
}

func (lexer *Lexer) nextSourceCodeIs(s string) bool {
	return strings.HasPrefix(lexer.sourceCode, s)
}

func (lexer *Lexer) skipSourceCode(n int) {
	lexer.sourceCode = lexer.sourceCode[n:]
}

func (lexer *Lexer) isIgnored() bool {
	isIgored := false
	// target pattern
	isNewLine := func(c byte) bool {
		return c == '\r' || c == '\n'
	}
	// not need to add line num
	isWhiteSpace := func(c byte) bool {
		switch c {
		case '\t', '\n', '\v', '\r', ' ':
			return true
		}
		return false
	}
	// matching
	for len(lexer.sourceCode) > 0 {
		if lexer.nextSourceCodeIs("\r\n") || lexer.nextSourceCodeIs("\n\r") {
			lexer.skipSourceCode(2)
			lexer.lineNum += 1
			isIgored = true
		} else if isNewLine(lexer.sourceCode[0]) {
			lexer.skipSourceCode(1)
			lexer.lineNum += 1
			isIgored = true
		} else if isWhiteSpace(lexer.sourceCode[0]) {
			lexer.skipSourceCode(1)
			isIgored = true
		} else {
			break
		}
	}

	return isIgored
}

func (lexer *Lexer) GetLineNum() int {
	return lexer.lineNum
}

// 这个函数用于断言下一个 Token 是什么. 并且由于内部执行了 GetNextToken(), 所以游标会自动向前移动
func (lexer *Lexer) NextTokenIs(tokenType int) (lineNum int, token string) {
	nowLineNum, nowTokenType, nowToken := lexer.GetNextToken()
	// syntax error
	if tokenType != nowTokenType {
		err := fmt.Sprintf("NextTokenIs(): syntax error near '%s', expected token: {%s} but got {%s}.", tokenNameMap[nowTokenType], tokenNameMap[tokenType], tokenNameMap[nowTokenType])
		panic(err)
	}
	return nowLineNum, nowToken
}

// 获取下一个token的属性
func (laxer *Lexer) GetNextToken() (lineNum int, tokenType int, token string) {
	// next token already loaded
	if laxer.nextTokenLineNum > 0 {
		lineNum = laxer.nextTokenLineNum
		tokenType = laxer.nextTokenType
		token = laxer.nextToken
		laxer.lineNum = laxer.nextTokenLineNum
		laxer.nextTokenLineNum = 0
		return
	}

	return laxer.MatchToken()
}

// 看下一个token的类型,如果不是预期的则跳过
func (lexer *Lexer) LookAheadAndSkip(expectedType int) {
	// get next token
	nowLineNum := lexer.lineNum
	lineNum, tokenType, token := lexer.GetNextToken()
	// if not is expected type, reverse cursor
	if tokenType != expectedType {
		lexer.lineNum = nowLineNum
		lexer.nextTokenLineNum = lineNum
		lexer.nextTokenType = tokenType
		lexer.nextToken = token
	}
}

//
func (lexer *Lexer) LookAhead() int {
	// lexer.nextToken already setted
	if lexer.nextTokenLineNum > 0 {
		return lexer.nextTokenType
	}

	// set it
	nowLineNum := lexer.lineNum
	lineNum, tokenType, token := lexer.GetNextToken()
	// 游标移动后再移动回来
	lexer.lineNum = nowLineNum
	lexer.nextTokenLineNum = lineNum
	lexer.nextTokenType = tokenType
	lexer.nextToken = token

	return tokenType
}

func (lexer *Lexer) MatchToken() (lineNum int, tokenType int, token string) {
	// check token
	switch lexer.sourceCode[0] {
	case '$':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_VAR_PREFIX, "$"
	case '(':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_LEFT_PAREN, "("
	case ')':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_RIGHT_PAREN, ")"
	case '=':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_EQUAL, "="
	case '"':
		if lexer.nextSourceCodeIs("\"\"") {
			lexer.skipSourceCode(2)
			return lexer.lineNum, TOKEN_DUOQUOTE, "\"\""
		}
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_QUOTE, "\""
	}

	// check multiple characters token
	if lexer.sourceCode[0] == '_' || isLetter(lexer.sourceCode[0]) {
		token := lexer.scanName()
		if tokenType, isMatch := keywords[token]; isMatch {
			return lexer.lineNum, tokenType, token
		} else {
			return lexer.lineNum, TOKEN_NAME, token
		}
	}

	// unexpected symbol
	err := fmt.Sprintf("MatchToken(): unexpected symbol near '%q'.", lexer.sourceCode[0])
	panic(err)
	return
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}


// eg.
// Print ::= "print" "(" Ignored Variable Ignored ")" Ignored
//func parsePrint(lexer *Lexer) (*Print, error) {
//	var print Print
//	var err   error
//
//	print.LineNum = lexer.GetLineNum()
//	lexer.NextTokenIs(TOKEN_PRINT)
//	lexer.NextTokenIs(TOKEN_LEFT_PAREN)
//	lexer.LookAheadAndSkip(TOKEN_IGNORED)
//	if print.Variable, err = parseVariable(lexer); err != nil {
//		return nil, err
//	}
//	lexer.LookAheadAndSkip(TOKEN_IGNORED)
//	lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
//	lexer.LookAheadAndSkip(TOKEN_IGNORED)
//	return &print, nil
//}