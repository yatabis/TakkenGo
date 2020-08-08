package scheduler

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/kawasin73/htask/cron"
)

type Scheduler struct {
	wg *sync.WaitGroup
	cron *cron.Cron
}

func Init() *Scheduler {
	wg := new(sync.WaitGroup)
	c := cron.NewCron(wg, cron.Option{})

	loc := time.FixedZone(os.Getenv("TZ"), 9 * 60 * 60)
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), now.Hour() + 1, 0, 0, 0, loc)

	_, err := c.Every(1).Hour().From(start).Run(training)
	if err != nil {
		log.Printf("failed to set the training task for scheduler: %e\n", err)
	}

	_, err = c.Every(10).Minute().From(start).Run(snooze)
	if err != nil {
		log.Printf("failed to set the snooze task for scheduler: %e\n", err)
	}

	return &Scheduler{wg: wg, cron: c}
}

func (s *Scheduler) Close() {
	err := s.cron.Close()
	if err != nil {
		log.Printf("Error occurs with Scheduler close: %e\n", err)
	}
	s.wg.Wait()
}
