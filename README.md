# hashstructure

Since mitchellh has continued his maintenance of the project, I refer to his [repository](https://github.com/mitchellh/hashstructure).

## Features

  * Hash any arbitrary Go value, including complex types.

  * Tag a struct field to ignore it and not affect the hash value.

  * Tag a slice type struct field to treat it as a set where ordering
    doesn't affect the hash code but the field itself is still taken into
    account to create the hash value.

  * Optionally, specify a custom hash function to optimize for speed, collision
    avoidance for your data set, etc.
  
  * Optionally, hash the output of `.String()` on structs that implement fmt.Stringer,
    allowing effective hashing of time.Time

  * Optionally, override the hashing process with a `hash` field or by implementing `hashstructure.Hashable`.

## Installation

```bash
go get -u github.com/gofunky/hashstructure
```

## Example

<!-- add-file: ./hashstructure_examples_test.go -->
``` go markdown-add-files
package hashstructure

import (
	"fmt"
)

func ExampleHash() {
	type ComplexStruct struct {
		Name     string
		Age      uint
		Metadata map[string]interface{}
	}

	v := ComplexStruct{
		Name: "gofunky",
		Age:  64,
		Metadata: map[string]interface{}{
			"car":      true,
			"location": "California",
			"siblings": []string{"Bob", "John"},
		},
	}

	hash, err := Hash(v, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d", hash)
	// Output:
	// 12836943650294093551
}

```
