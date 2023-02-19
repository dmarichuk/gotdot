# gotdot
Simple module to export env variables from .env-like files

# Installation

`go get github.com/dmarichuk/gotdot`

# Examples

```
package main

import (
    "os"

    "github.com/dmarichuk/gotdot"
)

func main() {
	
	c := gotdot.Config // you can assign config to a variable
	// or one can use gotdot.Config variable as is
	
	// By default, gotdot looking for .env file in current directory.
	// To change the source file, change Path parameter in Config struct
	c.Path = "./.env.dev"

	c.Load() // Parses and loads file to ENV variables and inner config struct

	// Now you can use os.Expandenv to include your env variables to
	example := os.Expandenv("$VAR1 - $VAR2")
	
	// If neede, one can cast env variable to another type
	
	envValue, _ := c.Get("VAR1") // Get variable from inner struct, error if not exist
	envValue.Cast("int") // as int64
	// you can also cast to "bool"(as bool) or to "float"(as float64)
	
	castedValue := envValue.Export() // assign casted value to a variable
	// one can also use pipe-like syntax for this
	newVar := envValue.Cast("float").Import()
}
```
