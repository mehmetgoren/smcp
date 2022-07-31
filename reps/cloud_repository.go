package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"golang.org/x/oauth2"
	"log"
	"smcp/models"
	"smcp/utils"
	"strconv"
)

type CloudRepository struct {
	Connection *redis.Client
}

func getTelegramUserKey(id string) string {
	return "cloud:telegram:users:" + id
}

func getTelegramBotKey() string {
	return "cloud:telegram:bot"
}

func getTelegramEnabledKey() string {
	return "cloud:telegram:enabled"
}

func (c *CloudRepository) IsTelegramIntegrationEnabled() bool {
	conn := c.Connection
	ctx := context.Background()
	result, err := conn.Get(ctx, getTelegramEnabledKey()).Result()
	if err != nil {
		log.Println(err.Error())
		_, err = conn.Set(ctx, getTelegramEnabledKey(), "0", 0).Result()
		if err != nil {
			log.Println(err.Error())
		}
		return false
	}

	return result == "1"
}

func (c *CloudRepository) GetTelegramUsers() ([]*models.TelegramUser, error) {
	ctx := context.Background()
	ret := make([]*models.TelegramUser, 0)
	conn := c.Connection
	keys, err := conn.Keys(ctx, getTelegramUserKey("*")).Result()
	if err != nil {
		return ret, err
	}
	for _, key := range keys {
		var tbUser models.TelegramUser
		err := conn.HGetAll(ctx, key).Scan(&tbUser)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &tbUser)
	}

	return ret, err
}

func (c *CloudRepository) AddTelegramUser(user *models.TelegramUser) (int64, error) {
	if user == nil {
		return 0, nil
	}
	list, err := c.GetTelegramUsers()
	if err != nil || list == nil {
		return 0, err
	}
	conn := c.Connection
	for _, addedUser := range list {
		if addedUser.ID == user.ID {
			return 0, nil
		}
	}
	result, err := conn.HSet(context.Background(), getTelegramUserKey(strconv.FormatInt(user.ID, 10)), Map(user)).Result()
	list = append(list, user)

	return result, err
}

func (c *CloudRepository) RemoveTelegramUserById(telegramUserId int64) error {
	conn := c.Connection
	_, err := conn.Del(context.Background(), getTelegramUserKey(strconv.FormatInt(telegramUserId, 10))).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudRepository) GetTelegramBot() (*models.TelegramBot, error) {
	conn := c.Connection
	bot := &models.TelegramBot{}
	err := conn.HGetAll(context.Background(), getTelegramBotKey()).Scan(bot)
	return bot, err
}

// ********************* GDrive *********************

func getGdriveEnabledKey() string {
	return "cloud:gdrive:enabled"
}

func getGdriveTokenKey() string {
	return "cloud:gdrive:token"
}

func getGdriveCredentialsKey() string {
	return "cloud:gdrive:credentials"
}

func getGdriveUrlKey() string {
	return "cloud:gdrive:url"
}

func getGdriveAuthCodeKey() string {
	return "cloud:gdrive:authcode"
}

func (c *CloudRepository) getValue(key string) (string, error) {
	result, err := c.Connection.Get(context.Background(), key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			_, err := c.Connection.Set(context.Background(), key, "", 0).Result()
			return "", err
		} else {
			log.Println(err.Error())
			return "", err
		}
	}

	return result, nil
}
func (c *CloudRepository) setValue(key string, value string) (string, error) {
	return c.Connection.Set(context.Background(), key, value, 0).Result()
}

func (c *CloudRepository) IsGdriveIntegrationEnabled() bool {
	conn := c.Connection
	ctx := context.Background()
	result, err := conn.Get(ctx, getGdriveEnabledKey()).Result()
	if err != nil {
		log.Println(err.Error())
		_, err = conn.Set(ctx, getGdriveEnabledKey(), "0", 0).Result()
		if err != nil {
			log.Println(err.Error())
		}
		return false
	}

	return result == "1"
}

func (c *CloudRepository) GetGdriveToken() (*oauth2.Token, error) {
	js, err := c.getValue(getGdriveTokenKey())
	if err != nil {
		return nil, err
	}
	var ret oauth2.Token
	err = utils.DeserializeJson(js, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *CloudRepository) SaveGdriveToken(token *oauth2.Token) error {
	js, err := utils.SerializeJson(token)
	if err != nil {
		return err
	}
	_, err = c.setValue(getGdriveTokenKey(), js)
	return err
}

func (c *CloudRepository) GetGdriveCredentials() (string, error) {
	return c.getValue(getGdriveCredentialsKey())
}

func (c *CloudRepository) SaveGdriveUrl(url string) error {
	_, err := c.setValue(getGdriveUrlKey(), url)
	return err
}

func (c *CloudRepository) GetGdriveAuthCode() (string, error) {
	return c.getValue(getGdriveAuthCodeKey())
}

func (c *CloudRepository) SaveGdriveAuthCode(authCode string) error {
	_, err := c.setValue(getGdriveAuthCodeKey(), authCode)
	return err
}
