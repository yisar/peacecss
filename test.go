package nextcss

import (
	"fmt"
)

func main() {
	parser := NewParser()

	s := []byte(".a{color:#fff;}")

	ast := parser.Parse(s)

	out := ast.ToPrettyJSON()

	fmt.Printf("ast: %s\n", out)

}