package line

import (
	"log"
	"net/http"

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
		if event.Type != linebot.EventTypeMessage {
			break
		}
		switch event.Message.(type) {
		case *linebot.TextMessage:
			bot.ReplyOtherText(token)
		default:
			bot.ReplyOtherType(token)
		}
	}
	return c.String(http.StatusOK, "OK")
}
