package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	nextcss "github.com/yisar/nextcss"
)

func JSONStringfy(data []*nextcss.CSSDefinition) string {
	ret, _ := json.MarshalIndent(data, "", "  ")
	return string(ret)
}

func main() {
	// parser := nextcss.NewParser()

	// s := []byte(".a{color:#fff;}")

	// ast := parser.Parse(s)
	
	// ast.Walk(func (node *nextcss.CSSDefinition){
	// 	fmt.Printf("before: %v\n", node)
		
	// 	node.Selector.Selector = ".b"
	
	// 	fmt.Printf("after: %v\n", node)
	// })

	// ast.Minisize()

	// json := ast.ToPrettyJSON()

	// fmt.Printf("ast: %s\n", json)

	js.Global().Set("nextcssParse", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		parser := nextcss.NewParser()
		s := []byte(args[0].String())
		ast := parser.Parse(s)
		out := JSONStringfy(ast.GetData())
		return out
	}))

	js.Global().Set("nextcssMinisize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		parser := nextcss.NewParser()
		s := []byte(args[0].String())
		ast := parser.Parse(s)
		out := ast.Minisize()
		return out.String()
	}))

	select {}
}
