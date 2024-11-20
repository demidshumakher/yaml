package main

import (
	"fmt"
	"github.com/demidshumakher/yaml/internal/lexer"
)

func main() {
	str := `---
canonical: "kasjfd \" skaljf"
...
`

	test := lexer.Scan([]rune(str))
	fmt.Println(test)

}
