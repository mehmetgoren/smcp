package rd

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

type RedisOptions struct {
	Client *redis.Client
}

type RedisRepository struct {
	*RedisOptions
}

var redisKeyUsers = "users"

func (r *RedisRepository) GetAllUsers() []*tb.User {
	jsonList, err := r.Client.Get(context.Background(), redisKeyUsers).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			emptyList := make([]*tb.User, 0)
			r.Client.Set(context.Background(), redisKeyUsers, emptyList, 0)
			return emptyList
		} else {
			log.Println(err.Error())
			return nil
		}
	}
	var list []*tb.User
	if len(jsonList) == 0 {
		list = make([]*tb.User, 0)
	} else {
		err := json.Unmarshal([]byte(jsonList), &list)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
	}

	return list
}

func (r *RedisRepository) sendList(list []*tb.User) {
	jsonBytes, _ := json.Marshal(list)
	jsonList := string(jsonBytes)
	r.Client.Set(context.Background(), redisKeyUsers, jsonList, 0)
}

func (r *RedisRepository) AddUser(user *tb.User) *RedisRepository {
	if user == nil {
		return r
	}
	list := r.GetAllUsers()
	if list == nil {
		return r
	}
	for _, addedUser := range list {
		if addedUser.ID == user.ID {
			return r
		}
	}
	list = append(list, user)
	r.sendList(list)

	return r
}

func (r *RedisRepository) RemoveUser(user *tb.User) *RedisRepository {
	if user == nil {
		return r
	}
	list := r.GetAllUsers()
	if list == nil {
		return r
	}
	removeFn := func(s []*tb.User, i int) []*tb.User {
		s[i] = s[len(s)-1]
		return s[:len(s)-1]
	}

	userIndex := -1
	for i, currentUser := range list {
		if currentUser.ID == user.ID {
			userIndex = i
			break
		}
	}
	if userIndex > -1 {
		list = removeFn(list, userIndex)
		r.sendList(list)
	}

	return r
}

func (r *RedisRepository) GetValue(key string) (string, error) {
	result, err := r.Client.Get(context.Background(), key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			r.Client.Set(context.Background(), key, "", 0)
			return "", err
		} else {
			log.Println(err.Error())
			return "", err
		}
	}

	return result, nil
}

func (r *RedisRepository) SetValue(key string, value string) *redis.StatusCmd {
	return r.Client.Set(context.Background(), key, value, 0)
}
