package cachettl

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*

# Универсальный потокобезопасный кэш с TTL, очисткой и JSON-сериализацией

Реализация in-memory кэша на Go с расширенными возможностями: автоматическое удаление ключей, очистка и сериализация данных.

---

## **Основные возможности**

**TTL (Time-To-Live)**
- Автоматическое удаление ключей по истечении времени жизни.

**Очистка кэша**
- Мгновенное удаление всех данных одной командой.

**Сериализация в JSON**
- Преобразование актуальных данных в JSON-формат.

**Потокобезопасность**
- Использование `sync.RWMutex` для конкурентного доступа.

**Универсальное хранение**
- Поддержка любых типов данных через `interface{}`.
---

## **Методы**
### **Базовые операции**

 - `Set(key string, value interface{}, ttl time.Duration)` Добавляет значение с указанным TTL
 - `Get(key string) (interface{}, bool)`  Возвращает значение (с проверкой TTL)
 - `Delete(key string)`  Удаляет конкретный ключ
 - `Exists(key string) bool` Проверяет наличие непросроченного ключа

### **Расширенные функции**
 - `Clear()` Полностью очищает кэш
 - `ToJSON() ([]byte, error)` Сериализует данные в JSON
 - `GetAs[T any](key string) (T, error)`  Типизированное получение
*/

var (
	ErrStopped       = errors.New("cache stopped")
	ErrNotStopped    = errors.New("cant restart: not stopped")
	ErrNotFound      = errors.New("not found")
	ErrIncorrectType = errors.New("value type dont match")
)

var (
	ttlCheckInterval = time.Minute * 3
)

type cacheValue struct {
	val interface{}
	exp time.Time
}

type Cache struct {
	cache map[string]*cacheValue
	mu    sync.RWMutex

	processMu sync.Mutex
	stopped   int32
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

func NewCache() *Cache {
	c := &Cache{
		mu:     sync.RWMutex{},
		stopCh: make(chan struct{}),

		processMu: sync.Mutex{},
		wg:        sync.WaitGroup{},
	}

	c.initNewCacheStore()

	c.wg.Add(1)
	c.runObserveTtl()

	return c
}

func GetAs[T any](cache *Cache, key string) (T, error) {
	var out T

	v, ex := cache.Get(key)
	if !ex {
		return out, ErrNotFound
	}

	val, ok := v.(T)
	if !ok {
		return out, fmt.Errorf("%w: value type: %T, got: %T", ErrIncorrectType, v, out)
	}

	return val, nil
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cVal, ex := c.cache[key]
	if !ex {
		return nil, false
	}

	if c.isExpiredValue(cVal) {
		defer func() {
			go c.delete(key)
		}()

		return nil, false
	}

	return cVal.val, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	c.processMu.Lock()
	defer c.processMu.Unlock()

	if c.isStopped() {
		return ErrStopped
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &cacheValue{
		val: value,
		exp: time.Now().Add(ttl),
	}

	return nil
}

func (c *Cache) Delete(key string) {
	if !c.Exists(key) {
		return
	}

	c.delete(key)
}

func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cVal, ex := c.cache[key]
	if !ex {
		return false
	}

	if c.isExpiredValue(cVal) {
		defer func() {
			go c.delete(key)
		}()

		return false
	}

	return true
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.initNewCacheStore()
}

func (c *Cache) ToJSON() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	jsMap := make(map[string]interface{}, len(c.cache))
	for k, cv := range c.cache {
		jsMap[k] = cv.val
	}

	return json.Marshal(jsMap)
}

func (c *Cache) Restart() error {
	c.processMu.Lock()
	defer c.processMu.Unlock()

	if !c.isStopped() {
		return ErrNotStopped
	}

	c.setStopped(false)
	c.stopCh = make(chan struct{})

	c.wg.Add(1)
	c.runObserveTtl()

	return nil
}

func (c *Cache) Stop() error {
	c.processMu.Lock()
	defer c.processMu.Unlock()

	if c.isStopped() {
		return ErrStopped
	}

	c.setStopped(true)
	close(c.stopCh)
	c.wg.Wait()

	return nil
}

func (c *Cache) runObserveTtl() {
	go func() {
		defer c.wg.Done()

		ticker := time.NewTicker(ttlCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-c.stopCh:
				return
			case <-ticker.C:
				c.clearExpiredValues()
			}
		}
	}()
}

func (c *Cache) clearExpiredValues() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, cv := range c.cache {
		select {
		case <-c.stopCh:
			return
		default:
			if c.isExpiredValue(cv) {
				delete(c.cache, k)
			}
		}
	}
}

func (c *Cache) isExpiredValue(v *cacheValue) bool {
	return v.exp.Before(time.Now())
}

func (c *Cache) setStopped(status bool) {
	var stopVal int32
	if status {
		stopVal = 1
	}

	atomic.StoreInt32(&c.stopped, stopVal)
}

func (c *Cache) isStopped() bool {
	return atomic.LoadInt32(&c.stopped) == 1
}

func (c *Cache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

func (c *Cache) initNewCacheStore() {
	c.cache = make(map[string]*cacheValue)
}
