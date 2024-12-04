package ast

import (
	"fmt"
	"github.com/demidshumakher/yaml/pkg/linked_list"
	"iter"
)

type NodeType int

const (
	STRING NodeType = iota + 1
	INTEGER
	FLOAT
	FLOAT_POSITIVE_INF
	FLOAT_NEGATIVE_INF
	FLOAT_NAN
	NULL
	ARRAY
	MAP
	ARRAY_ELEMENT
	MAP_KEY
	MAP_VALUE
	DOCUMENT_START
	TRUE
	FALSE
	FILE_END
	TIMESTAMP
	TAG
	SCALAR
	ANCHOR
	ALIAS
)

func (n NodeType) String() string {
	switch n {
	case STRING:
		return "STRING"
	case INTEGER:
		return "INTEGER"
	case FLOAT:
		return "FLOAT"
	case FLOAT_POSITIVE_INF:
		return "FLOAT_POSITIVE_INF"
	case FLOAT_NEGATIVE_INF:
		return "FLOAT_NEGATIVE_INF"
	case FLOAT_NAN:
		return "FLOAT_NAN"
	case NULL:
		return "NULL"
	case ARRAY:
		return "ARRAY"
	case MAP_KEY:
		return "MAP_KEY"
	case MAP_VALUE:
		return "MAP_VALUE"
	case DOCUMENT_START:
		return "DOCUMENT_START"
	case FILE_END:
		return "FILE_END"
	case TIMESTAMP:
		return "TIMESTAMP"
	case TAG:
		return "TAG"
	case SCALAR:
		return "SCALAR"
	case ARRAY_ELEMENT:
		return "ARRAY_ELEMENT"
	case MAP:
		return "MAP"
	case ANCHOR:
		return "ANCHOR"
	case ALIAS:
		return "ALIAS"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	}
	return "UNKNOWN"
}

type AST struct {
	ll *linked_list.LinkedList[NodeValue]
}

type NodeValue struct {
	Type  NodeType
	Value string
}

func (nd NodeValue) String() string {
	return fmt.Sprintf("{%s=%s}", nd.Type, nd.Value)
}

func NewAST() *AST {
	return &AST{
		ll: linked_list.New[NodeValue](),
	}
}

func NewNodeValue(tp NodeType, value string) NodeValue {
	return NodeValue{
		Type:  tp,
		Value: value,
	}
}

func (a *AST) Add(nd NodeValue) {
	a.ll.PushBack(nd)
}

func (a *AST) AddChild(nd *AST) {
	if nd == nil {
		return
	}
	a.ll.AddChild(nd.ll)
}

func (a *AST) Iterate() iter.Seq[*linked_list.LinkedListNode[NodeValue]] {
	return func(yield func(*linked_list.LinkedListNode[NodeValue]) bool) {
		for el := range a.ll.Iterate() {
			if !yield(el) {
				return
			}
		}
	}
}

func (a *AST) GetHead() *linked_list.LinkedListNode[NodeValue] {
	return a.ll.GetHead()
}

func (a *AST) String() string {
	return a.ll.String()
}
