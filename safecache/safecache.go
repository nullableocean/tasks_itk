package safecache

import "sync"

/*
## Реализация потокобезопасного кеша

### Описание задачи
Ваша задача — реализовать потокобезопасный кеш для хранения данных в формате ключ-значение. Кеш должен безопасно обрабатывать одновременные операции записи и чтения из множества горутин.

### Требования
1. Реализовать структуру `SafeCache` с методами:
   - `Set(key string, value string)` — добавляет значение в кеш.
   - `Get(key string) (string, bool)` — возвращает значение по ключу.
2. Гарантировать отсутствие data race при параллельном доступе.
*/

type SafeCache struct {
	store map[string]string
	mu    sync.RWMutex
}

func NewSafeCache() *SafeCache {
	return &SafeCache{
		store: make(map[string]string),
		mu:    sync.RWMutex{},
	}
}

func (sc *SafeCache) Set(key, value string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.store[key] = value
}

func (sc *SafeCache) Get(key string) (string, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	v, ex := sc.store[key]
	return v, ex
}
