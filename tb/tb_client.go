package tb

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"smcp/models"
	"smcp/reps"
	"time"
)

type TelegramBotClient struct {
	DbBot           *models.TelegramBot
	Bot             *tb.Bot
	CloudRepository *reps.CloudRepository
}

func CreateTelegramBot(cr *reps.CloudRepository) (*TelegramBotClient, error) {
	if cr == nil {
		return nil, errors.New("insufficient arguments")
	}

	users, err := cr.GetTelegramUsers()
	log.Println(users)

	dbBot, err := cr.GetTelegramBot()
	if err != nil || dbBot == nil || len(dbBot.Token) == 0 {
		if err == nil {
			err = errors.New("no token has been found, exiting")
		}
		return nil, err
	}

	bot, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL:    "https://api.telegram.org",
		Token:  dbBot.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Println("a telegram bot couldn't be created")
		log.Println(err)
		return nil, errors.Unwrap(err)
	} else {
		bot.Handle("/start", func(m *tb.Message) {
			dbUser := &models.TelegramUser{}
			result, err := cr.AddTelegramUser(dbUser.MapFrom(m.Sender))
			if err != nil {
				log.Println(err.Error())
				return
			}
			if err == nil && result == 0 {
				_, err = bot.Send(m.Sender, "You've already been registered.")
			} else {
				if err != nil {
					log.Println(err.Error())
				}
				_, err = bot.Send(m.Sender, "Feniks® Telegram© Cloud Integration has been completed.")
			}
			if err != nil {
				log.Println(err.Error())
			}
			log.Println("a telegram bot has been created")
		})
		bot.Handle("/stop", func(m *tb.Message) {
			_, err := bot.Send(m.Sender, "You will never get any message from me anymore.")
			if err != nil {
				log.Println(err.Error())
				return
			}
			err = cr.RemoveTelegramUserById(m.Sender.ID)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Println("a telegram bot has been removed")
		})
		go func() {
			log.Println("telegram bot has been started -> " + bot.Me.Username)
			bot.Start()
		}()
	}

	return &TelegramBotClient{Bot: bot, DbBot: dbBot, CloudRepository: cr}, err
}
