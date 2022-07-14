package models

import tb "gopkg.in/tucnak/telebot.v2"

type TelegramBot struct {
	Token string `json:"token" redis:"token"`
	URL   string `json:"url" redis:"url"`
}

type TelegramUser struct {
	ID int64 `json:"id" redis:"id"`

	FirstName    string `json:"first_name" redis:"first_name"`
	LastName     string `json:"last_name" redis:"last_name"`
	Username     string `json:"username" redis:"username"`
	LanguageCode string `json:"language_code" redis:"language_code"`
	IsBot        bool   `json:"is_bot" redis:"is_bot"`

	// Returns only in getMe
	CanJoinGroups   bool `json:"can_join_groups" redis:"can_join_groups"`
	CanReadMessages bool `json:"can_read_all_group_messages" redis:"can_read_all_group_messages"`
	SupportsInline  bool `json:"supports_inline_queries" redis:"supports_inline_queries"`
}

func (u *TelegramUser) MapFrom(user *tb.User) *TelegramUser {
	u.ID = user.ID

	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Username = user.Username
	u.LanguageCode = user.LanguageCode
	u.IsBot = user.IsBot

	u.CanJoinGroups = user.CanJoinGroups
	u.CanReadMessages = user.CanReadMessages
	u.SupportsInline = user.SupportsInline

	return u
}

func (u *TelegramUser) MapTo() *tb.User {
	ret := &tb.User{}
	ret.ID = u.ID

	ret.FirstName = u.FirstName
	ret.LastName = u.LastName
	ret.Username = u.Username
	ret.LanguageCode = u.LanguageCode
	ret.IsBot = u.IsBot

	ret.CanJoinGroups = u.CanJoinGroups
	ret.CanReadMessages = u.CanReadMessages
	ret.SupportsInline = u.SupportsInline

	return ret
}
