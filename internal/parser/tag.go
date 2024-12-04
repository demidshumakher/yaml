package parser

import (
	"encoding/base64"
	"github.com/demidshumakher/yaml/internal/ast"
)

type TagFunction func(ast.NodeValue) ast.NodeValue

func BinaryTagFunction(value ast.NodeValue) ast.NodeValue {
	value.Type = ast.STRING
	data, _ := base64.StdEncoding.DecodeString(value.Value)
	value.Value = string(data)
	return value
}

func StringTagFunction(value ast.NodeValue) ast.NodeValue {
	value.Type = ast.STRING
	return value
}
