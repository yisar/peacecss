# cuss
Postcss alternative.

### [Playground](https://yisar.github.io/cuss/)

### Use

```go
package main

import (
	"fmt"
	"github.com/yisar/cuss"
)

func main() {
	parser := nextcss.NewParser()

	s := []byte(".a{color:#fff;}")

	ast := parser.Parse(s)
	
	ast.Walk(func (node *nextcss.CSSDefinition){
		fmt.Printf("%v", node)
		
		node.Selector.Selector = ".b"
	
		fmt.Printf("%v", node)
	})

	mini := ast.Minisize()

	fmt.Printf("%s", mini.String())
	
}
```
