package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/kawasin73/htask/cron"
)

type Scheduler struct {
	wg *sync.WaitGroup
	cron *cron.Cron
}

func Init() *Scheduler {
	time.Local = time.FixedZone("Asia/Tokyo", 9 * 60 * 60)
	wg := new(sync.WaitGroup)
	c := cron.NewCron(wg, cron.Option{})
	s := &Scheduler{wg: wg, cron: c}
	s.Set()
	return s
}

func (s *Scheduler) Set() {
	c := s.cron

	start := time.Now()
	start = start.Add(1 * time.Hour)
	start = start.Add(-time.Duration(start.Minute()) * time.Minute)
	start = start.Add(-time.Duration(start.Second()) * time.Second)
	start = start.Add(-time.Duration(start.Nanosecond()) * time.Nanosecond)

	_, err := c.Every(20).Minute().From(time.Now()).Run(ping)
	if err != nil {
		log.Printf("failed to set the ping task for scheduler: %e\n", err)
	}

	_, err = c.Every(1).Hour().From(start).Run(training)
	if err != nil {
		log.Printf("failed to set the training task for scheduler: %e\n", err)
	}

	_, err = c.Every(10).Minute().From(start).Run(snooze)
	if err != nil {
		log.Printf("failed to set the snooze task for scheduler: %e\n", err)
	}
}

func (s *Scheduler) Close() {
	err := s.cron.Close()
	if err != nil {
		log.Printf("Error occurs with Scheduler close: %e\n", err)
	}
	s.wg.Wait()
}
