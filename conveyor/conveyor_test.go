package conveyor

import (
	"context"
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	ctx, cl := context.WithCancel(context.Background())
	defer cl()

	naturals := genNaturals(ctx)
	squares := calcSquare(ctx, naturals)

	count := 10
	for range count {
		fmt.Println(<-squares)
	}
}
