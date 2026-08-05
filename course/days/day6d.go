package days

import (
	"context"
	"log/slog"
	"os"
	"time"
)

func Day6d() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	// ex 16, 17
	slog.Info("user logged in", "user_id", 42, "role", "admin")

	// ex 18
	t1 := time.Now()

	t2 := time.Date(2012, time.February, 6, 12, 11, 40, 12, time.Local)

	diff := t2.Sub(t1)
	slog.Info("time difference", "diff", diff.Abs().Seconds())

	// ex 19
	slog.Info("Current time", "time", t1.Format("2006-01-02"))

	// ex 20
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()

	for i := 0; i < 5; i++ {
		select {
		case <-ticker.C:
			slog.Info("ticker")
		case <-ctx.Done():
			slog.Info("ticker stopped")
			return
		}
	}

}
