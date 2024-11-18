package lexer

import (
	"fmt"
	"slices"
	"strings"
)

type position struct {
	line        int
	col         int
	indentLevel int
}

type context struct {
	pos                 position
	isFirstCharAtLine   bool
	isAnchor            bool
	isFlow              int
	previousChar        rune
	previousIndentLevel int
	indentNum           int
}

type Lexer struct {
	Tokens  []*Token
	buffer  []rune
	context context
	index   int
}

var whitespaces = []rune{'\x20', '\x09', '\n', '\r'}

func New(src []rune) *Lexer {
	return &Lexer{
		buffer: src,
		Tokens: []*Token{},
		context: context{
			isFirstCharAtLine: false,
			isAnchor:          false,
			isFlow:            0,
			previousChar:      0,
			indentNum:         0,
			pos: position{
				line:        0,
				col:         0,
				indentLevel: 0,
			},
		},
	}
}

func (l *Lexer) nextChar() rune {
	l.index++
	return l.buffer[l.index-1]
}

func (l *Lexer) canPeek() bool {
	return l.canPeekN(0)
}

func (l *Lexer) canPeekN(n int) bool {
	return l.index+n < len(l.buffer)
}

func (l *Lexer) peek(n int) rune {
	return l.buffer[l.index+n]
}

func (l *Lexer) peekNextChar() rune {
	return l.peek(0)
}

func (l *Lexer) peekCurrentChar() rune {
	return l.peek(-1)
}

func (l *Lexer) peekPreviousChar() rune {
	return l.peek(-2)
}

func (l *Lexer) skipChar(n int) {
	l.index += n
}

func (l *Lexer) addToken(token *Token) {
	l.Tokens = append(l.Tokens, token)
}

func (l *Lexer) scanFlowMapStart() {
	l.context.isFlow++
	l.addToken(flowMapStartToken(l.context.pos.indentLevel))
}

func (l *Lexer) scanFlowMapEnd() bool {
	if l.context.isFlow != 0 {
		l.addToken(flowMapEndToken(l.context.pos.indentLevel))
		l.context.isFlow--
		return true
	}
	return false
}

func (l *Lexer) scanMapDelimiter() {
	l.addToken(mapDelimiterToken(l.context.pos.indentLevel))
}

func (l *Lexer) scanFlowArrayStart() {
	l.context.isFlow++
	l.addToken(flowArrayStartToken(l.context.pos.indentLevel))
}

func (l *Lexer) scanFlowArrayEnd() bool {
	if l.context.isFlow != 0 {
		l.context.isFlow--
		l.addToken(flowArrayEndToken(l.context.pos.indentLevel))
		return true
	}
	return false
}

func (l *Lexer) calculateIndentLevel() {
	n := 0
	for l.canPeekN(n) && slices.Contains(whitespaces, l.peek(n)) {
		if l.peek(n) == '\r' || l.peek(n) == '\n' {
			l.skipChar(n + 1)
			l.calculateIndentLevel()
		}
		n++
	}

	if l.context.previousIndentLevel == 0 {
		l.context.indentNum = n
	}
	fmt.Println(l.context.previousIndentLevel, l.context.indentNum, n)
	current := 0
	if n != 0 {
		// todo fix
		current = n / l.context.indentNum
	}
	l.context.pos.indentLevel = current
}

func (l *Lexer) scanNewLine() {
	l.context.isFirstCharAtLine = true
	l.context.pos.line++
	l.context.previousIndentLevel = l.context.pos.indentLevel
	l.calculateIndentLevel()
}

func (l *Lexer) scanDocumentEnd() {
	l.skipChar(2)
	l.addToken(documentEndToken(l.context.pos.indentLevel))
}

func (l *Lexer) scanDocumentStart() bool {
	if l.peekNextChar() != '-' || l.peek(1) != '-' {
		return false
	}
	l.skipChar(2)
	l.addToken(documentStartToken(l.context.pos.indentLevel))
	return true
}

func (l *Lexer) scanComment() bool {
	if !slices.Contains(whitespaces, l.peekPreviousChar()) {
		return false
	}
	for l.peekNextChar() != '\n' || l.peekNextChar() != '\r' {
		l.skipChar(1)
	}
	return true
}

func (l *Lexer) scanAnchor() {
	var b strings.Builder
	for !slices.Contains(whitespaces, l.peekNextChar()) {
		b.WriteRune(l.nextChar())
	}
	l.addToken(anchorToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scanAlias() {
	var b strings.Builder
	for !slices.Contains(whitespaces, l.peekNextChar()) {
		b.WriteRune(l.nextChar())
	}
	l.addToken(aliasToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scanQuote() {
	quotationMark := l.peekCurrentChar()
	var b strings.Builder
	for l.peekNextChar() != quotationMark {
		b.WriteRune(l.nextChar())
		if l.peekCurrentChar() == '\\' {
			b.WriteRune(l.nextChar())
		}
		if l.peekCurrentChar() == '\n' {
			l.context.pos.line++
			currentIndentLevel := l.context.pos.indentLevel
			l.calculateIndentLevel()
			l.skipChar(currentIndentLevel * l.context.indentNum)
		}

	}

	l.addToken(quoteToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scanSequence() {
	l.addToken(sequenceToken(l.context.pos.indentLevel))
	l.context.previousIndentLevel = l.context.pos.indentLevel
}

func (l *Lexer) scanMapKey() {
	l.addToken(complexMapKeyToken(l.context.pos.indentLevel))
}

func (l *Lexer) scanMultilineString() {
	sep := l.peekCurrentChar()
	stringIndentLevel := l.context.pos.indentLevel
	var b strings.Builder
	for {
		c := l.nextChar()
		if c == '\r' {
			c = l.nextChar()
		}
		if c == '\n' {
			var t rune = '\n'
			if sep == '>' {
				t = ' '
			}
			l.context.pos.line++
			l.calculateIndentLevel()
			if l.context.pos.indentLevel > stringIndentLevel {
				b.WriteRune(t)
				b.WriteString(strings.Repeat(" ", l.context.pos.indentLevel-stringIndentLevel))
				l.skipChar(stringIndentLevel * l.context.indentNum)
				continue
			}
			break
		}
		b.WriteRune(c)
	}
	l.addToken(multiLineStringToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scanTag() {
	var b strings.Builder
	for !slices.Contains(whitespaces, l.peekNextChar()) {
		b.WriteRune(l.nextChar())
	}
	l.addToken(tagToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scanFlowDelimiter() {
	l.addToken(flowDelimiterToken(l.context.pos.indentLevel))
}

var alp = []rune{':', '{', '}', '[', ']', '#', '\n', '\r'}

func (l *Lexer) scanScalar() {
	var b strings.Builder
	for !slices.Contains(alp, l.peekCurrentChar()) {
		b.WriteRune(l.peekCurrentChar())
		l.nextChar()
	}
	l.skipChar(-1)
	l.addToken(scalarToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scan() {
	for l.canPeek() == true {
		c := l.nextChar()
		switch c {
		case '\x20', '\x09':
			continue
		case ':':
			l.scanMapDelimiter()
			continue
		case '{':
			l.scanFlowMapStart()
			continue
		case '}':
			if l.scanFlowMapEnd() {
				continue
			}
		case '[':
			l.scanFlowArrayStart()
			continue
		case ']':
			if l.scanFlowArrayEnd() {
				continue
			}
		case ',':
			l.scanFlowDelimiter()
			continue
		case '.':
			l.scanDocumentEnd()
			continue
		case '-':
			if l.scanDocumentStart() {
				continue
			}
			l.scanSequence()
			continue
		case '#':
			if l.scanComment() {
				continue
			}
		case '\n', '\r':
			l.scanNewLine()
			continue
		case '<':
			// todo
		case '|', '>':
			l.scanMultilineString()
		case '!':
			l.scanTag()
		case '?':
			l.scanMapKey()
			continue
		case '&':
			l.scanAnchor()
			continue
		case '*':
			l.scanAlias()
			continue
		case '\'', '"':
			l.scanQuote()
			continue

		}

		l.scanScalar()

		l.context.previousChar = c
	}
}

func Scan(src []rune) []*Token {
	lx := New(src)
	lx.scan()
	return lx.Tokens
}
