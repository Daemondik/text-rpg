package telegram_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strings"
	"wer/cmd/game/helpers"
	"wer/cmd/game/src"
	"wer/cmd/game/src/structures"
)

type Api struct {
}

var bot *tgbotapi.BotAPI
var gameService *src.GameService
var tgConf *tgConfig
var profile *structures.Profile

func (t Api) Start() {
	setConfig()

	var err error
	bot, err = tgbotapi.NewBotAPI(tgConf.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, _ := tgbotapi.NewWebhook(fmt.Sprintf("%s/%s", tgConf.WebhookUrl, bot.Token))

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	gameService = src.NewGameService()

	updates := bot.ListenForWebhook("/" + bot.Token)
	go func() {
		err := http.ListenAndServe(":3000", http.HandlerFunc(handler))
		if err != nil {
			log.Fatal(err)
		}
	}()

	for update := range updates {
		log.Printf("---%+v\n", update)
		//if update.CallbackQuery != nil {
		//	callbackQueryHandler(update.CallbackQuery)
		//	continue
		//}
	}
}

func handler(res http.ResponseWriter, req *http.Request) {
	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	err := handleMsg(body)
	if err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}
}

func callbackQueryHandler(query *tgbotapi.CallbackQuery) {
	fmt.Println(query)
	split := strings.Split(query.Data, ":")
	if split[0] == "selectStory" {
		fmt.Println("selected story", split[1:])
		return
	}
}

func sendMessage(reqBody *sendMessageReqBody) error {
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post("https://api.telegram.org/bot"+tgConf.Token+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}
	return nil
}

func setConfig() {
	tgConf = &tgConfig{
		Token:      helpers.Env("TOKEN_TG", ""),
		WebhookUrl: helpers.Env("WEBHOOK_URL", ""),
	}
}
