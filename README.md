# Sequence Generator GO library

### Install
```yaml
import:
- package: bitbucket.lzd.co/lgo/sequence_template
  version: x.x.x
```

### Examples

##### 1) Postgress genenerator

```go
package main

import (
	"net/http"

	"bitbucket.org/snapmartinc/sequence_template"
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
