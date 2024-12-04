package main

import (
	"fmt"
	"github.com/demidshumakher/yaml/internal/backend/json"
	"github.com/demidshumakher/yaml/internal/lexer"
	"github.com/demidshumakher/yaml/internal/parser"
	"io"
	"os"
)

func myConverter(input []rune, output string) error {
	ast := parser.Parse(lexer.Scan(input))
	fl, err := os.Create(output)
	if err != nil {
		return err
	}
	defer fl.Close()
	json.NewJsonBackend(fl, ast).Run()
	//toml.NewTomlBackend(fl, ast).Run()
	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("give a input file")
		return
	}
	input := os.Args[1]
	var output string
	if len(os.Args) < 3 {
		output = "out.json"
	} else {
		output = os.Args[2]
	}

	fl, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	bf, err := io.ReadAll(fl)
	if err != nil {
		panic(err)
	}
	myConverter([]rune(string(bf)), output)
}
