package waitmap

import (
	"slices"
	"sync"
	"testing"
	"time"
)

func TestWaitMap(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		m := New[string, string]()
		m.Set("foo", "bar")
		want := "bar"
		if got := m.Get("foo"); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Get and Set", func(t *testing.T) {
		m := New[string, string]()
		go func() {
			want := "bar"
			if got := m.Get("foo"); got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}()
		time.Sleep(100 * time.Millisecond)
		m.Set("foo", "bar")
	})

	t.Run("Set Set Get and Get", func(t *testing.T) {
		m := New[string, string]()
		m.Set("foo", "bar")
		m.Set("foo", "baz")
		want := "baz"
		if got := m.Get("foo"); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got := m.Get("foo"); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Get Get and Set", func(t *testing.T) {
		m := New[string, string]()
		want := "bar"
		go func() {
			if got := m.Get("foo"); got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}()
		go func() {
			if got := m.Get("foo"); got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}()
		time.Sleep(100 * time.Millisecond)
		m.Set("foo", "bar")
	})

	t.Run("Delete no key", func(t *testing.T) {
		m := New[string, string]()
		m.Delete("foo")
	})

	t.Run("Delete Set and Get", func(t *testing.T) {
		m := New[string, string]()
		m.Delete("foo")
		m.Set("foo", "bar")
		want := "bar"
		if got := m.Get("foo"); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Set Delete Get and Set", func(t *testing.T) {
		m := New[string, string]()
		m.Set("foo", "bar")
		m.Delete("foo")
		go func() {
			want := "baz"
			if got := m.Get("foo"); got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}()
		time.Sleep(100 * time.Millisecond)
		m.Set("foo", "baz")
	})

	t.Run("Multi Get and Set", func(t *testing.T) {
		keys := []string{"foo", "bar", "baz"}
		m := New[string, string]()
		sg := sync.WaitGroup{}
		for _, k := range keys {
			k := k
			sg.Add(1)
			go func() {
				want := "ok"
				if got := m.Get(k); got != want {
					t.Errorf("got %v, want %v", got, want)
				}
				sg.Done()
			}()
		}
		for _, k := range keys {
			m.Set(k, "ok")
		}
		sg.Wait()
	})
}

func TestTryGet(t *testing.T) {
	m := New[string, string]()
	v, ok := m.TryGet("foo")
	if ok {
		t.Errorf("got %v, want %v", v, nil)
	}
	m.Set("foo", "bar")
	{
		v, ok := m.TryGet("foo")
		if !ok {
			t.Errorf("got %v, want %v", v, nil)
		}
		want := "bar"
		if got := v; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestKeys(t *testing.T) {
	m := New[string, string]()
	m.Set("foo", "bar")
	m.Set("baz", "qux")
	want := []string{"foo", "baz"}
	got := m.Keys()
	if len(got) != len(want) {
		t.Errorf("got %v, want %v", got, want)
	}
	for _, v := range got {
		if !slices.Contains(want, v) {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestChan(t *testing.T) {
	m := New[string, string]()
	done := make(chan struct{})

	go func() {
		for {
			select {
			case got := <-m.Chan("foo"):
				want := "bar"
				if got != want {
					t.Errorf("got %v, want %v", got, want)
				}
				close(done)
				return
			case <-time.After(100 * time.Millisecond):
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	m.Set("foo", "bar")

	got := <-m.Chan("foo")
	want := "bar"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	<-done
}
