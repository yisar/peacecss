package main

import (
	"fmt"
	// "syscall/js"
	nextcss "github.com/yisar/nextcss"
)

func main() {
	parser := nextcss.NewParser()

	s := []byte(".a{color:#fff;}")

	ast := parser.Parse(s)
	
	ast.Walk(func (node *nextcss.CSSDefinition){
		fmt.Printf("before: %v\n", node)
		
		node.Selector.Selector = ".b"
	
		fmt.Printf("after: %v\n", node)
	})

	ast.Minisize()

	// json := ast.ToPrettyJSON()

	// fmt.Printf("ast: %s\n", json)

	// js.Global().Set("cssnextParse", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	parser := cssnext.NewParser()
	// 	s := []byte(args[0].String())
	// 	ast := parser.Parse(s)

	// 	out := ast.ToPrettyJSONString()
	// 	return out
	// }))

	// select {}
}
