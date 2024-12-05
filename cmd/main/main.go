package main

import (
	"fmt"
	"github.com/demidshumakher/yaml/internal/backend/json"
	"github.com/demidshumakher/yaml/internal/backend/toml"
	"github.com/demidshumakher/yaml/internal/lexer"
	"github.com/demidshumakher/yaml/internal/parser"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"os"
)

func myConverter(input []rune, output string, Istoml bool) error {
	tokens := lexer.Scan(input)
	ast := parser.Parse(tokens)
	fl, err := os.Create(output)
	if err != nil {
		return err
	}
	defer fl.Close()
	fmt.Println(ast)
	if Istoml {
		toml.NewTomlBackend(fl, ast).Run()
	} else {
		json.NewJsonBackend(fl, ast).Run()
	}
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
	Istoml := false
	lib := false
	if len(os.Args) > 3 {
		Istoml = os.Args[3] == "toml"
		lib = os.Args[3] == "lib"
	}

	if lib {
		lib_Sol(input, output)
		return
	}
	fl, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	bf, err := io.ReadAll(fl)
	if err != nil {
		panic(err)
	}
	myConverter([]rune(string(bf)), output, Istoml)
}

func lib_Sol(input, output string) {
	file, _ := os.Open(input)
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	data, _ = yaml.YAMLToJSON(data)
	fl, _ := os.Create(output)
	defer fl.Close()
	fl.Write(data)
}
