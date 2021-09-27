package main

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

type TelegramBotOptions struct {
	Token string //"1944447440:AAF8C0vJ2rjd__9CWT7PVcg9cON8QixdAMs"
	Repository *RedisRepository
}

type TelegramBotClient struct {
	Bot    *tb.Bot
	repository *RedisRepository
}

func empty() TelegramBotClient {
	return TelegramBotClient{}
}

func(t TelegramBotClient) GetUsers() []*tb.User {
	return t.repository.GetAllUsers()
}

func CreateTelegramBot(token string, rep *RedisRepository) (TelegramBotClient, error) {
	if len(token) == 0 || rep == nil {
		return empty(), errors.New("insufficient arguments")
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
			//
		})
		go func() {
			log.Println("telegram bot has been started -> " + bot.Me.Username)
			bot.Start()
		}()
	}

	return TelegramBotClient{Bot: bot, repository: rep}, err
}
