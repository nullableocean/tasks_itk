package slicestask

/*

### Задача

Реализуйте структуру стека с использованием слайсов, удовлетворяющую следующему интерфейсу:

```go
type Stacker interface {
    Push(v int)
    Pop() int
}
```

### Требования к реализации

1. Операция Push(v int)
    Должна добавлять целочисленное значение v в стек.

2. Операция Pop() int Должна возвращать последний добавленный элемент, реализуя поведение LIFO (последним пришёл — первым ушёл).
    Если стек пуст, вызов метода Pop() должен приводить к панике.

3. Конструктор
    Реализуйте функцию New() *stack, возвращающую новый экземпляр стека.

4. Реализация должна находится в main.go
5. Реализация должна успешно проходить тесты. Для их запуска введите команду `go test ./...` в этой директории

*/

type Stacker interface {
	Push(v int)
	Pop() int
}

type stack struct {
	lifo []int
}

func (s *stack) Push(v int) {
	s.lifo = append(s.lifo, v)
}

func (s *stack) Pop() int {
	l := len(s.lifo)
	if l == 0 {
		panic("empty stack")
	}

	val := s.lifo[l-1]
	s.lifo = s.lifo[:l-1]

	return val
}

func NewStack() *stack {
	return &stack{
		lifo: make([]int, 0),
	}
}
