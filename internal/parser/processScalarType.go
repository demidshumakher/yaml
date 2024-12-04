package parser

import (
	"github.com/demidshumakher/yaml/internal/ast"
	"slices"
	"strings"
	"unicode"
)

type ProcessScalarTypeFunction func(ast.NodeValue) (ast.NodeValue, bool)

func ProcessScalarType(nd ast.NodeValue) ast.NodeValue {
	if nd, ok := toNull(nd); ok {
		return nd
	}
	if nd, ok := toFloatNan(nd); ok {
		return nd
	}
	if nd, ok := toFloatNegativeInf(nd); ok {
		return nd
	}
	if nd, ok := toFloatPossitiveInf(nd); ok {
		return nd
	}
	if nd, ok := toBoolean(nd); ok {
		return nd
	}
	if nd, ok := toTimeStamp(nd); ok {
		return nd
	}
	if nd, ok := toInteger(nd); ok {
		return nd
	}
	if nd, ok := toFloat(nd); ok {
		return nd
	}

	nd, _ = toString(nd)
	return nd
}

// only xxxx-xx-xx
func toTimeStamp(nd ast.NodeValue) (ast.NodeValue, bool) {
	v := strings.TrimSpace(nd.Value)
	if len(v) == 10 {
		for idx, el := range v {
			switch idx {
			case 0, 1, 2, 3, 5, 6, 8, 9:
				if !unicode.IsDigit(el) {
					return nd, false
				}
			default:
				if el != '-' {
					return nd, false
				}
			}
		}
	} else {
		return nd, false
	}
	nd.Value = v
	nd.Type = ast.TIMESTAMP
	return nd, true
}

func toInteger(nd ast.NodeValue) (ast.NodeValue, bool) {
	v := strings.TrimSpace(nd.Value)
	res := ""
	for idx, r := range v {
		if idx == 0 && (r == '+' || r == '-') {
			if r == '-' {
				res += string(r)
			}
			continue
		}
		if unicode.IsDigit(r) {
			res += string(r)
		} else {
			return nd, false
		}
	}
	nd.Value = res
	nd.Type = ast.INTEGER
	return nd, true
}

func toFloat(nd ast.NodeValue) (ast.NodeValue, bool) {
	v := strings.TrimSpace(nd.Value)
	idx := strings.Index(v, ".")
	if idx == -1 {
		return nd, false
	}
	cpy := nd
	cpy.Value = v[:idx]
	if _, ok := toInteger(cpy); !ok {
		return nd, false
	}
	intpart := v[:idx+1]
	for i, el := range v[idx+1:] {
		if unicode.IsDigit(el) {
			intpart += string(el)
		} else if el == 'e' || el == 'E' {
			intpart += string(el)
			index := idx + 1 + i
			if v[index+1] == '+' || v[index+1] == '-' {
				intpart += string(v[index+1])
				index++
			}
			for _, el := range v[index+1:] {
				if !unicode.IsDigit(el) {
					return nd, false
				}
				intpart += string(el)
			}
			nd.Type = ast.FLOAT
			nd.Value = intpart
			return nd, true
		} else {
			return nd, false
		}
	}
	nd.Value = intpart
	nd.Type = ast.FLOAT
	return nd, true
}

func toString(nd ast.NodeValue) (ast.NodeValue, bool) {
	nd.Type = ast.STRING
	return nd, true
}

func toBoolean(nd ast.NodeValue) (ast.NodeValue, bool) {
	v := strings.TrimSpace(nd.Value)
	fs := []string{"n", "N", "no", "No", "NO", "false", "False", "FALSE", "off", "Off", "OFF"}
	tr := []string{"y", "Y", "yes", "Yes", "YES", "true", "True", "TRUE", "on", "On", "ON"}
	if slices.Contains(fs, v) {
		nd.Type = ast.FALSE
		nd.Value = ""
		return nd, true
	}
	if slices.Contains(tr, v) {
		nd.Type = ast.TRUE
		nd.Value = ""
		return nd, true
	}

	return nd, false
}

func toFloatPossitiveInf(nd ast.NodeValue) (ast.NodeValue, bool) {
	v := strings.TrimSpace(nd.Value)
	if v == ".inf" {
		nd.Type = ast.FLOAT_POSITIVE_INF
		nd.Value = ""
		return nd, true
	}
	return nd, false
}

func toFloatNegativeInf(nd ast.NodeValue) (ast.NodeValue, bool) {
	v := strings.TrimSpace(nd.Value)
	if v == "-.inf" {
		nd.Type = ast.FLOAT_NEGATIVE_INF
		nd.Value = ""
		return nd, true
	}
	return nd, false
}

func toFloatNan(nd ast.NodeValue) (ast.NodeValue, bool) {
	v := strings.TrimSpace(nd.Value)
	if v == ".nan" {
		nd.Type = ast.FLOAT_NAN
		nd.Value = ""
		return nd, true
	}
	return nd, false
}

func toNull(nd ast.NodeValue) (ast.NodeValue, bool) {
	v := strings.TrimSpace(nd.Value)
	if v == "null" {
		nd.Type = ast.NULL
		nd.Value = ""
		return nd, true
	}
	return nd, false
}
