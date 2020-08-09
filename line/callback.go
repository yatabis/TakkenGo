package line

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

func Callback(c echo.Context) error {
	bot := New()
	if bot == nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	events, err := bot.Client.ParseRequest(c.Request())
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Printf("receive the bad request: %e\n", err)
			return c.JSON(http.StatusBadRequest, err)
		} else {
			log.Printf("unexpected error occurs: %e\n", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	for _, event := range events {
		token := event.ReplyToken
		switch event.Type {
		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case "answer":
				bot.ReplyText(token, "指定された問題に解答して、点数を入力してください。\nshortcuts://run-shortcut?name=takken-go\n\nこの機能は未実装です。")
			case "snooze":
				bot.ReplyText(token, "この機能は未実装です。")
			default:
				bot.ReplyOtherPostback(token, event.Postback.Data)
			}
		case linebot.EventTypeMessage:
			if user := event.Source.UserID; user != os.Getenv("USER_ID") {
				bot.ReplyOtherUser(token, user)
				break
			}
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				bot.ReplyOtherText(token, message.Text)
			default:
				bot.ReplyOtherType(token)
			}
		default:
			break
		}
	}
	return c.String(http.StatusOK, "OK")
}
