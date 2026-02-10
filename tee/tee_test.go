package tee

import "testing"

func Test_TeeWritePrintWork(t *testing.T) {
	repsN := 4
	reps := make([]Replica, 0, repsN)

	for i := range repsN {
		reps = append(reps, Replica{i + 1})
	}

	in := make(chan int)

	go func() {
		for i := range 10 {
			in <- i
		}
		close(in)
	}()

	Write(in, reps)
}
