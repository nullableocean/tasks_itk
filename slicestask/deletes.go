package slicestask

/*

### SLICE TASK 3
# 1. Курс Go: Удаление элементов из слайса

Это задание поможет освоить работу со слайсами в Go, фокусируясь на операциях удаления элементов с учетом эффективности, порядка и управления памятью.

## Цели
- Научиться удалять элементы из слайса с сохранением и без сохранения порядка.
- Понять, как избежать утечек памяти при работе с указателями.
- Оптимизировать использование памяти слайса.
- Реализовать продвинутые операции (удаление дубликатов, фильтрация).

## Задание

### Часть 1: Базовое удаление
Реализуйте функции:
- `RemoveUnordered(s []T, i int) []T` — удаление без сохранения порядка.
- `RemoveOrdered(s []T, i int) []T` — удаление с сохранением порядка.

### Часть 2: Удаление по значению
Реализуйте функцию:
- `RemoveAllByValue(s []T, value T) []T` — удаление всех вхождений `value`.

### Часть 3: Работа с памятью
1. Обнуляйте удаленные элементы-указатели.
2. Сокращайте вместимость (`capacity`) слайса при сильном уменьшении размера.

### Часть 4: Дополнительные задачи
1. `RemoveDuplicates(s []T) []T` — удаление дубликатов.
2. `RemoveIf(s []T, predicate func(T) bool) []T` — удаление по условию.
package main



*/

// RemoveUnordered удаляет элемент по индексу без сохранения порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveUnordered[T any](s []T, i int) []T {
	if i >= len(s) {
		return s
	}

	newS := append(s[:i], s[i+1:]...)
	return newS
}

// RemoveOrdered удаляет элемент по индексу с сохранением порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveOrdered[T any](s []T, i int) []T {
	if i >= len(s) {
		return s
	}

	clear(s[i : i+1])
	return s
}

// RemoveAllByValue удаляет все вхождения указанного значения.
func RemoveAllByValue[T comparable](s []T, value T) []T {
	newS := make([]T, 0, len(s))

	for _, v := range s {
		if v != value {
			newS = append(newS, v)
		}
	}

	return newS[:len(newS):len(newS)]
}

// RemoveDuplicates оставляет только уникальные элементы (сохраняет порядок).
func RemoveDuplicates[T comparable](s []T) []T {
	vals := make(map[T]struct{}, len(s))
	duples := make(map[T]struct{}, 0)

	for _, v := range s {
		if _, ex := vals[v]; ex {
			duples[v] = struct{}{}
		} else {
			vals[v] = struct{}{}
		}
	}

	newS := make([]T, 0, len(s))
	for _, v := range s {
		if _, ex := duples[v]; !ex {
			newS = append(newS, v)
		}
	}

	return newS[:len(newS):len(newS)]
}

// RemoveIf удаляет элементы, удовлетворяющие условию predicate.
func RemoveIf[T any](s []T, predicate func(T) bool) []T {
	newS := make([]T, 0, len(s))

	for _, v := range s {
		if !predicate(v) {
			newS = append(newS, v)
		}
	}

	return newS[:len(newS):len(newS)]
}

// RemoveOrderedWithNil удаляет элемент по индексу (для слайса указателей),
// обнуляя удаляемый элемент для предотвращения утечек памяти.
func RemoveOrderedWithNil[T any](s []*T, i int) []*T {
	if i >= len(s) {
		return s
	}

	s[i] = nil

	return s
}

// ShrinkCapacity сокращает вместимость слайса, если она превышает
// удвоенную длину после удаления элементов.
func ShrinkCapacity[T any](s []T) []T {
	if cap(s) < len(s)*2 {
		return s
	}

	return s[:len(s):len(s)]
}
