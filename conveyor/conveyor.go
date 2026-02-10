package conveyor

/*
**Цель задания**:
Реализовать конвейер для обработки чисел с использованием горутин и каналов. Числа из первого канала должны читаться по мере поступления, обрабатываться (например, возводиться в квадрат) и записываться во второй канал.

### Описание задачи

Даны два канала:
- `naturals` (для передачи исходных чисел),
- `squares` (для передачи обработанных чисел).

Необходимо:
1. **Генерировать** числа и отправлять их в канал `naturals`.
2. **Читать** числа из `naturals`, обрабатывать их (возводить в квадрат) и отправлять результат в `squares`.
3. **Выводить** результаты из `squares` в консоль.
*/

import (
	"context"
	"math/rand"
)

func genNaturals(ctx context.Context) chan int {
	outCh := make(chan int)
	naturalLimit := 10000

	go func() {
		defer close(outCh)
		for {
			select {
			case <-ctx.Done():
				return
			case outCh <- rand.Intn(naturalLimit):
			}
		}
	}()

	return outCh
}

func calcSquare(ctx context.Context, in <-chan int) chan int {
	outCh := make(chan int)

	go func() {
		defer close(outCh)
		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
					return
				case outCh <- n * n:
				}
			}
		}
	}()

	return outCh
}
