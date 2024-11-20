package parser

import (
	"github.com/demidshumakher/yaml/internal/lexer"
	"unicode/utf8"
)

type context struct {
	increaseLevel int
}

type node int

type Parser struct {
	//ast    linked_list.LinkedList[node] // TODO
	tokens  []*lexer.Token
	index   int
	context context
}

func initParser(src []rune) *Parser {
	p := &Parser{
		tokens: lexer.Scan(src),
		index:  0,
		context: context{
			increaseLevel: 0,
		},
	}
	return p
}

func castRawStringToString(value string) string {
	res := make([]rune, 0, utf8.RuneCountInString(value))
	for _, r := range value {
		if r == '\n' {
			res = append(res, '\\', 'n')
		} else {
			res = append(res, r)
		}
	}
	return string(res)
}

func (p *Parser) heighAdjustment() {
	for _, val := range p.tokens {
		switch val.Type {
		case lexer.NEW_LINE:
			p.context.increaseLevel = 0
		case lexer.MAP_DELIMITER, lexer.SEQUENCE:
			p.context.increaseLevel++
			val.Level += p.context.increaseLevel
		}
	}
}

func (p *Parser) castTokensToNodes() {

}
