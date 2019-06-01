package main

import (
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const TELEGRAM_APITOKEN = "840467541:AAFlgUSKpv4ZvK2jKsaVF0SgSvsTh248iO8"

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Help"),
		tgbotapi.NewKeyboardButton("Say Hi"),
		tgbotapi.NewKeyboardButton("Status"),
	),
)

type DateRange struct {
	Start string
	End   string
}

type Conversation struct {
	User   *tgbotapi.User
	Step   int
	From   string
	To     string
	Depart *DateRange
	Return *DateRange
}

func NewConversation(User *tgbotapi.User) *Conversation {
	return &Conversation{
		User: User,
	}
}

var conversations = map[int]*Conversation{}

func main() {
	bot, err := tgbotapi.NewBotAPI(TELEGRAM_APITOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Error("Произошла ошибка во время получения обновлений")
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}
		//
		User := update.Message.From
		UserID := User.ID
		//
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		conv := conversations[UserID]

		if conv != nil {
			// conversation already exist
			switch update.Message.Command() {
			case "stop":
				delete(conversations, UserID)
				msg.Text = "Разговор был успешно удален"
			default:
				msg.Text = "Разговор уже есть, если хочешь его отменить набери: /stop"
			}

		} else {
			//
			switch update.Message.Command() {
			case "start":
				conversations[UserID] = NewConversation(User)
				msg.Text = "Новый разговор был создан"
			default:
				msg.Text = "Я не знаю такой комманды. Ты можешь начать разговор использую комманду /start"
			}
		}

		// send
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
