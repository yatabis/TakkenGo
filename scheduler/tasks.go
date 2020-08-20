package scheduler

import (
	"log"
	"time"

	"gopkg.in/resty.v1"

	"TakkenGo/database"
	"TakkenGo/line"
)

func ping() {
	r, err := resty.R().Get("https://takken-go.herokuapp.com/ping")
	if err == nil {
		log.Printf("ping: %+v\n", r)
	} else {
		log.Printf("ping error: %e\n", err)
	}
}

func reset() {
	database.ResetScores()
}

func training() {
	bot := line.New()
	hour := time.Now().Hour()
	if hour < 9 {
		return
	}
	bot.Training()
}

func snooze() {}
