package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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
		ast.Traverse(rex2rem)
		out := JSONStringfy(ast.GetData())
		return out
	}))

	js.Global().Set("peacecssMinisize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		parser := peacecss.NewParser()
		s := []byte(args[0].String())
		ast := parser.Parse(s)
		ast.Traverse(rex2rem)
		out := ast.Minisize()
		return out.String()
	}))

	select {}
}

func rex2rem(node *peacecss.CSSDefinition) {
	for _, r := range node.Rules {
		reg, _ := regexp.Compile("([0-9]+)rpx")
		r.Value.Value = reg.ReplaceAllStringFunc(r.Value.Value, func(s string) string {
			num, _ := strconv.ParseFloat(s[:len(s) - 3], 64)
			rem := num / 75
			return fmt.Sprintf("%.2f", rem) + "rem"
		})
	}
}

func test() {
	parser := peacecss.NewParser()

	s := []byte(".a{width:75rpx}")

	ast := parser.Parse(s)

	ast.Traverse(rex2rem)

	ast.Minisize()

	json := JSONStringfy(ast.GetData())

	fmt.Printf("ast: %s\n", json)
}

func main() {
	registerWasm()
	// test()
}
