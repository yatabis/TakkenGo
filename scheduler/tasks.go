package scheduler

import (
	"time"

	"TakkenGo/line"
)

func training() {
	bot := line.New()
	hour := time.Now().Hour()
	if hour < 9 {
		return
	}
	bot.Training()
}

func snooze() {}
