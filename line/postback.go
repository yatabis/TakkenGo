package line

import (
	"net/url"
	"strconv"
)

type Action string

const (
	AnswerAction Action = "answer"
	SnoozeAction Action = "snooze"
	ScoreAction  Action = "score"
)

type PostbackData struct {
	action     Action
	questionId int
	time       int
	score      int
}

func ParsePostback(postback string) (data PostbackData) {
	query, err := url.ParseQuery(postback)
	if err != nil {
		return
	}

	action := getValue(query, "action")
	if action == "" {
		return
	}

	questionId := getValueInt(query, "id")
	if questionId < 1 || questionId > 50 {
		return
	}

	if Action(action) == ScoreAction {
		score := getValueInt(query, "score")
		if 10 <= score && score <= 100 && score % 10 == 0 {
			data.score = score
		} else {
			return
		}
	}

	data.action = Action(action)
	data.questionId = questionId

	t := getValueInt(query, "time")
	if t >= 9 {
		data.time = t
	}
	return
}

func getValue(data map[string][]string, key string) string {
	value, ok := data[key]
	if !ok || len(value) != 1 {
		return ""
	}
	return value[0]
}

func getValueInt(data map[string][]string, key string) int {
	value := getValue(data, key)
	if value == "" {
		return 0
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return num
}
