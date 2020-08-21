package line

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"

	"TakkenGo/database"
)

type ButtonColor string

const (
	ButtonColorRed    ButtonColor = "#FF8080"
	ButtonColorOrange ButtonColor = "#FFBF80"
	ButtonColorYellow ButtonColor = "#FFFF80"
	ButtonColorGreen  ButtonColor = "#80FF80"
	ButtonColorBlue   ButtonColor = "#80BFFF"
	ButtonColorPurple ButtonColor = "#BF80FF"
)

type AnswerButton struct {
	id    int
	time  int
	score int
	color ButtonColor
}

func (c *AnswerButton) MarshalJSON() ([]byte, error) {
	text := linebot.TextComponent{
		Type:       linebot.FlexComponentTypeText,
		Text:       strconv.Itoa(c.score) + "点",
		Align:      linebot.FlexComponentAlignTypeCenter,
		Gravity:    linebot.FlexComponentGravityTypeCenter,
	}
	data := NewPostbackData(ScoreAction, c.id, c.time, c.score).Unmarshal()
	return json.Marshal(&struct {
		Type            string                  `json:"type"`
		Layout          string                  `json:"layout"`
		Contents        []linebot.FlexComponent `json:"contents"`
		BackgroundColor ButtonColor             `json:"backgroundColor"`
		CornerRadius    string                  `json:"cornerRadius"`
		PaddingAll      string                  `json:"paddingAll"`
		Action          linebot.Action          `json:"action"`
	}{
		Type:            "box",
		Layout:          "horizontal",
		Contents:        []linebot.FlexComponent{&text},
		BackgroundColor: c.color,
		CornerRadius:    "md",
		PaddingAll:      "md",
		Action:          linebot.NewPostbackAction("", data, "", ""),
	})
}

func (*AnswerButton) FlexComponent() {}

func ButtonRow(row, id, time int) *linebot.BoxComponent {
	color := []ButtonColor{
		ButtonColorRed,
		ButtonColorOrange,
		ButtonColorYellow,
		ButtonColorGreen,
		ButtonColorBlue,
		ButtonColorPurple,
	}
	return &linebot.BoxComponent{
		Type:            linebot.FlexComponentTypeBox,
		Layout:          linebot.FlexBoxLayoutTypeHorizontal,
		Contents:        []linebot.FlexComponent{
			NewAnswerButton(id, time, 20 * row + 10, color[row]),
			NewAnswerButton(id, time, 20 * row + 20, color[row + 1]),
		},
		Spacing:         linebot.FlexComponentSpacingTypeMd,
	}
}

func NewAnswerButton(id, time, score int, color ButtonColor) *AnswerButton {
	return &AnswerButton{
		id:    id,
		time:  time,
		score: score,
		color: color,
	}
}

func NewQuestionText(text string) *linebot.TextComponent {
	return &linebot.TextComponent{
		Type:     linebot.FlexComponentTypeText,
		Text:     text,
		Size:     linebot.FlexTextSizeTypeLg,
		Wrap:     true,
		Weight:   linebot.FlexTextWeightTypeBold,
	}
}

func NewTrainingButton(label string, action Action, id int, style linebot.FlexButtonStyleType) *linebot.ButtonComponent {
	data := NewPostbackData(action, id, time.Now().Hour(), 0).Unmarshal()
	postback := linebot.NewPostbackAction(label, data, "", "")
	return &linebot.ButtonComponent{
		Type: linebot.FlexComponentTypeButton,
		Action: postback,
		Height: linebot.FlexButtonHeightTypeSm,
		Style: style,
	}
}

func NewTrainingMessage() *linebot.FlexMessage {
	id, chapter, section := database.GetQuestionByRate()
	text := "【" + chapter + "】\n" + section

	head := linebot.TextComponent{
		Type:     linebot.FlexComponentTypeText,
		Text:     "次の問題に解答してください。",
	}

	body := linebot.BoxComponent{
		Type:     linebot.FlexComponentTypeBox,
		Layout:   linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{&head, NewQuestionText(text)},
		Spacing:  linebot.FlexComponentSpacingTypeMd,
	}

	// TODO: あまり有効でない関数の使用 (NewTrainingButton)
	answer := NewTrainingButton("解答", AnswerAction, id, linebot.FlexButtonStyleTypePrimary)
	snooze := &linebot.ButtonComponent{
		Type: linebot.FlexComponentTypeButton,
		Action: linebot.NewURIAction("延期", os.Getenv("ORIGIN") + "/snooze"),
		Height: linebot.FlexButtonHeightTypeSm,
		Style: linebot.FlexButtonStyleTypeSecondary,
	}

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

func NewAnswerMessage(id, time int) *linebot.FlexMessage {
	chapter, section := database.GetQuestionById(id)
	text := "【" + chapter + "】\n" + section
	header := linebot.BoxComponent{
		Type:            linebot.FlexComponentTypeBox,
		Layout:          linebot.FlexBoxLayoutTypeVertical,
		Contents:        []linebot.FlexComponent{NewQuestionText(text)},
	}

	body := linebot.BoxComponent{
		Type:            linebot.FlexComponentTypeBox,
		Layout:          linebot.FlexBoxLayoutTypeVertical,
		Contents:        []linebot.FlexComponent{
			ButtonRow(0, id, time),
			ButtonRow(1, id, time),
			ButtonRow(2, id, time),
			ButtonRow(3, id, time),
			ButtonRow(4, id, time),
		},
		Spacing:         linebot.FlexComponentSpacingTypeMd,
	}

	message := linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Header:    &header,
		Body:      &body,
	}
	return linebot.NewFlexMessage("点数を入力", &message)
}
