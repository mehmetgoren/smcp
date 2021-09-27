package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"strings"
	"time"
)

type RedisOptions struct {
	Context context.Context
	Client *redis.Client
}

type RedisRepository struct {
	*RedisOptions
	Key string
}

func (r RedisRepository) GetAllUsers() []*tb.User{
	jsonList, err := r.Client.Get(r.Context, r.Key).Result()
	if err != nil{
		if err.Error() == "redis: nil"{
			emptyList := make([]*tb.User, 0)
			r.Client.Set(r.Context, r.Key, emptyList, 0)
			return emptyList
		}else{
			log.Println(err.Error())
			return nil
		}
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

func (r RedisRepository) AddUser(user *tb.User)  *RedisRepository{
	if user == nil{
		return &r
	}
	list := r.GetAllUsers()
	if list == nil{
		return &r
	}
	for _, addedUser := range list{
		if addedUser.ID == user.ID{
			return &r
		}
	}
	list = append(list, user)
	jsonBytes, _ := json.Marshal(list)
	jsonList := string(jsonBytes)
	r.Client.Set(r.Context, r.Key, jsonList, 0)

	return &r
}

func toMyFormat(t *time.Time)  string{
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t.Year()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(int(t.Month())))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Day()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Hour()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Minute()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Second()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Nanosecond()))

	return sb.String()
}

func (r RedisRepository) Heartbeat(time *time.Time){
	r.Client.Set(r.Context, "heartbeat_smcp", toMyFormat(time), 0)
}