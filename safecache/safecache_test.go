package safecache

import (
	"sync"
	"testing"
)

// go test -race .
func Test_ConcurrentAccess(t *testing.T) {
	sc := &SafeCache{store: make(map[string]string)}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := "key"
			if idx%2 == 0 {
				key = "qwe"
			}
			sc.Set(key, "value")
		}(i)

		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := "key"
			if idx%2 == 0 {
				key = "qwe"
			}
			sc.Get(key)
		}(i)
	}

	wg.Wait()
}
