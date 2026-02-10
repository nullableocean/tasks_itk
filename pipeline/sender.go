package pipeline

import (
	"context"
	"sync"
)

var (
	defaultWorkersLimits = 5
	boundForStartWorkers = 100
)

type sender struct {
	prefix string
}

func (s *sender) send(data string) string {
	return s.prefix + data
}

type Sender struct {
	*sender

	workersLimit int
	wg           sync.WaitGroup
}

func NewSender(workersLimit int) *Sender {
	return &Sender{
		sender: &sender{prefix: "sent - "},

		workersLimit: workersLimit,
		wg:           sync.WaitGroup{},
	}
}

func (sender *Sender) Send(ctx context.Context, inChans []<-chan string) <-chan string {
	out := make(chan string)

	// если входных каналов много, создаём ограниченное количество воркеров и распределяем между ними каналы
	// или запускаем для  каждого канала отдельную горутину
	if len(inChans) >= boundForStartWorkers {
		sender.startWorkers(ctx, inChans, out)
	} else {
		for _, input := range inChans {
			sender.wg.Add(1)
			go sender.sendSingleJob(ctx, input, out)
		}
	}

	go func() {
		sender.wg.Wait()
		close(out)
	}()

	return out
}

func (sender *Sender) startWorkers(ctx context.Context, inChans []<-chan string, out chan<- string) {
	inputsLen := len(inChans)
	remain := inputsLen % sender.workersLimit
	batchCount := inputsLen / sender.workersLimit

	start := 0
	end := batchCount

	for range sender.workersLimit {
		if end+remain >= inputsLen {
			end = inputsLen
		}

		sender.wg.Add(1)
		go sender.sendWorker(ctx, inChans[start:end], out)

		start = end
		end += batchCount
	}
}

func (sender *Sender) sendWorker(ctx context.Context, inChans []<-chan string, out chan<- string) {
	defer sender.wg.Done()

	inputsLen := len(inChans)
	closed := make(map[int]struct{}, inputsLen)
	ind := -1

WORK_LOOP:
	for {
		ind = (ind + 1) % inputsLen

		if len(closed) == inputsLen {
			return
		}
		if _, ex := closed[ind]; ex {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case s, ok := <-inChans[ind]:
			if !ok {
				closed[ind] = struct{}{}
				continue WORK_LOOP
			}

			sentVal := sender.send(s)
			select {
			case <-ctx.Done():
				return
			case out <- sentVal:
			}
		default:
		}
	}
}

func (sender *Sender) sendSingleJob(ctx context.Context, in <-chan string, out chan<- string) {
	defer sender.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case s, ok := <-in:
			if !ok {
				return
			}

			sentVal := sender.prefix + s
			select {
			case <-ctx.Done():
				return
			case out <- sentVal:
			}
		}
	}
}
