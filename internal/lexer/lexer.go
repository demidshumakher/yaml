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
	isFirstScalarsChar  bool
	isAnchor            bool
	isFlow              int
	previousIndentLevel int
	indentNum           int
	characterChecked    bool
}

type Lexer struct {
	tokens  []*Token
	buffer  []rune
	context context
	index   int
}

var whitespaces = []rune{'\x20', '\x09', '\n', '\r'}

func New(src []rune) *Lexer {
	return &Lexer{
		buffer: src,
		tokens: []*Token{},
		context: context{
			isFirstCharAtLine:  true,
			isFirstScalarsChar: true,
			isAnchor:           false,
			isFlow:             0,
			indentNum:          0,
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
	l.context.characterChecked = false
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

func (l *Lexer) canPeekPrevious() bool {
	return l.index > 1
}

func (l *Lexer) skipChar(n int) {
	l.index += n
}

func (l *Lexer) addToken(token *Token) {
	l.tokens = append(l.tokens, token)
}

func (l *Lexer) scanFlowMapStart() bool {
	//if l.context.isFirstScalarsChar {
	//	l.context.isFlow++
	//	l.context.isFirstScalarsChar = true
	//	l.addToken(flowMapStartToken(l.context.pos.indentLevel))
	//	return true
	//}
	//return false
	l.context.isFlow++
	l.context.isFirstScalarsChar = true
	l.addToken(flowMapStartToken(l.context.pos.indentLevel))
	return true
}

func (l *Lexer) scanFlowMapEnd() bool {
	if l.context.isFlow != 0 {
		l.context.isFirstScalarsChar = false
		l.addToken(flowMapEndToken(l.context.pos.indentLevel))
		l.context.isFlow--
		return true
	}
	return false
}

func (l *Lexer) scanMapDelimiter() bool {
	//if !l.context.isFirstScalarsChar {
	//	return false
	//}
	l.addToken(mapDelimiterToken(l.context.pos.indentLevel))
	l.context.isFirstScalarsChar = true
	return true
}

func (l *Lexer) scanFlowArrayStart() bool {
	l.context.isFlow++
	l.context.isFirstScalarsChar = true
	l.addToken(flowArrayStartToken(l.context.pos.indentLevel))
	return true
}

func (l *Lexer) scanFlowArrayEnd() bool {
	l.context.isFirstScalarsChar = false
	if l.context.isFlow != 0 {
		l.context.isFirstScalarsChar = false
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
		} else {
			n++
		}
	}
	l.context.characterChecked = false

	if l.context.previousIndentLevel == 0 {
		l.context.indentNum = n
	}
	current := 0
	if n != 0 {
		// todo fix
		current = n / l.context.indentNum
	}
	l.context.pos.indentLevel = current
}

func (l *Lexer) scanNewLine() {
	l.context.isFirstCharAtLine = true
	l.context.isFirstScalarsChar = true
	l.context.pos.line++
	l.context.previousIndentLevel = l.context.pos.indentLevel
	l.calculateIndentLevel()
	l.addToken(newLineToken(l.context.pos.indentLevel))
}

func (l *Lexer) scanDocumentEnd() bool {
	if l.context.pos.indentLevel != 0 || !l.context.isFirstCharAtLine || (l.canPeek() && l.peekNextChar() != '.') || (l.canPeekN(1) && l.peek(1) != '.') {
		return false
	}
	l.context.isFirstScalarsChar = false
	l.skipChar(2)
	l.addToken(documentEndToken(l.context.pos.indentLevel))
	return true
}

func (l *Lexer) scanDocumentStart() bool {
	if l.context.pos.indentLevel != 0 || !l.context.isFirstCharAtLine || (l.canPeek() && l.peekNextChar() != '-') || (l.canPeekN(1) && l.peek(1) != '-') {
		return false
	}
	l.context.isFirstScalarsChar = false
	l.skipChar(2)
	l.addToken(documentStartToken(l.context.pos.indentLevel))
	return true
}

func (l *Lexer) scanComment() bool {
	if l.canPeekPrevious() && !slices.Contains(whitespaces, l.peekPreviousChar()) {
		return false
	}
	for l.canPeek() && l.peekNextChar() != '\n' && l.peekNextChar() != '\r' {
		l.skipChar(1)
	}
	return true
}

func (l *Lexer) scanAnchor() {
	var b strings.Builder
	for !slices.Contains(whitespaces, l.peekNextChar()) {
		b.WriteRune(l.nextChar())
	}
	l.context.isFirstScalarsChar = false
	l.addToken(anchorToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scanAlias() {
	var b strings.Builder
	for !slices.Contains(whitespaces, l.peekNextChar()) {
		b.WriteRune(l.nextChar())
	}
	l.context.isFirstScalarsChar = false
	l.addToken(aliasToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scanQuote() {
	quotationMark := l.peekCurrentChar()
	l.context.isFirstScalarsChar = false
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
	l.skipChar(1)
	l.addToken(quoteToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scanSequence() bool {
	if (!l.context.isFirstScalarsChar && !l.context.isFirstCharAtLine) || (l.canPeek() && !slices.Contains(whitespaces, l.peekNextChar())) {
		return false
	}
	l.context.isFirstScalarsChar = true
	l.addToken(sequenceToken(l.context.pos.indentLevel))
	l.context.previousIndentLevel = l.context.pos.indentLevel
	return true
}

func (l *Lexer) scanMapKey() bool {
	if !l.context.isFirstCharAtLine {
		return false
	}
	l.addToken(complexMapKeyToken(l.context.pos.indentLevel))
	return true
}

func (l *Lexer) scanMultilineString() {
	sep := l.peekCurrentChar()
	l.context.isFirstScalarsChar = false
	stringIndentLevel := l.context.pos.indentLevel
	var b strings.Builder
	for {
		c := l.nextChar()

		if c == '\r' {
			c = l.nextChar()
		}
		if c == '\n' {
			t := "\\n"
			if sep == '>' {
				t = " "
			}
			l.context.pos.line++
			l.calculateIndentLevel()
			if l.context.pos.indentLevel > stringIndentLevel {
				b.WriteString(t)
				b.WriteString(strings.Repeat(" ", l.context.pos.indentLevel-stringIndentLevel))
				l.skipChar(stringIndentLevel * l.context.indentNum)
				continue
			}
			break
		}
		b.WriteRune(c)
		l.context.isFirstCharAtLine = false
	}
	l.addToken(multiLineStringToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scanTag() bool {
	if !l.context.isFirstScalarsChar && !l.context.isFirstCharAtLine {
		return false
	}
	l.context.isFirstScalarsChar = false
	var b strings.Builder
	for !slices.Contains(whitespaces, l.peekNextChar()) {
		b.WriteRune(l.nextChar())
	}
	l.addToken(tagToken(l.context.pos.indentLevel, b.String()))
	return true
}

func (l *Lexer) scanFlowDelimiter() bool {
	if l.context.isFlow != 0 {
		l.context.isFirstScalarsChar = false
		l.addToken(flowDelimiterToken(l.context.pos.indentLevel))
		return true
	}
	return false
}

var alp = []rune{':', '{', '}', '[', ']', '#', '\n', '\r', ',', '.', '-'}

func (l *Lexer) scanScalar() {
	var b strings.Builder
	for !slices.Contains(alp, l.peekCurrentChar()) || l.context.characterChecked {
		b.WriteRune(l.peekCurrentChar())
		if !l.canPeekN(1) {
			l.addToken(scalarToken(l.context.pos.indentLevel, b.String()))
			return
		}
		l.nextChar()
	}
	l.skipChar(-1)
	if len(strings.TrimSpace(b.String())) == 0 {
		return
	}
	l.addToken(scalarToken(l.context.pos.indentLevel, b.String()))
}

func (l *Lexer) scan() {
Loop:
	for l.canPeek() == true {
		c := l.nextChar()
		switch c {
		case '\x20', '\x09':
			continue
		case ':':
			if l.scanMapDelimiter() {
				continue
			}
		case '{':
			if l.scanFlowMapStart() {
				continue
			}
		case '}':
			if l.scanFlowMapEnd() {
				continue
			}
		case '[':
			if l.scanFlowArrayStart() {
				continue
			}
		case ']':
			if l.scanFlowArrayEnd() {
				continue
			}
		case ',':
			if l.scanFlowDelimiter() {
				continue
			}
		case '.':
			if l.scanDocumentEnd() {
				break Loop
			}
		case '-':
			if l.scanDocumentStart() {
				continue
			}
			if l.scanSequence() {
				continue
			}
		case '#':
			if l.scanComment() {
				continue
			}
		case '\n', '\r':
			l.scanNewLine()
			continue
		case '|', '>':
			l.scanMultilineString()
		case '!':
			if l.scanTag() {
				continue
			}
		case '?':
			if l.scanMapKey() {
				continue
			}
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
		l.context.characterChecked = true
		l.scanScalar()

	}
}

func (l *Lexer) increaseLevel() {
	inc := 0
	l.context.isFlow = 0
	for i, el := range l.tokens {
		switch el.Type {
		case NEW_LINE:
			inc = 0
		case FLOW_MAP_START, FLOW_ARRAY_START:
			inc++
			l.context.isFlow++
			el.Level--
		case FLOW_MAP_END, FLOW_ARRAY_END:
			inc--
			l.context.isFlow--
		case SEQUENCE, COMMPLEX_MAP_KEY:
			inc++
			el.Level--
		case MAP_DELIMITER:
			if l.context.isFlow == 0 {
				inc++
				el.Level--
			} else if i+1 < len(l.tokens) {
				l.tokens[i+1].Level++
			}
		}
		if el.Value == "map" {
			fmt.Println(inc)
		}
		el.Level += inc

	}
}

func (l *Lexer) mergeScalars() {
	new_tokens := make([]*Token, 0, len(l.tokens))
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i].Type == SCALAR {
			current_scalar := l.tokens[i]
			if i+1 < len(l.tokens) {
				for nx := l.tokens[i+1]; nx.Type == SCALAR && nx.Level == current_scalar.Level; nx = l.tokens[i+1] {
					current_scalar.Value += nx.Value
					i++
					if i+1 >= len(l.tokens) {
						break
					}
				}
			}
			new_tokens = append(new_tokens, current_scalar)
		} else if l.tokens[i].Type != FLOW_DELIMITER {
			new_tokens = append(new_tokens, l.tokens[i])
		}
	}
	l.tokens = new_tokens
}

func (l *Lexer) mergeQuotes() {
	new_tokens := make([]*Token, 0, len(l.tokens))
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i].Type == QUOTE {
			current_scalar := l.tokens[i]
			if i+1 < len(l.tokens) {
				for nx := l.tokens[i+1]; nx.Type == QUOTE && nx.Level == current_scalar.Level; nx = l.tokens[i+1] {
					current_scalar.Value += nx.Value
					i++
					if i+1 >= len(l.tokens) {
						break
					}
				}
			}
			new_tokens = append(new_tokens, current_scalar)
		} else if l.tokens[i].Type != FLOW_DELIMITER {
			new_tokens = append(new_tokens, l.tokens[i])
		}
	}
	l.tokens = new_tokens
}

func Scan(src []rune) []*Token {
	lx := New(src)
	lx.scan()
	lx.mergeScalars()
	lx.mergeQuotes()
	lx.increaseLevel()
	return lx.tokens
}
