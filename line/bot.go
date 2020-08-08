package line

import (
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LINE struct {
	Client *linebot.Client
}

func New() *LINE {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Printf("failed to create new line bot.: %e\n", err)
		return nil
	}
	return &LINE{Client: bot}
}

func (b *LINE) Training() {
	user := os.Getenv("USER_ID")
	_, err := b.Client.PushMessage(user, linebot.NewTextMessage("テスト")).Do()
	if err != nil {
		log.Printf("failed to send training: %e\n", err)
	}
}

func (b *LINE) ReplyOtherType(token string) {
	text := "テキスト以外のメッセージタイプには対応していません。"
	res, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(text)).Do()
	if err == nil {
		log.Printf("replied to other message type: %+v\n", res)
	} else {
		log.Printf("failed to replying to other message type: %e\n", err)
	}
}

func (b *LINE) ReplyOtherText(token string) {
	text := "テキストを受け取りました。\n解答を登録する場合は、解答する問題の「解答」ボタンを押してから点数を入力してください。"
	res, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(text)).Do()
	if err == nil {
		log.Printf("replied to other text message: %+v\n", res)
	} else {
		log.Printf("failed to replying to other text message: %e\n", err)
	}
}
