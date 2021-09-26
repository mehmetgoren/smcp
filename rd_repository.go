package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

type RedisOptions struct {
	Context context.Context
	Client *redis.Client
}

type RedisRepository struct {
	*RedisOptions
	Key string
}

func (r RedisRepository) GetAll() []*tb.User{
	jsonList, err := r.Client.Get(r.Context, r.Key).Result()
	if err != nil{
		log.Println(err.Error())
		return nil
	}
	var list []*tb.User
	if len(jsonList) == 0{
		list = make([]*tb.User, 0)
	}else{
		err := json.Unmarshal([]byte(jsonList), &list)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
	}

	return list
}

func (r RedisRepository) Add(user *tb.User)  *RedisRepository{
	list := r.GetAll()
	list = append(list, user)
	jsonBytes, _ := json.Marshal(list)
	jsonList := string(jsonBytes)
	r.Client.Set(r.Context, r.Key, jsonList, 0)

	return &r
}