package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	peacecss "github.com/yisar/peacecss"
)

func JSONStringfy(data []*peacecss.CSSDefinition) string {
	ret, _ := json.MarshalIndent(data, "", "  ")
	return string(ret)
}

func registerWasm() {
	js.Global().Set("peacecssParse", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		parser := peacecss.NewParser()
		s := []byte(args[0].String())
		ast := parser.Parse(s)
		out := JSONStringfy(ast.GetData())
		return out
	}))

	js.Global().Set("peacecssMinisize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		parser := peacecss.NewParser()
		s := []byte(args[0].String())
		ast := parser.Parse(s)
		out := ast.Minisize()
		return out.String()
	}))

	select {}
}

func test() {
	parser := peacecss.NewParser()

	s := []byte(".a{color:#fff;}")

	ast := parser.Parse(s)

	ast.Walk(func (node *peacecss.CSSDefinition){
		fmt.Printf("before: %v\n", node)

		node.Selector.Selector = ".b"

		fmt.Printf("after: %v\n", node)
	})

	ast.Minisize()

	json := JSONStringfy(ast.GetData())

	fmt.Printf("ast: %s\n", json)
}

func main() {
	registerWasm()
}
