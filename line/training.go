package line

import (
	"strconv"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"

	"TakkenGo/database"
)

func NewTrainingMessage() *linebot.FlexMessage {
	id, chapter, section := database.GetQuestions()
	text := "【" + chapter + "】\n" + section

	head := linebot.TextComponent{
		Type:     linebot.FlexComponentTypeText,
		Text:     "次の問題に解答してください。",
	}

	question := linebot.TextComponent{
		Type:     linebot.FlexComponentTypeText,
		Text:     text,
		Size:     linebot.FlexTextSizeTypeLg,
		Wrap:     true,
		Weight:   linebot.FlexTextWeightTypeBold,
	}

	body := linebot.BoxComponent{
		Type:     linebot.FlexComponentTypeBox,
		Layout:   linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{&head, &question},
		Spacing:  linebot.FlexComponentSpacingTypeMd,
	}

	answer := TrainingButton("解答", AnswerAction, id, linebot.FlexButtonStyleTypePrimary)
	snooze := TrainingButton("延期", SnoozeAction, id, linebot.FlexButtonStyleTypeSecondary)

	footer := linebot.BoxComponent{
		Type:     linebot.FlexComponentTypeBox,
		Layout:   linebot.FlexBoxLayoutTypeHorizontal,
		Contents: []linebot.FlexComponent{answer, snooze},
		Spacing:  linebot.FlexComponentSpacingTypeLg,
	}

	message := linebot.BubbleContainer{
		Type:     linebot.FlexContainerTypeBubble,
		Body:     &body,
		Footer:   &footer,
	}

	return linebot.NewFlexMessage(text, &message)
}

func TrainingButton(label string, action Action, id int, style linebot.FlexButtonStyleType) *linebot.ButtonComponent {
	data := "action=" + string(action)
	data += "&id=" + strconv.Itoa(id)
	data += "&time=" + strconv.Itoa(time.Now().Hour())
	postback := linebot.NewPostbackAction(label, data, "", "")
	return &linebot.ButtonComponent{
		Type: linebot.FlexComponentTypeButton,
		Action: postback,
		Height: linebot.FlexButtonHeightTypeSm,
		Style: style,
	}
}
