package pipeline

import (
	"context"
	"fmt"
)

/*

# Задание: Реализация конвейерной обработки данных (Pipeline паттерн)

**Цель задания**:
Создать конвейер из трех этапов для обработки строковых данных:
1. **Парсинг** — добавление метки "parsed" к данным.
2. **Разделение** — распределение данных между N каналами (round-robin).
3. **Отправка** — параллельная обработка данных в N горутинах с добавлением метки "sent".

---

## Описание задачи

Ваша задача — реализовать систему, которая:
- Обрабатывает данные в строгом порядке: **Parse → Split → Send**.
- Корректно закрывает все каналы после завершения работы.
- Гарантирует потокобезопасность и отсутствие утечек горутин.

### Этапы конвейера

1. **Parse**:
   - Принимает канал сырых данных (`<-chan string`).
   - Добавляет к каждой строке префикс "parsed - ".
   - Возвращает канал обработанных данных.

2. **Split**:
   - Принимает канал данных и число `N` (количество выходных каналов).
   - Распределяет данные между `N` каналами в порядке round-robin.
   - Возвращает слайс каналов (`[]<-chan string`).

3. **Send**:
   - Принимает слайс каналов и запускает `N` горутин.
   - Каждая горутина добавляет к данным префикс "sent - ".
   - Возвращает объединенный канал результатов.
*/

type Processor struct {
	parser  *Parser
	spliter *Spliter
	sender  *Sender
}

func NewProcessor(parser *Parser, spliter *Spliter, sender *Sender) *Processor {
	return &Processor{
		parser:  parser,
		spliter: spliter,
		sender:  sender,
	}
}

type Option struct {
	SplitNum int
}

func (p *Processor) Process(ctx context.Context, input <-chan string, opt Option) {
	if opt.SplitNum <= 0 || opt.SplitNum > defaultSplitLimit {
		opt.SplitNum = defaultSplitNum
	}

	parsedInput := p.parser.Parse(ctx, input)
	splitedInputs := p.spliter.Split(ctx, parsedInput, opt.SplitNum)
	resultsChan := p.sender.Send(ctx, splitedInputs)

	p.handleResultChan(ctx, resultsChan)
}

func (p *Processor) handleResultChan(ctx context.Context, in <-chan string) {
HANDLING_LOOP:
	for {
		select {
		case <-ctx.Done():
			break HANDLING_LOOP
		case s, ok := <-in:
			if !ok {
				break HANDLING_LOOP
			}

			fmt.Println("HANDLED: ", s)
		}
	}

	fmt.Println("PROCESSING CLOSED")
}
