package parser

type NodeType int

const (
	STRING NodeType = iota + 1
	INTEGER
	FLOAT
	FLOAT_POSITIVE_INF
	FLOAT_NEGATIVE_INF
	FLOAT_NAN
	ARRAY_START
	ARRAY_END
	//ARRAY_DE
	MAP_KEY
	MAP_VALUE
	DOCUMENT_START
	FILE_END
	TIMESTAMP
	_TAGTYPE
	_BLOCK_OF_SCALARS
)

type NodeValue struct {
	Type  NodeType
	Value string
}
