package json

import (
	"fmt"
	"github.com/demidshumakher/yaml/internal/ast"
	"github.com/demidshumakher/yaml/pkg/linked_list"
	"io"
)

type JsonBackend struct {
	out   io.Writer
	nodes *ast.AST
}

func NewJsonBackend(out io.Writer, nd *ast.AST) *JsonBackend {
	return &JsonBackend{
		out:   out,
		nodes: nd,
	}
}

func (bc *JsonBackend) writeString(el *linked_list.LinkedListNode[ast.NodeValue]) {
	str := fmt.Sprintf("\"%s\"", el.Value.Value)
	bc.out.Write([]byte(str))
}

func (bc *JsonBackend) writeInteger(el *linked_list.LinkedListNode[ast.NodeValue]) {
	bc.out.Write([]byte(el.Value.Value))
}

func (bc *JsonBackend) writeFloat(el *linked_list.LinkedListNode[ast.NodeValue]) {
	bc.out.Write([]byte(el.Value.Value))
}

func (bc *JsonBackend) writeFloatNegatifeInf(el *linked_list.LinkedListNode[ast.NodeValue]) {
	bc.out.Write([]byte("-Infinity"))
}

func (bc *JsonBackend) writeFloatPossitiveInf(el *linked_list.LinkedListNode[ast.NodeValue]) {
	bc.out.Write([]byte("Infinity"))
}

func (bc *JsonBackend) writeFloatNan(el *linked_list.LinkedListNode[ast.NodeValue]) {
	bc.out.Write([]byte("NaN"))
}

func (bc *JsonBackend) writeNull(el *linked_list.LinkedListNode[ast.NodeValue]) {
	bc.out.Write([]byte("null"))
}

func (bc *JsonBackend) writeTrue(el *linked_list.LinkedListNode[ast.NodeValue]) {
	bc.out.Write([]byte("true"))
}

func (bc *JsonBackend) writeFalse(el *linked_list.LinkedListNode[ast.NodeValue]) {
	bc.out.Write([]byte("false"))
}

func (bc *JsonBackend) writeTimestamp(el *linked_list.LinkedListNode[ast.NodeValue]) {
	str := fmt.Sprintf("\"%s\"", el.Value.Value)
	bc.out.Write([]byte(str))
}

func (bc *JsonBackend) writeArrayElement(el *linked_list.LinkedListNode[ast.NodeValue]) *linked_list.LinkedListNode[ast.NodeValue] {
	bc.out.Write([]byte("[ "))
	fs := true
	for el != nil && el.Value.Type == ast.ARRAY_ELEMENT {
		if !fs {
			bc.out.Write([]byte(","))
		}
		bc._run(el.Child)
		el = el.Next
		fs = false
	}
	bc.out.Write([]byte("]"))
	return el
}

func (bc *JsonBackend) writeMapKey(el *linked_list.LinkedListNode[ast.NodeValue]) *linked_list.LinkedListNode[ast.NodeValue] {
	bc.out.Write([]byte("{"))
	fs := true
	for el != nil && el.Value.Type == ast.MAP_KEY {
		if !fs {
			bc.out.Write([]byte(","))
		}
		str := fmt.Sprintf("\"%s\":", el.Value.Value)
		bc.out.Write([]byte(str))
		if el.Next == nil || el.Next.Value.Type != ast.MAP_VALUE {
			bc.writeNull(el.Child)
		} else {
			el = el.Next
			bc._run(el.Child)
		}
		el = el.Next
		fs = false
	}
	bc.out.Write([]byte("}"))
	return el
}

func (bc *JsonBackend) writeArray(el *linked_list.LinkedListNode[ast.NodeValue]) {
	bc.out.Write([]byte("["))
	fs := true
	el = el.Child
	for el != nil {
		if !fs {
			bc.out.Write([]byte(","))
		}
		el = bc.write(el)
		fs = false
	}
	bc.out.Write([]byte("]"))
}

//func (bc *JsonBackend) writeMap(el *linked_list.LinkedListNode[ast.NodeValue]) {
//
//}

func (bc *JsonBackend) write(el *linked_list.LinkedListNode[ast.NodeValue]) *linked_list.LinkedListNode[ast.NodeValue] {
	if el == nil {
		return nil
	}
	//fmt.Println(el.Value)
	switch el.Value.Type {
	case ast.ARRAY_ELEMENT:
		return bc.writeArrayElement(el)
	case ast.MAP_KEY:
		return bc.writeMapKey(el)
	case ast.ARRAY:
		bc.writeArray(el)
	case ast.MAP:
		bc.writeMapKey(el.Child)
	case ast.STRING:
		bc.writeString(el)
	case ast.INTEGER:
		bc.writeInteger(el)
	case ast.FLOAT:
		bc.writeFloat(el)
	case ast.FLOAT_POSITIVE_INF:
		bc.writeFloatPossitiveInf(el)
	case ast.FLOAT_NEGATIVE_INF:
		bc.writeFloatNegatifeInf(el)
	case ast.FLOAT_NAN:
		bc.writeFloatNan(el)
	case ast.NULL:
		bc.writeNull(el)
	case ast.TRUE:
		bc.writeTrue(el)
	case ast.FALSE:
		bc.writeFalse(el)
	case ast.TIMESTAMP:
		bc.writeTimestamp(el)
	default:
		bc._run(el.Child)
	}
	return el.Next
}

func (bc *JsonBackend) _run(el *linked_list.LinkedListNode[ast.NodeValue]) {
	if el == nil {
		return
	}
	nx := bc.write(el)
	bc._run(nx)
}

func (bc *JsonBackend) run() {
	bc._run(bc.nodes.GetHead())
}

func (bc *JsonBackend) Run() {
	bc.run()
}
