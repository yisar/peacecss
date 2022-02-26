package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	cuss "github.com/yisar/cuss"
)

func JSONStringfy(data []*cuss.CSSDefinition) string {
	ret, _ := json.MarshalIndent(data, "", "  ")
	return string(ret)
}

func registerWasm() {
	js.Global().Set("cussParse", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		parser := cuss.NewParser()
		s := []byte(args[0].String())
		ast := parser.Parse(s)
		out := JSONStringfy(ast.GetData())
		return out
	}))

	js.Global().Set("cussMinisize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		parser := cuss.NewParser()
		s := []byte(args[0].String())
		ast := parser.Parse(s)
		out := ast.Minisize()
		return out.String()
	}))

	select {}
}

func test() {
	parser := cuss.NewParser()

	s := []byte(".a{color:#fff;}")

	ast := parser.Parse(s)

	ast.Walk(func (node *cuss.CSSDefinition){
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
