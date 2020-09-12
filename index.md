---
title: Overview
---
# hashstructure

[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/gofunky/hashstructure/build/master?style=for-the-badge)](https://github.com/gofunky/hashstructure/actions)
[![Codecov](https://img.shields.io/codecov/c/github/gofunky/hashstructure?style=for-the-badge)](https://codecov.io/gh/gofunky/hashstructure)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue?style=for-the-badge)](https://pkg.go.dev/github.com/gofunky/hashstructure?tab=doc)
[![Renovate Status](https://img.shields.io/badge/renovate-enabled-green?style=for-the-badge&logo=renovatebot&color=1a1f6c)](https://app.renovatebot.com/dashboard#github/gofunky/hashstructure)
[![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/gofunky/hashstructure?style=for-the-badge)](https://libraries.io/github/gofunky/hashstructure)
[![CodeFactor](https://www.codefactor.io/repository/github/gofunky/hashstructure/badge?style=for-the-badge)](https://www.codefactor.io/repository/github/gofunky/hashstructure)
[![Go Report Card](https://goreportcard.com/badge/github.com/gofunky/hashstructure?style=for-the-badge)](https://goreportcard.com/report/github.com/gofunky/hashstructure)
[![GitHub License](https://img.shields.io/github/license/gofunky/hashstructure.svg?style=for-the-badge)](https://github.com/gofunky/hashstructure/blob/master/LICENSE)
[![Fossa](https://img.shields.io/badge/OSS-compliant-green?style=for-the-badge&logo=fossa)](https://app.fossa.com/reports/b3739086-45a2-4fc3-987d-2871e321c849)
[![GitHub last commit](https://img.shields.io/github/last-commit/gofunky/hashstructure.svg?style=for-the-badge&color=9cf)](https://github.com/gofunky/hashstructure/commits/master)

a go library for creating a unique hash value for arbitrary values

This can be used to key values in a hash (for use in a map, set, etc.) that are complex.
The most common use case is comparing two values without sending data across the network, caching values locally (de-dup), etc.

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
``` go 
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
