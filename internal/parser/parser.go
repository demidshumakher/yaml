package parser

import (
	"fmt"
	"github.com/demidshumakher/yaml/internal/ast"
	"github.com/demidshumakher/yaml/internal/lexer"
	"github.com/demidshumakher/yaml/pkg/linked_list"
)

type context struct {
}

type node int

type Parser struct {
	tokens  []*lexer.Token
	index   int
	context context
	nodes   *ast.AST
	anchors map[string]*linked_list.LinkedListNode[ast.NodeValue]
}

func newParser(src []*lexer.Token) *Parser {
	p := &Parser{
		tokens:  src,
		index:   0,
		context: context{},
		nodes:   ast.NewAST(),
		anchors: make(map[string]*linked_list.LinkedListNode[ast.NodeValue]),
	}
	return p
}

func (p *Parser) parse() *ast.AST {
	if len(p.tokens) == 0 {
		return nil
	}
	//var prevsc *lexer.Token = nil
	for i := 0; i < len(p.tokens); i++ {
		el := p.tokens[i]
		// todo complex map key пока что пофиг, их все равно в json не запихнуть
		//fmt.Println(el)
		switch el.Type {
		case lexer.QUOTE, lexer.MULTILINE_STRING:
			if i+1 < len(p.tokens) && p.tokens[i+1].Type == lexer.MAP_DELIMITER {
				p.nodes.Add(ast.NewNodeValue(ast.MAP_KEY, el.Value))
			} else {
				p.nodes.Add(ast.NewNodeValue(ast.STRING, el.Value))
			}
		case lexer.SCALAR:
			if i+1 < len(p.tokens) && p.tokens[i+1].Type == lexer.MAP_DELIMITER {
				p.nodes.Add(ast.NewNodeValue(ast.MAP_KEY, el.Value))
			} else {
				p.nodes.Add(ast.NewNodeValue(ast.SCALAR, el.Value))
			}
		case lexer.MAP_DELIMITER:
			p.nodes.Add(ast.NewNodeValue(ast.MAP_VALUE, ""))
			level := el.Level
			tempArr := []*lexer.Token{}
			for i+1 < len(p.tokens) && p.tokens[i+1].Level > level {
				i++
				tempArr = append(tempArr, p.tokens[i])
			}
			chAst := newParser(tempArr).parse()
			p.nodes.AddChild(chAst)
		case lexer.FLOW_MAP_START:
			p.nodes.Add(ast.NewNodeValue(ast.MAP, ""))
			t := []*lexer.Token{}
			cnt := 1
			for i++; i < len(p.tokens) && (p.tokens[i].Type != lexer.FLOW_MAP_END || cnt != 1); i++ {
				t = append(t, p.tokens[i])
				if p.tokens[i].Type == lexer.FLOW_MAP_START {
					cnt++
				}
				if p.tokens[i].Type == lexer.FLOW_MAP_END {
					cnt--
				}
			}
			chAst := newParser(t).parse()
			p.nodes.AddChild(chAst)

		case lexer.FLOW_ARRAY_START:
			p.nodes.Add(ast.NewNodeValue(ast.ARRAY, ""))
			t := []*lexer.Token{}
			cnt := 1
			for i++; i < len(p.tokens) && (p.tokens[i].Type != lexer.FLOW_ARRAY_END || cnt != 1); i++ {
				t = append(t, p.tokens[i])
				if p.tokens[i].Type == lexer.FLOW_ARRAY_END {
					cnt--
				}
				if p.tokens[i].Type == lexer.FLOW_ARRAY_START {
					cnt++
				}
			}
			chAst := newParser(t).parse()
			p.nodes.AddChild(chAst)
		case lexer.ANCHOR:
			p.nodes.Add(ast.NewNodeValue(ast.ANCHOR, el.Value))
			level := el.Level
			tarr := []*lexer.Token{}
			for i+1 < len(p.tokens) && p.tokens[i+1].Level >= level {
				i++
				tarr = append(tarr, p.tokens[i])
			}
			chAst := newParser(tarr).parse()
			p.nodes.AddChild(chAst)
		case lexer.ALIAS:
			p.nodes.Add(ast.NewNodeValue(ast.ALIAS, el.Value))
		case lexer.SEQUENCE:
			p.nodes.Add(ast.NewNodeValue(ast.ARRAY_ELEMENT, ""))
			tarr := []*lexer.Token{}
			level := el.Level
			for i+1 < len(p.tokens) && p.tokens[i+1].Level > level {
				i++
				tarr = append(tarr, p.tokens[i])
			}
			chAst := newParser(tarr).parse()
			p.nodes.AddChild(chAst)
		case lexer.TAG:
			p.nodes.Add(ast.NewNodeValue(ast.TAG, el.Value))
			level := el.Level
			tarr := []*lexer.Token{}
			for i+1 < len(p.tokens) && p.tokens[i+1].Level >= level {
				i++
				tarr = append(tarr, p.tokens[i])
			}
			chAst := newParser(tarr).parse()
			p.nodes.AddChild(chAst)
		case lexer.DOCUMENT_END:
			p.nodes.Add(ast.NewNodeValue(ast.FILE_END, ""))
		case lexer.DOCUMENT_START:
			p.nodes.Add(ast.NewNodeValue(ast.DOCUMENT_START, ""))
		case lexer.FLOW_ARRAY_END, lexer.FLOW_MAP_END, lexer.NEW_LINE:
			continue
		default:
			panic(fmt.Sprintf("unsupported token type: %v at %d", el.Type, i))
		}
	}
	return p.nodes
}

func _applyTagFunction(value ast.NodeValue, tag string) ast.NodeValue {
	switch tag {
	case "!str":
		return StringTagFunction(value)
	case "!binary":
		return BinaryTagFunction(value)
	default:
		value.Value = ""
		value.Type = ast.STRING
		return value
	}
}

func applyTag(tag string, nd *linked_list.LinkedListNode[ast.NodeValue]) {
	if nd == nil {
		return
	}
	if nd.Value.Type == ast.SCALAR {
		nd.Value = _applyTagFunction(nd.Value, tag)
	}
	applyTag(tag, nd.Child)
	applyTag(tag, nd.Next)
}

func findTags(nd *linked_list.LinkedListNode[ast.NodeValue]) {
	if nd == nil {
		return
	}
	if nd.Value.Type == ast.TAG {
		applyTag(nd.Value.Value, nd.Child)
		return
	}
	findTags(nd.Next)
	findTags(nd.Child)
}

func (p *Parser) applyTags() {
	for el := range p.nodes.Iterate() {
		if el.Value.Type == ast.TAG {
			applyTag(el.Value.Value, el.Child)
		}
		findTags(el.Child)
	}
}

func _changeTypes(el *linked_list.LinkedListNode[ast.NodeValue]) {
	if el == nil {
		return
	}
	if el.Value.Type == ast.SCALAR {
		el.Value = ProcessScalarType(el.Value)
	}
	_changeTypes(el.Child)
	_changeTypes(el.Next)
}

func (p *Parser) changeTypes() {
	_changeTypes(p.nodes.GetHead())
}

func (p *Parser) _findAnchors(el *linked_list.LinkedListNode[ast.NodeValue]) {
	if el == nil {
		return
	}
	if el.Value.Type == ast.ANCHOR {
		p.anchors[el.Value.Value] = el.Child
	}
	p._findAnchors(el.Child)
	p._findAnchors(el.Next)
}

func (p *Parser) findAnchors() {
	p._findAnchors(p.nodes.GetHead())
}

func (p *Parser) _applyAnchors(el *linked_list.LinkedListNode[ast.NodeValue]) {
	if el == nil {
		return
	}
	if el.Value.Type == ast.ALIAS {
		if val, ok := p.anchors[el.Value.Value]; ok {
			el.Child = val
		}
	}
	p._applyAnchors(el.Child)
	p._applyAnchors(el.Next)
}

func (p *Parser) applyAnchors() {
	p._applyAnchors(p.nodes.GetHead())
}

func Parse(src []*lexer.Token) *ast.AST {
	p := newParser(src)
	p.parse()
	p.applyTags()
	p.changeTypes()
	p.findAnchors()
	p.applyAnchors()
	return p.nodes
}
