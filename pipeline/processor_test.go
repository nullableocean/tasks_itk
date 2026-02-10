package pipeline

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_ProcessorWork(t *testing.T) {
	processor := NewProcessor(NewParser(), NewSpliter(), NewSender(defaultWorkersLimits))

	ctx, cl := context.WithCancel(context.Background())
	defer cl()

	input := make(chan string)
	go func() {
		for i := range 20 {
			testMessage := fmt.Sprintf("message_id_%d", i)
			input <- testMessage
		}

		close(input)
	}()

	fmt.Println("DEFAULT WORK")
	processor.Process(ctx, input, Option{})
	fmt.Println()
}

func Test_ProcessContextTimeline(t *testing.T) {
	processor := NewProcessor(NewParser(), NewSpliter(), NewSender(defaultWorkersLimits))

	ctx, cl := context.WithTimeout(context.Background(), time.Second*2)
	defer cl()

	input := make(chan string)
	go func() {
		for i := range 20 {
			testMessage := fmt.Sprintf("message_id_%d", i)
			input <- testMessage
			time.Sleep(time.Millisecond * 500)
		}

		close(input)
	}()

	fmt.Println("WITH CONTEXT TIMEOUT")
	processor.Process(ctx, input, Option{})
}
