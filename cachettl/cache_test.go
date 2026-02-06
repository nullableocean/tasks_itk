package cachettl

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCache_SetGet(t *testing.T) {
	c := NewCache()
	defer c.Stop()

	key := "testk"
	value := "testv"
	ttl := time.Second * 10

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("error to set value: %v", err)
	}

	val, ex := c.Get(key)
	if !ex {
		t.Fatal("value should exist")
	}
	if val != value {
		t.Fatalf("expected %v, got %v", value, val)
	}
}

func TestCache_GetNonExistent(t *testing.T) {
	c := NewCache()
	defer c.Stop()

	val, ex := c.Get("notexist")
	if ex {
		t.Fatal("value should not exist")
	}

	if val != nil {
		t.Fatalf("expected nil, got %v", val)
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCache()
	defer c.Stop()

	key := "testk"
	value := "testv"
	ttl := time.Second * 10

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("error to set value: %v", err)
	}

	val, ex := c.Get(key)
	if !ex {
		t.Fatal("value should exist")
	}

	c.Delete(key)
	val, ex = c.Get(key)
	if ex {
		t.Fatal("value should not exist after deletion")
	}
	if val != nil {
		t.Fatalf("expected nil after deletion, got %v", val)
	}

	// ok with undefined key
	notExistKey := "notexkey"
	c.Delete(notExistKey)
}

func TestCache_Exists(t *testing.T) {
	c := NewCache()
	defer c.Stop()

	key := "testk"
	value := "testv"
	ttl := time.Second * 10

	if c.Exists(key) {
		t.Fatal("key should not exist after init")
	}

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("error to set value: %v", err)
	}

	if !c.Exists(key) {
		t.Fatal("key should exist after set")
	}

	c.Delete(key)
	if c.Exists(key) {
		t.Fatal("key should not exist after deletion")
	}
}

func TestCache_Clear(t *testing.T) {
	c := NewCache()
	defer c.Stop()

	keys := []string{"keyt_1", "keyt_2", "keyt_3"}
	values := []string{"value1", "value2", "value3"}
	ttl := time.Second * 10

	for i, key := range keys {
		err := c.Set(key, values[i], ttl)
		if err != nil {
			t.Fatalf("error to set value for key %s: %v", key, err)
		}
	}

	for _, key := range keys {
		if !c.Exists(key) {
			t.Fatalf("key %s should exist", key)
		}
	}

	c.Clear()

	for _, key := range keys {
		if c.Exists(key) {
			t.Fatalf("key %s should not exist after clear", key)
		}
	}
}

func TestCache_TTLExpiration(t *testing.T) {
	c := NewCache()
	defer c.Stop()

	key := "testk"
	value := "testv"
	ttl := time.Second * 1

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("error to set value: %v", err)
	}

	if !c.Exists(key) {
		t.Fatal("value should exist after set")
	}

	time.Sleep(ttl + time.Millisecond*50)

	if c.Exists(key) {
		t.Fatal("value should not exist after TTL expired")
	}

	val, exists := c.Get(key)
	if exists {
		t.Fatal("value should not exist after TTL expired")
	}
	if val != nil {
		t.Fatalf("expected nil after TTL expired, got %v", val)
	}
}

type testUser struct {
	Name string `json:"name"`
}

func TestCache_ToJSON(t *testing.T) {
	c := NewCache()
	defer c.Stop()

	tUser := testUser{Name: "OKNAME"}

	testData := map[string]interface{}{
		"name":    "nameval",
		"count":   42,
		"okey":    true,
		"balance": 32.14,

		"user": tUser,
	}

	for key, value := range testData {
		err := c.Set(key, value, time.Second*10)
		if err != nil {
			t.Fatalf("error to set value for key %s: %v", key, err)
		}
	}

	jsonData, err := c.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert cache to JSON: %v", err)
	}

	var parsedData map[string]interface{}
	err = json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		t.Fatalf("error to parse JSON: %v", err)
	}

	for key, expectedValue := range testData {
		actualValue, exists := parsedData[key]
		if !exists {
			t.Errorf("key %s should exist in JSON", key)
		}

		if _, ok := expectedValue.(testUser); ok {
			structMap, ok := actualValue.(map[string]interface{})
			if !ok {
				t.Errorf("expected map[string]interface{} for struct from json, got %v", actualValue)
				continue
			}

			if structMap["name"] != tUser.Name {
				t.Errorf("expected json name data = struct name field, got %v", structMap["name"])
			}

			continue
		}

		if expectInt, ok := expectedValue.(int); ok {
			expectedValue = float64(expectInt)
		}

		if expectedValue != actualValue {
			t.Errorf("expected %v for key %s, got %v", expectedValue, key, actualValue)
		}
	}
}

func TestCache_GetAs(t *testing.T) {
	c := NewCache()
	defer c.Stop()

	stringKey := "stringKey"
	stringValue := "stringVal"
	err := c.Set(stringKey, stringValue, time.Second*10)
	if err != nil {
		t.Fatalf("error to set value: %v", err)
	}

	result, err := GetAs[string](c, stringKey)
	if err != nil {
		t.Fatalf("error to get string value: %v", err)
	}
	if result != stringValue {
		t.Fatalf("expected %s, got %s", stringValue, result)
	}

	_, err = GetAs[int](c, stringKey)
	if !errors.Is(err, ErrIncorrectType) {
		t.Errorf("expected ErrIncorrectType, got %v", err)
	}

	intKey := "intKey"
	intValue := 42
	err = c.Set(intKey, intValue, time.Second*5)
	if err != nil {
		t.Fatalf("error to set value: %v", err)
	}

	resultInt, err := GetAs[int](c, intKey)
	if err != nil {
		t.Fatalf("error to get int value: %v", err)
	}
	if resultInt != intValue {
		t.Fatalf("expected %d, got %d", intValue, resultInt)
	}

	nonExistKey := "notexist"
	_, err = GetAs[string](c, nonExistKey)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCache_StopAndRestart(t *testing.T) {
	c := NewCache()

	key := "testk"
	value := "testv"
	err := c.Set(key, value, time.Second*5)
	if err != nil {
		t.Fatalf("error to set value: %v", err)
	}

	err = c.Stop()
	if err != nil {
		t.Fatalf("error to stop cache: %v", err)
	}

	newValue := "newValue"
	err = c.Set(key, newValue, time.Second*5)
	if !errors.Is(err, ErrStopped) {
		t.Fatalf("expected %v, got %v", ErrStopped, err)
	}

	err = c.Restart()
	if err != nil {
		t.Fatalf("error to restart cache: %v", err)
	}

	err = c.Set(key, newValue, time.Second*5)
	if err != nil {
		t.Fatalf("error to set value after restart: %v", err)
	}

	val, exists := c.Get(key)
	if !exists {
		t.Error("value should exist after restart")
	}
	if val != newValue {
		t.Errorf("expected %v, got %v", newValue, val)
	}

	c.Stop()
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := NewCache()
	defer c.Stop()

	numGoroutines := 10
	numOperations := 100

	errChan := make(chan error, numGoroutines*numOperations)

	wg := sync.WaitGroup{}
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				value := fmt.Sprintf("value_%d_%d", id, j)

				err := c.Set(key, value, time.Second*5)
				if err != nil {
					errChan <- fmt.Errorf("goroutine %d, operation %d: set failed: %v", id, j, err)
					continue
				}

				val, exists := c.Get(key)
				if !exists {
					errChan <- fmt.Errorf("goroutine %d, operation %d: get failed - value doesn't exist", id, j)
					continue
				}
				if val != value {
					errChan <- fmt.Errorf("goroutine %d, operation %d: get failed - expected %s, got %s", id, j, value, val)
					continue
				}

				if !c.Exists(key) {
					errChan <- fmt.Errorf("goroutine %d, operation %d: exists failed - value should exist", id, j)
					continue
				}

				c.Delete(key)
				if c.Exists(key) {
					errChan <- fmt.Errorf("goroutine %d, operation %d: delete failed - value should not exist", id, j)
					continue
				}
			}
		}(i)
	}

	wg.Wait()
	close(errChan)
	for err := range errChan {
		t.Error(err)
	}
}
