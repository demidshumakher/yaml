package lexer

import "fmt"

type TokenType int

const (
	SCALAR TokenType = iota + 1
	MAP_DELIMITER
	FLOW_MAP_START
	FLOW_MAP_END
	FLOW_ARRAY_START
	FLOW_ARRAY_END
	DOCUMENT_START
	DOCUMENT_END
	ANCHOR
	ALIAS
	QUOTE
	SEQUENCE
	COMMPLEX_MAP_KEY
	MULTILINE_STRING
	TAG
	FLOW_DELIMITER
	NEW_LINE
)

var TYPESNAMES = map[TokenType]string{
	MAP_DELIMITER:    "MAP_DELIMITER",
	FLOW_MAP_START:   "FLOW_MAP_START",
	FLOW_MAP_END:     "FLOW_MAP_END",
	FLOW_ARRAY_START: "FLOW_ARRAY_START",
	FLOW_ARRAY_END:   "FLOW_ARRAY_END",
	DOCUMENT_START:   "DOCUMENT_START",
	DOCUMENT_END:     "DOCUMENT_END",
	ANCHOR:           "ANCHOR",
	ALIAS:            "ALIAS",
	QUOTE:            "QUOTE",
	SEQUENCE:         "SEQUENCE",
	COMMPLEX_MAP_KEY: "COMMPLEX_MAP_KEY",
	MULTILINE_STRING: "MULTILINE_STRING",
	TAG:              "TAG",
	FLOW_DELIMITER:   "FLOW_DELIMITER",
	SCALAR:           "SCALAR",
	NEW_LINE:         "NEW_LINE",
}

func (t TokenType) String() string {
	return TYPESNAMES[t]
}

type Token struct {
	Value string
	Type  TokenType
	Level int
}

func (t Token) String() string {
	return fmt.Sprintf(" { %s\t%s\t%d } ", t.Value, t.Type, t.Level)
}

func flowMapStartToken(level int) *Token {
	return &Token{
		Type:  FLOW_MAP_START,
		Level: level,
	}
}

func flowMapEndToken(level int) *Token {
	return &Token{
		Type:  FLOW_MAP_END,
		Level: level,
	}
}

func mapDelimiterToken(level int) *Token {
	return &Token{
		Type:  MAP_DELIMITER,
		Level: level,
	}
}

func flowArrayStartToken(level int) *Token {
	return &Token{
		Type:  FLOW_ARRAY_START,
		Level: level,
	}
}

func flowArrayEndToken(level int) *Token {
	return &Token{
		Type:  FLOW_ARRAY_END,
		Level: level,
	}
}

func documentEndToken(level int) *Token {
	return &Token{
		Type:  DOCUMENT_END,
		Level: level,
	}
}

func documentStartToken(level int) *Token {
	return &Token{
		Type:  DOCUMENT_START,
		Level: level,
	}
}

func anchorToken(level int, name string) *Token {
	return &Token{
		Type:  ANCHOR,
		Value: name,
		Level: level,
	}
}

func aliasToken(level int, name string) *Token {
	return &Token{
		Type:  ALIAS,
		Value: name,
		Level: level,
	}
}

func quoteToken(level int, value string) *Token {
	return &Token{
		Type:  QUOTE,
		Value: value,
		Level: level,
	}
}

func sequenceToken(level int) *Token {
	return &Token{
		Type:  SEQUENCE,
		Level: level,
	}
}

func complexMapKeyToken(level int) *Token {
	return &Token{
		Type:  COMMPLEX_MAP_KEY,
		Level: level,
	}
}

func multiLineStringToken(level int, value string) *Token {
	return &Token{
		Type:  MULTILINE_STRING,
		Value: value,
		Level: level,
	}
}

func tagToken(level int, value string) *Token {
	return &Token{
		Type:  TAG,
		Value: value,
		Level: level,
	}
}

func flowDelimiterToken(level int) *Token {
	return &Token{
		Type:  FLOW_DELIMITER,
		Level: level,
	}
}

func scalarToken(level int, value string) *Token {
	return &Token{
		Type:  SCALAR,
		Value: value,
		Level: level,
	}
}

func newLineToken(level int) *Token {
	return &Token{
		Type:  NEW_LINE,
		Level: level,
	}
}
