# Sequence Generator GO library

### Examples

##### 1) Postgress genenerator

```go
package main

import (
	"net/http"

	"github.com/best-expendables/sequence_template"
)

func main() {
	gen, err := sequencetemplate.NewPostgesGenerator()
	if err != nil {
		panic(err)
	}
    
    sequenceName := "order_number_code_sequence"
    prefix := "LG_ORDER_"
	seq, err := gen.Generate(sequenceName, prefix, 5)
	if err != nil {
		panic(err)
	}
	//Sample sequence: LG_ORDER_00001
}
```
