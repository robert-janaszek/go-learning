package course

import (
	"context"
	"fmt"
	"time"
)

func processRequest(ctx context.Context) {
	// reqId := ctx.Value("request_id")

	// switch reqId.(type) {
	// case string:
	// 	fmt.Println("This is string")
	// }

	id, ok := ctx.Value("request_id").(string)

	if ok {
		fmt.Println(id)
	}
}

func run(ctx context.Context) {
	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("stopped", ctx.Err())
			return
		default:
			fmt.Println("work...")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func cancelFunc(cancel func()) {
	time.Sleep(1000 * time.Millisecond)
	cancel()
}

func Day6b() {
	// ex 6
	ctx := context.Background()

	// ex 7
	ctx = context.WithValue(ctx, "request_id", "abc-123")

	processRequest(ctx)

	// ex 8
	ctx, cancel := context.WithCancel(context.Background())

	go cancelFunc(cancel)

	run(ctx)

	// ex 9
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)

	defer cancel()

	run(ctx)
}
