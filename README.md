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
	parser := parser.NewParser()

	s := []byte(".a{color:#fff;}")

	ast := parser.Parse(s)

	fmt.Printf("ast: %s\n", ast.ToPrettyJSON())
	
}
```
