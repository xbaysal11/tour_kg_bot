package main

import (
	"strconv"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	TELEGRAM_APITOKEN      = "840467541:AAFlgUSKpv4ZvK2jKsaVF0SgSvsTh248iO8"
	StepFrom          Step = iota
	StepTo
)

var step = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🛩 Начать"),
		tgbotapi.NewKeyboardButton("❌ Выйти"),
	),
)

var step1 = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Бишкек"),
		tgbotapi.NewKeyboardButton("Ош"),
		tgbotapi.NewKeyboardButton("Ыссык-Кол"),
	),
)

var step2 = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Бишкек"),
		tgbotapi.NewKeyboardButton("Ош"),
		tgbotapi.NewKeyboardButton("Ыссык-Кол"),
	),
)

type DateRange struct {
	Start string
	End   string
}
type Step int

type Conversation struct {
	User   *tgbotapi.User
	Step   Step
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
		// if update.Message == nil {
		// 	continue
		// }

		// if !update.Message.IsCommand() {
		// 	continue
		// }

		User := update.Message.From
		UserName := User.FirstName
		UserID := User.ID
		ChatID := strconv.FormatInt(update.Message.Chat.ID, 10)

		//
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		conv := conversations[UserID]

		if conv != nil {
			// conversation already exist
			// switch update.Message.Command() {
			// case "from":
			// 	msg.Text = "Пожалуйста, выберите откуда собираетесь вылететь:"
			// 	// msg.ReplyMarkup = step2
			// case "stop":
			// 	delete(conversations, UserID)
			// 	msg.Text = "Разговор был успешно удален. Чтобы начать разговор использую комманду /start. До встречи " + UserName
			// default:
			// 	msg.Text = "Разговор уже есть, если хочешь его отменить набери: /stop"
			// }
			if update.Message.IsCommand() {
				cmdText := update.Message.Command()
				if cmdText == "stop" {
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					msg.Text = "✅ Разговор был успешно удален.\n До встречи " + UserName + "\n\nЧтобы начать новый разговор набери: /start."
					delete(conversations, UserID)
				} else if cmdText == "from" {
					msg.Text = "Пожалуйста, выберите откуда собираетесь вылететь: "
					msg.ReplyMarkup = step1
				} else {
					msg.Text = "Разговор уже существует!\nДля того чтобы отменить набери: /stop"
				}
			} else {
				if update.Message.Text == step.Keyboard[0][0].Text {
					msg.Text = "Отлично!\nНаберите: `/from` "
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				} else if update.Message.Text == step.Keyboard[0][1].Text {
					msg.Text = "✅ Разговор был успешно удален.\n До встречи " + UserName + "\n\nЧтобы начать новый разговор набери: /start."
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					delete(conversations, UserID)
				} else if update.Message.Text == step1.Keyboard[0][0].Text {
					msg.Text = "🔵 Tour From: " + step1.Keyboard[0][0].Text + "\n🔴 Tour To: "

					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					// conv.Step = StepFrom
				} else {
					// other messages
					msg.Text = "ok"
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					conv, ok := conversations[UserID]
					if ok {

						if conv.Step == StepFrom {
							conv.From = update.Message.Text
							msg.Text = "Введите телефон:"
							conv.Step = 1
							msg.ReplyMarkup = step1
						} else if conv.Step == StepTo {
							conv.To = update.Message.Text
							conv.Step = 2
							msg.Text = "Введите course:"
							msg.ReplyMarkup = step2
						}
					} else {
						// other messages
						msg.Text = "ok"
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					}
				}
			}
		} else {
			if update.Message.IsCommand() {
				cmdText := update.Message.Command()
				if cmdText == "start" {
					conversations[UserID] = NewConversation(User)
					msg.Text = "✋ Здравствуй, " + UserName + ".\nНовый разговор был создан ChatID: " + ChatID
					msg.ReplyMarkup = step
				} else {
					msg.Text = "Я телеграм бот 🤖.\nЯ не знаю такой комманды.\nТы можешь начать разговор использую комманду /start"
				}
			}
		}
		// send
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
