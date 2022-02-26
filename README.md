# nextcss
Postcss alternative.

### [Playground](https://yisar.github.io/nextcss/)

### Use

```go
package main

import (
	"fmt"
	"github.com/yisar/nextcss"
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

	mini := ast.Minisize()

	fmt.Printf("ast: %s\n", mini.String())
	
}
```