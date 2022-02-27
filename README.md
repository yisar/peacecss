# peacecss
Postcss alternative.

### [Playground](https://yisar.github.io/peacecss/)

### Use

```go
package main

import (
	"fmt"
	"github.com/yisar/peacecss"
)

func main() {
	parser := peacecss.NewParser()

	s := []byte(".a{color:#fff;}")

	ast := parser.Parse(s)
	
	ast.Traverse(func (node *peacecss.CSSDefinition){
		fmt.Printf("%v", node)
		
		node.Selector.Selector = ".b"
	
		fmt.Printf("%v", node)
	})

	mini := ast.Minisize()

	fmt.Printf("%s", mini.String())
	
}
```
