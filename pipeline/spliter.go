package pipeline

import "context"

var (
	defaultSplitLimit = 1000
	defaultSplitNum   = 4
)

type Spliter struct{}

func NewSpliter() *Spliter {
	return &Spliter{}
}

// splited channels have buffer size = 1
func (spliter *Spliter) Split(ctx context.Context, in <-chan string, num int) []<-chan string {
	chans := make([]chan string, 0, num)

	for range num {
		ch := make(chan string, 1)
		chans = append(chans, ch)
	}

	closeFn := func() {
		for _, ch := range chans {
			close(ch)
		}
	}

	go func() {
		defer closeFn()

		ind := 0
		for {
			select {
			case <-ctx.Done():
				return
			case s, ok := <-in:
				if !ok {
					return
				}

				// можем отправлять в вых. канал в горутине и использовать вайт группу, только нужно добавить лимиты на такие горутины
				select {
				case <-ctx.Done():
					return
				case chans[ind] <- s:
				}

				ind = (ind + 1) % num
			}
		}
	}()

	outChans := make([]<-chan string, num)
	for i, ch := range chans {
		outChans[i] = ch
	}

	return outChans
}
