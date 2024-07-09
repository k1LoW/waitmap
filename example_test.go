package waitmap_test

import (
	"fmt"
	"time"

	"github.com/k1LoW/waitmap"
)

func Example() {
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
