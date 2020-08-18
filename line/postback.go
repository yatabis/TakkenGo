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

func NewPostbackData(action Action, questionId, time, score int) *PostbackData {
	return &PostbackData{
		action:     action,
		questionId: questionId,
		time:       time,
		score:      score,
	}
}

func (p *PostbackData) Unmarshal() (data string) {
	if p.action == "" || p.questionId == 0 {
		return
	}
	if p.action == ScoreAction && p.score == 0 {
		return
	}
	data += "action=" + string(p.action) + "&id=" + strconv.Itoa(p.questionId)
	if p.time >= 9 {
		data += "&time=" + strconv.Itoa(p.time)
	}
	if p.action == ScoreAction {
		if 10 <= p.score && p.score <= 100 && p.score % 10 == 0 {
			data += "&score=" + strconv.Itoa(p.score)
		} else {
			return ""
		}
	}
	return
}

func ParsePostbackData(postback string) (data PostbackData) {
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
