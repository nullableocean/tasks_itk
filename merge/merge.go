package merge

/*

**Цель задания**:
Написать функцию `mergeChannels`, которая объединяет данные из нескольких каналов в один общий канал, используя паттерн `FAN-IN`.

---

### Описание задачи

Дано:
- `n` каналов типа `<-chan int`.
- Функция должна вернуть канал `<-chan int`, в который попадают все значения из исходных каналов.

Требования:
1. Все значения из входных каналов должны быть отправлены в выходной канал.
2. Выходной канал должен быть закрыт после завершения всех входных каналов.
3. Решение должно быть потокобезопасным и эффективным.

*/

import (
	"sync"
)

func mergeChannels(channels ...<-chan int) <-chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, ch := range channels {
		go func() {
			defer wg.Done()

			for n := range ch {
				out <- n
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
