# hashstructure

[![GoDoc](https://godoc.org/github.com/gofunky/hashstructure?status.svg)](https://godoc.org/github.com/gofunky/hashstructure)

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
