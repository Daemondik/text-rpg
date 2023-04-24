package telegram_api

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	structures "wer/cmd/game/src/structures"
)

const (
	commStart         = "start"
	commStories       = "stories"
	actionSelectStory = "selectStory"
)

func handleMsg(body *webhookReqBody) error {
	if len(body.Message.Text) == 0 && len(body.CallbackQuery.Data) == 0 {
		return nil
	}
	var err error
	if len(body.Message.Text) > 0 {

		if body.Message.Text[1:] != commStart {
			ps := &structures.ProfileSource{
				SourceID: int(body.Message.From.ID),
				Source:   structures.UserTG,
			}
			profile, err = gameService.ProfileBySource(*ps)
			if err != nil {
				return err
			}
			gameService.SetProfile(profile)
		}

		if body.Message.Text[:1] == "/" {
			switch body.Message.Text[1:] {
			case commStart:
				err := startBot(body)
				if err != nil {
					return err
				}
				err = storiesList(body)
				if err != nil {
					return err
				}
			case commStories:
				err := storiesList(body)
				if err != nil {
					return err
				}
			default:
				fmt.Println("default")
			}
		} else {
			if profile.ProfileSource.SourceState == structures.SourceStateReadingStory {
				storyChoice, err := strconv.Atoi(body.Message.Text)
				if err != nil {
					return err
				}
				progress, err := gameService.LastProfileProgress()
				if err != nil {
					return err
				}

				if len(progress.StoryLine.StoryLineChoices) >= storyChoice {
					for choiceCount, choice := range progress.StoryLine.StoryLineChoices {
						choiceCount++
						if choiceCount == storyChoice {
							progress.Choice = &choice
						}
					}
					err = gameService.UpdateLastProfileProgress(*progress)

					err = showNextStoryLine(body, *progress.Choice)
					if err != nil {
						return err
					}
				} else {
					return errors.New("нет такого варианта ответа")
				}
			} else {
				return errors.New("неправильное состаяние клиента")
			}
		}
	}
	if len(body.CallbackQuery.Data) > 0 {
		callbackData := &callbackData{}
		err := json.Unmarshal([]byte(body.CallbackQuery.Data), callbackData)
		if err != nil {
			return err
		}

		ps := &structures.ProfileSource{
			SourceID: callbackData.ChatId,
			Source:   structures.UserTG,
		}
		profile, err = gameService.ProfileBySource(*ps)
		if err != nil {
			return err
		}
		gameService.SetProfile(profile)

		switch callbackData.Action {
		case actionSelectStory:
			err := selectStory(callbackData)
			if err != nil {
				return err
			}
			err = showLastStoryLine(callbackData)
			if err != nil {
				return err
			}

		default:
			fmt.Println("default")
		}
	}
	return nil
}

func startBot(body *webhookReqBody) error {
	ps := structures.ProfileSource{
		SourceID:    int(body.Message.From.ID),
		Source:      structures.UserTG,
		SourceState: structures.SourceStateSelectingStory,
	}
	var err error
	profile, err = gameService.CreateProfileIfNotExist(ps)
	if err != nil {
		return err
	}
	return nil
}

func storiesList(body *webhookReqBody) error {
	stories, err := gameService.StoriesList()
	if err != nil {
		return err
	}

	err = gameService.UpdateProfileState(structures.SourceStateSelectingStory)
	if err != nil {
		return err
	}

	var text string
	var btns [][]tgbotapi.InlineKeyboardButton
	for counter, story := range stories {
		counter++
		counterStr := strconv.Itoa(counter)
		text += fmt.Sprintf("%s. %s\n", counterStr, story.Name)
		callbackData := callbackData{
			Action: actionSelectStory,
			Value:  strconv.FormatUint(uint64(story.ID), 10),
			ChatId: int(body.Message.Chat.ID),
		}
		callbackDataJson, err := json.Marshal(callbackData)
		if err != nil {
			return err
		}
		btns = append(btns, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(counterStr, string(callbackDataJson))))
	}

	reqBody := &sendMessageReqBody{
		ChatID: body.Message.Chat.ID,
		Text:   text,
		KeyboardMarkup: &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: btns,
		},
	}

	err = sendMessage(reqBody)
	if err != nil {
		return err
	}

	return nil
}

func selectStory(callback *callbackData) error {
	storyId, err := strconv.ParseUint(callback.Value, 10, 32)
	if err != nil {
		return err
	}
	ps := structures.ProfileSource{
		SourceID:    callback.ChatId,
		Source:      structures.UserTG,
		SourceState: structures.SourceStateReadingStory,
	}
	err = gameService.SelectStory(uint(storyId), ps)
	if err != nil {
		return err
	}
	return nil
}

func showLastStoryLine(callback *callbackData) error {
	sl, err := gameService.LastStoryLine(*profile)
	if err != nil {
		return err
	}

	var btns [][]tgbotapi.KeyboardButton
	text := sl.Text + "\n"

	for count, choices := range sl.StoryLineChoices {
		count++
		text += fmt.Sprintf("\n%d. %s", count, choices.Text)
		btns = append(btns, tgbotapi.NewKeyboardButtonRow(tgbotapi.KeyboardButton{
			Text: strconv.Itoa(count),
		}))
	}

	reqBody := &sendMessageReqBody{
		ChatID: int64(callback.ChatId),
		Text:   text,
		KeyboardMarkup: &tgbotapi.ReplyKeyboardMarkup{
			Keyboard:        btns,
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
	}

	err = sendMessage(reqBody)
	if err != nil {
		return err
	}

	return nil
}

func showNextStoryLine(body *webhookReqBody, slc structures.StoryLineChoice) error {
	sl, err := gameService.NextStoryLineByChoice(slc)
	if err != nil {
		return err
	}

	pp := structures.ProfileProgress{
		Profile:   *profile,
		Story:     sl.Story,
		StoryLine: *sl,
	}

	err = gameService.AddProfileProgress(pp)
	if err != nil {
		return err
	}

	var btns [][]tgbotapi.KeyboardButton
	text := sl.Text + "\n"

	for count, choices := range sl.StoryLineChoices {
		count++
		text += fmt.Sprintf("\n%d. %s", count, choices.Text)
		btns = append(btns, tgbotapi.NewKeyboardButtonRow(tgbotapi.KeyboardButton{
			Text: strconv.Itoa(count),
		}))
	}

	reqBody := &sendMessageReqBody{
		ChatID: body.Message.From.ID,
		Text:   text,
		KeyboardMarkup: &tgbotapi.ReplyKeyboardMarkup{
			Keyboard:        btns,
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
	}

	err = sendMessage(reqBody)
	if err != nil {
		return err
	}

	return nil
}
