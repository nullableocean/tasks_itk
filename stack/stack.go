package stack

/*

Реализация простого обобщенного стека с методами `Push`, `Pop`, `Peek` и `IsEmpty`. Используются дженерики Go (1.18+).

## Требования

1. **Структура `Stack[T]`**:
   - Поле `elements` для хранения элементов (слайс `[]T`).
   - Методы:
        - `NewStack[T]() *Stack[T]`: конструктор.
        - `Push(value T)`: добавление элемента в стек.
        - `Pop() (T, bool)`: удаление и возврат верхнего элемента (с проверкой на пустоту).
        - `Peek() (T, bool)`: возврат верхнего элемента без удаления.
        - `IsEmpty() bool`: проверка стека на пустоту.

2. **Дополнительно**:
   - Гарантировать безопасность операций (например, `Pop` на пустом стеке возвращает `false`).
   - Использовать слайс для эффективного добавления/удаления элементов.
*/

type Stack[T any] struct {
	elems []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		elems: make([]T, 0),
	}
}

func (s *Stack[T]) Push(value T) {
	s.elems = append(s.elems, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	var out T
	if s.IsEmpty() {
		return out, false
	}

	l := len(s.elems)
	out = s.elems[l-1]
	s.elems = s.elems[:l-1]

	return out, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var out T
	if s.IsEmpty() {
		return out, false
	}

	return s.elems[len(s.elems)-1], true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.elems) == 0
}
