package main

import (
	"fmt"
	"github.com/demidshumakher/yaml/internal/backend/json"
	"github.com/demidshumakher/yaml/internal/lexer"
	"github.com/demidshumakher/yaml/internal/parser"
	"io"
	"os"
)

type temp struct {
	wr io.Writer
}

func (t temp) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return t.wr.Write(p)
}

func main() {
	fl, err := os.Open("tests/anchors.yaml")
	if err != nil {
		panic(err)
	}
	bf, err := io.ReadAll(fl)
	if err != nil {
		panic(err)
	}

	tokens := lexer.Scan([]rune(string(bf)))
	ast := parser.Parse(tokens)
	file, err := os.Create("out.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bck := json.NewJsonBackend(file, ast)
	fmt.Println("OK")
	bck.Run()
	fmt.Println(ast)
	//fmt.Println(tokens)
	fmt.Println("END")

}
