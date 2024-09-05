package cron

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Action func()

func NewCron(interval time.Duration, action Action, ctx context.Context) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println("-cron: error creating new scheduler:", err)
		return err
	}

	j, err := s.NewJob(gocron.DurationJob(interval), gocron.NewTask(action))
	if err != nil {
		fmt.Println("-cron: error creating new job:", err)
		return err
	}

	fmt.Println("-cron: running job with id:", j.ID())
	s.Start()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM) // Add more "syscall"s ???

	shutdown := func() {
		err = s.Shutdown()
		if err != nil {
			fmt.Println("-cron: error in shutdown", err)
		}
	}

	select {
	case <-signalCh:
		fmt.Println("-cron: gracefully shutdown")
		shutdown()
	case <-ctx.Done():
		fmt.Println("-cron: shutdown with context call")
		shutdown()
	}

	return nil
}
