package tee

/*
 Задание: Реализация паттерна "Tee" для записи в несколько реплик БД

**Цель задания**:
Реализовать паттерн "Разветвитель", при котором данные из одного источника параллельно записываются в несколько реплик базы данных (имитированных каналами).

---

### Описание задачи

Есть сервис, который записывает данные в кластер БД, состоящий из нескольких реплик. Требуется:
1. Принимать данные из входного канала.
2. Параллельно отправлять их во все реплики (каналы).
3. Гарантировать, что данные записаны во все реплики.
4. Корректно закрыть реплики после завершения работы.

*/

import (
	"fmt"
	"sync"
	"time"
)

// Реплика БД (имитация)
func dbReplica(name string, in <-chan int) {
	for data := range in {
		fmt.Printf("Запись в %s: %d\n", name, data)
		time.Sleep(100 * time.Millisecond) // Имитация задержки записи
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}

type Replica struct {
	id int
}

// TODO create INPUT CANNEL FOR replica, copy to data to every input
func Write(in chan int, replics []Replica) {
	repChs := make([]chan int, 0, len(replics))

	for _, rep := range replics {
		repCh := make(chan int)
		repChs = append(repChs, repCh)

		repName := fmt.Sprintf("replica_id_%d", rep.id)
		go dbReplica(repName, repCh)
	}

	wg := sync.WaitGroup{}
	for data := range in {
		for _, ch := range repChs {
			wg.Add(1)

			go func() {
				defer wg.Done()
				ch <- data
			}()
		}

		wg.Wait()
	}

	for _, ch := range repChs {
		close(ch)
	}
}
