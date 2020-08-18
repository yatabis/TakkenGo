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
	_, err := b.Client.PushMessage(user, NewTrainingMessage()).Do()
	if err != nil {
		log.Printf("failed to send training: %e\n", err)
	}
}

func (b *LINE) ReplyOtherPostback(token, data string) {
	log.Printf("received the other postback event: %s\n", data)
	text := "不明なポストバックイベントを受け取りました。"
	b.ReplyText(token, text)
}

func (b *LINE) ReplyOtherUser(token, user string) {
	log.Printf("received messages from the other user: %s\n", user)
	text := "ユーザー登録をすべての機能をお使いいただけます。\n\nなお、現在はユーザー登録を受け付けておりません。"
	b.ReplyText(token, text)
}

func (b *LINE) ReplyOtherType(token string) {
	log.Printf("received the other message.\n")
	text := "テキスト以外のメッセージタイプには対応していません。"
	b.ReplyText(token, text)
}

func (b *LINE) ReplyOtherText(token, mes string) {
	log.Printf("received the other text message: %s\n", mes)
	text := "テキストを受け取りました。\n解答を登録する場合は、解答する問題の「解答」ボタンを押してから点数を入力してください。"
	b.ReplyText(token, text)
}

func (b *LINE) ReplyText(token, text string) {
	_, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(text)).Do()
	if err != nil {
		log.Printf("failed to reply text message: %e\n", err)
	}
}

func (b *LINE) ReplyMessages(token string, messages ...linebot.SendingMessage) {
	_, err := b.Client.ReplyMessage(token, messages...).Do()
	if err != nil {
		log.Printf("failed to reply messages: %e\n", err)
	}
}
