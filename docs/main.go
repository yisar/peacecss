package main

import (
	"syscall/js"
	cssnext "github.com/yisar/cssnext"
)

func main() {
	// parser := parser.NewParser()

	// s := []byte(".a{color:#fff;}")

	// ast := parser.Parse(s)

	// out := ast.ToPrettyJSON()

	// fmt.Printf("ast: %s\n", out)

	js.Global().Set("cssnextParse", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		parser := cssnext.NewParser()
		s := []byte(args[0].String())
		ast := parser.Parse(s)

		out := ast.ToPrettyJSONString()
		return out
	}))

	select {}
}
