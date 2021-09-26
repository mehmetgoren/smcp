package main

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

type SenderOptions struct {
	Id           int    //1608330992
	Username     string //ionian_gokalp
	LanguageCode string //en
}

type TelegramBotOptions struct {
	Token string //"1944447440:AAF8C0vJ2rjd__9CWT7PVcg9cON8QixdAMs"
}

type TelegramBotClient struct {
	Bot    *tb.Bot
	Sender *tb.User
}

func empty() TelegramBotClient {
	return TelegramBotClient{}
}

func CreateTelegramBot(bi *TelegramBotOptions, si *SenderOptions) (TelegramBotClient, error) {
	if bi == nil || si == nil {
		return empty(), errors.New("insufficient arguments")
	}
	bot, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL:    "https://api.telegram.org",
		Token:  bi.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	var sender *tb.User = nil
	if err != nil {
		log.Println("a telegram bot couldn't be created")
		log.Println(err)
		return empty(), errors.Unwrap(err)
	} else {
		log.Println("a telegram bot has been created")
		sender = &tb.User{ID: si.Id, Username: si.Username, LanguageCode: si.LanguageCode}
		bot.Handle("/start", func(m *tb.Message) {
			sender = m.Sender
			bot.Send(m.Sender, "Camera capturing is now starting")
			//
		})
		go func() {
			log.Println("telegram bot has been started -> " + bot.Me.Username)
			bot.Start()
		}()
	}

	return TelegramBotClient{Bot: bot, Sender: sender}, err
}
