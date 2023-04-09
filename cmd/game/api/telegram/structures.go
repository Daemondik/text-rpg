package telegram_api

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type webhookReqBody struct {
	Message       tgbotapi.Message       `json:"message"`
	CallbackQuery tgbotapi.CallbackQuery `json:"callback_query"`
}

type sendMessageReqBody struct {
	ChatID         int64       `json:"chat_id"`
	Text           string      `json:"text"`
	KeyboardMarkup interface{} `json:"reply_markup,omitempty"`
}

type callbackData struct {
	Action string `json:"action"`
	Value  string `json:"value"`
	ChatId int    `json:"chat_id"`
}

type tgConfig struct {
	Token      string
	WebhookUrl string
}
