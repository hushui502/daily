package lexer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//var reSpaces = regexp.MustCompile(`^\s+`)
var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")
var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)
var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)
var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)

var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

type Lexer struct {
	chunk string	// 源代码
	chunkName string // 源代码文件号
	line int	// 当前行数
	nextToken string
	nextTokenKind int
	nextTokenLine int
}

func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{chunk, chunkName, 1, "", 0, 0}
}

// 查看下一个token的类型
func (self *Lexer) LookAhead() int {
	// 缓存作用
	if self.nextTokenLine > 0 {
		return self.nextTokenKind
	}
	currentLine := self.line
	line, kind, token := self.NextToken()
	self.line = currentLine
	self.nextTokenLine = line
	self.nextTokenKind = kind
	self.nextToken = token

	return kind
}

// 获取指定类型的token，语法相关
// 游标会前进1个字符
func (self *Lexer) NextTokenOfKind(kind int) (line int, token string) {
	// 实际的下一个token kind
	line, _kind, token := self.NextToken()
	// 实际和预期应该一致，否则就是语法错误
	if kind != _kind {
		self.error("syntax error near '%s'", token)
	}

	return line, token
}

// 提取标识符
func (self *Lexer) NextIdentifier() (line int, token string) {
	return self.NextTokenOfKind(TOKEN_IDENTIFIER)
}

// 返回当前行号
func (self *Lexer) Line() int {
	return self.line
}

// 跳过空白符
func (self *Lexer) skipWhiteSpaces() {
	for len(self.chunk) > 0 {
		if self.test("--") {
			self.skipComment()
		} else if self.test("\r\n") || self.test("\n\r") {
			self.next(2)
			self.line += 1
		} else if isNewLine(self.chunk[0]) {
			self.next(1)
			self.line += 1
		} else if isWhiteSpace(self.chunk[0]) {
			self.next(1)
		} else {
			break
		}
	}
}

// 判断剩余代码是否以某种字符串开头
func (self *Lexer) test(s string) bool {
	return strings.HasPrefix(self.chunk, s)
}

// 代码跳过n个字符
func (self *Lexer) next(n int) {
	self.chunk = self.chunk[n:]
}

// 是否是空白字符
func isWhiteSpace(c byte) bool {
	switch c {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}

	return false
}

// 判断是否是回车符或者换行符
func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

// 跳过注释，这里根据不同的PL定义不同
func (self *Lexer) skipComment() {
	self.next(2)		// skip --

	// long comment
	if self.test("[") {
		self.scanLongString()
		return
	}
	// short comment
	for len(self.chunk) > 0 && !isNewLine(self.chunk[0]) {
		self.next(1)
	}
}

// 下一个token的行号，类型，token实际string
func (self *Lexer) NextToken() (line, kind int, token string) {
	// 如果有缓存直接返回
	if self.nextTokenLine > 0 {
		line = self.nextTokenLine
		kind = self.nextTokenKind
		token = self.nextToken
		// 将当前的缓存清空
		self.line = self.nextTokenLine
		self.nextTokenLine = 0
		return
	}

	switch self.chunk[0] {
	case ';':
		self.next(1)
		return self.line, TOKEN_SEP_SEMI, ";"
	case ',':
		self.next(1)
		return self.line, TOKEN_SEP_COMMA, ","
	case '(':
		self.next(1)
		return self.line, TOKEN_SEP_LPAREN, "("
	case ')':
		self.next(1)
		return self.line, TOKEN_SEP_RPAREN, ")"
	case ']':
		self.next(1)
		return self.line, TOKEN_SEP_RBRACK, "]"
	case '{':
		self.next(1)
		return self.line, TOKEN_SEP_LCURLY, "{"
	case '}':
		self.next(1)
		return self.line, TOKEN_SEP_RCURLY, "}"
	case '+':
		self.next(1)
		return self.line, TOKEN_OP_ADD, "+"
	case '-':
		self.next(1)
		return self.line, TOKEN_OP_MINUS, "-"
	case '*':
		self.next(1)
		return self.line, TOKEN_OP_MUL, "*"
	case '^':
		self.next(1)
		return self.line, TOKEN_OP_POW, "^"
	case '%':
		self.next(1)
		return self.line, TOKEN_OP_MOD, "%"
	case '&':
		self.next(1)
		return self.line, TOKEN_OP_BAND, "&"
	case '|':
		self.next(1)
		return self.line, TOKEN_OP_BOR, "|"
	case '#':
		self.next(1)
		return self.line, TOKEN_OP_LEN, "#"
	case ':':
		if self.test("::") {
			self.next(2)
			return self.line, TOKEN_SEP_LABEL, "::"
		} else {
			self.next(1)
			return self.line, TOKEN_SEP_COLON, ":"
		}
	case '/':
		if self.test("//") {
			self.next(2)
			return self.line, TOKEN_OP_IDIV, "//"
		} else {
			self.next(1)
			return self.line, TOKEN_OP_DIV, "/"
		}
	case '~':
		if self.test("~=") {
			self.next(2)
			return self.line, TOKEN_OP_NE, "~="
		} else {
			self.next(1)
			return self.line, TOKEN_OP_WAVE, "~"
		}
	case '=':
		if self.test("==") {
			self.next(2)
			return self.line, TOKEN_OP_EQ, "=="
		} else {
			self.next(1)
			return self.line, TOKEN_OP_ASSIGN, "="
		}
	case '<':
		if self.test("<<") {
			self.next(2)
			return self.line, TOKEN_OP_SHL, "<<"
		} else if self.test("<=") {
			self.next(2)
			return self.line, TOKEN_OP_LE, "<="
		} else {
			self.next(1)
			return self.line, TOKEN_OP_LT, "<"
		}
	case '>':
		if self.test(">>") {
			self.next(2)
			return self.line, TOKEN_OP_SHR, ">>"
		} else if self.test(">=") {
			self.next(2)
			return self.line, TOKEN_OP_GE, ">="
		} else {
			self.next(1)
			return self.line, TOKEN_OP_GT, ">"
		}
	case '.':
		if self.test("...") {
			self.next(3)
			return self.line, TOKEN_VARARG, "..."
		} else if self.test("..") {
			self.next(2)
			return self.line, TOKEN_OP_CONCAT, ".."
		} else if len(self.chunk) == 1 || !isDigit(self.chunk[1]) {
			self.next(1)
			return self.line, TOKEN_SEP_DOT, "."
		}
	}

	// 数字字面量
	c := self.chunk[0]
	if c == '.' || isDigit(c) {
		token := self.scanNumber()
		return self.line, TOKEN_NUMBER, token
	}

	if c == '_' || isLatter(c) {
		token := self.scanIdentifier()
		if kind, found := keywords[token]; found {
			return self.line, kind, token
		} else {
			return self.line, TOKEN_IDENTIFIER, token
		}
	}

	self.error("unexpected symbol near %q", c)
	return
}

func (self *Lexer) scanLongString() string {
	// 取得 [[
	openingLongBracket := reOpeningLongBracket.FindString(self.chunk)
	if openingLongBracket == "" {
		self.error("invalid long string delimiter near '%s'", self.chunk[0:2])
	}
	// 取得 ]]
	closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
	// 取得 ]]之前的index
	closingLongBracketIdx := strings.Index(self.chunk, closingLongBracket)
	if closingLongBracketIdx < 0 {
		self.error("unfinished long string or comment")
	}

	// 获取[[str]]的str
	str := self.chunk[len(openingLongBracket):closingLongBracketIdx]
	// 跳过这段long string
	self.next(closingLongBracketIdx + len(closingLongBracket))

	// 将这段long str中的换行回车符都统一替换成\n，因为它们的作用都是“换行 line+=1”
	str = reNewLine.ReplaceAllString(str, "\n")
	// 获取行数需要+多少
	self.line += strings.Count(str, "\n")
	// 如果第一个字符就是换行，则跳过获取str
	if len(str) > 0 && str[0] == '\n' {
		str = str[1:]
	}

	return str
}

// 短字符串
func (self *Lexer) scanShortString() string {
	if str := reShortStr.FindString(self.chunk); str != "" {
		self.next(len(str))
		str = str[1 : len(str)-1]
		if strings.Index(str, `\`) >= 0 {
			self.line += len(reNewLine.FindAllString(str, -1))
			str = self.escape(str)
		}
		return str
	}
	self.error("unfinished string")
	return ""
}

// 组装字符串
func (self *Lexer) escape(str string) string {
	var buf bytes.Buffer

	for len(str) > 0 {
		if str[0] != '\\' {
			buf.WriteByte(str[0])
			str = str[1:]
			continue
		}

		if len(str) == 1 {
			self.error("unfinished string")
		}

		switch str[1] {
		case 'a':
			buf.WriteByte('\a')
			str = str[2:]
			continue
		case 'b':
			buf.WriteByte('\b')
			str = str[2:]
			continue
		case 'f':
			buf.WriteByte('\f')
			str = str[2:]
			continue
		case 'n', '\n':
			buf.WriteByte('\n')
			str = str[2:]
			continue
		case 'r':
			buf.WriteByte('\r')
			str = str[2:]
			continue
		case 't':
			buf.WriteByte('\t')
			str = str[2:]
			continue
		case 'v':
			buf.WriteByte('\v')
			str = str[2:]
			continue
		case '"':
			buf.WriteByte('"')
			str = str[2:]
			continue
		case '\'':
			buf.WriteByte('\'')
			str = str[2:]
			continue
		case '\\':
			buf.WriteByte('\\')
			str = str[2:]
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
			if found := reDecEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[1:], 10, 32)
				if d <= 0xFF {
					buf.WriteByte(byte(d))
					str = str[len(found):]
					continue
				}
				self.error("decimal escape too large near '%s'", found)
			}
		case 'x': // \xXX
			if found := reHexEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[2:], 16, 32)
				buf.WriteByte(byte(d))
				str = str[len(found):]
				continue
			}
		case 'u': // \u{XXX}
			if found := reUnicodeEscapeSeq.FindString(str); found != "" {
				d, err := strconv.ParseInt(found[3:len(found)-1], 16, 32)
				if err == nil && d <= 0x10FFFF {
					buf.WriteRune(rune(d))
					str = str[len(found):]
					continue
				}
				self.error("UTF-8 value too large near '%s'", found)
			}
		case 'z':
			str = str[2:]
			for len(str) > 0 && isWhiteSpace(str[0]) { // todo
				str = str[1:]
			}
			continue
		}
		self.error("invalid escape sequence near '\\%c'", str[1])
	}

	return buf.String()
}

// 自定义的错误
func (self *Lexer) error(f string, a ...interface{}) {
	err := fmt.Sprintf(f, a...)
	err = fmt.Sprintf("%s:%d: %s", self.chunkName, self.line, err)

	panic(err)
}

// 是否是数字
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// 提取数字
func (self *Lexer) scanNumber() string {
	return self.scan(reNumber)
}

// 根据不同的正则表达式来匹配出不同的string
func (self *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(self.chunk); token != "" {
		self.next(len(token))
		return token
	}

	panic("unreachable!")
}

// 是否是字符
func isLatter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// 提取标识符
func (self *Lexer) scanIdentifier() string {
	return self.scan(reIdentifier)
}