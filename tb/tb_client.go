package tb

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"smcp/reps"
	"time"
)

type TelegramBotOptions struct {
	Token      string //"1944447440:AAF8C0vJ2rjd__9CWT7PVcg9cON8QixdAMs"
	Repository *reps.RedisRepository
}

type TelegramBotClient struct {
	Bot        *tb.Bot
	Repository *reps.RedisRepository
}

func empty() TelegramBotClient {
	return TelegramBotClient{}
}

func (t TelegramBotClient) GetUsers() []*tb.User {
	return t.Repository.GetAllUsers()
}

func CreateTelegramBot(rep *reps.RedisRepository) (TelegramBotClient, error) {
	if rep == nil {
		return empty(), errors.New("insufficient arguments")
	}

	token, err := rep.GetValue("telegram_bot_token")
	if err != nil || len(token) == 0 || rep == nil {
		if err == nil {
			err = errors.New("no token has been found, exiting...")
		}
		return empty(), err
	}

	bot, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL:    "https://api.telegram.org",
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Println("a telegram bot couldn't be created")
		log.Println(err)
		return empty(), errors.Unwrap(err)
	} else {
		log.Println("a telegram bot has been created")
		bot.Handle("/start", func(m *tb.Message) {
			rep.AddUser(m.Sender)
			bot.Send(m.Sender, "Camera capturing is now starting")
		})
		bot.Handle("/stop", func(m *tb.Message) {
			bot.Send(m.Sender, "You will never get any message from me")
			rep.RemoveUser(m.Sender)
		})
		go func() {
			log.Println("telegram bot has been started -> " + bot.Me.Username)
			bot.Start()
		}()
	}

	return TelegramBotClient{Bot: bot, Repository: rep}, err
}
