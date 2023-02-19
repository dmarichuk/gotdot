# gotdot
Simple module to export env variables from .env-like files

# Installation

`go get github.com/dmarichuk/gotdot`

# Examples

```
package main

import (
	"fmt"
    "os"

    "github.com/dmarichuk/gotdot"
)

func main() {
	
	c := gotdot.Config // you can assign config to a variable
	
	// By default, gotdot looking for _.env_ file in current directory.
	// To change the source file, change _Path_ parameter in Config struct
	c.Path = "./.env.dev"

	c.Load() // Parses and loads file to ENV variables and inner config struct

	// Now you can use os.Expandenv to include your env variables to
    example := os.Expandenv("$VAR1 - $VAR2")

    // If need to cast your env variable to another type

    envValue, _ := c.Get("VAR1") // Get variable from inner struct, error if not exist
    envValue.Cast("int") // as int64, you can also cast to "bool"(as _bool_) or "float"(as _float64_) or
    castedValue := envValue.Export() // assign casted value to a variable
    // one can also use pipe-like syntax for this
    newVar := envValue.Cast("float").Import()
}
```
