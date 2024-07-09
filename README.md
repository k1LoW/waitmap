# waitmap [![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/waitmap.svg)](https://pkg.go.dev/github.com/k1LoW/waitmap) [![CI](https://github.com/k1LoW/waitmap/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/waitmap/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/waitmap/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/waitmap/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/waitmap/time.svg)

`waitmap` is a concurrent and type safe map that allows you to wait for a key to be set.

## Usage

``` go
package main

import (
	"fmt"
	"time"

	"github.com/k1LoW/waitmap"
)

func main() {
	m := waitmap.New[string, int64]()

	ch := make(chan struct{})
	go func() {
		got := m.Get("foo")
		fmt.Println(got)
		close(ch)
	}()

	time.Sleep(100 * time.Millisecond)
	m.Set("foo", 123)
	<-ch

	// Output: 123
}
```
