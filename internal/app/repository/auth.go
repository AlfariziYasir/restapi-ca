package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"restapi/internal/app/model"
	"restapi/internal/db/redis"
	"time"
)

type AuthRepo interface {
	CreateAuth(map[string]interface{}, *model.TokenDetails) error
	FetchAuth(tokenUuid string) (map[string]interface{}, error)
	DeleteRefresh(string) error
	DeleteTokens(*model.AccessDetails) error
}

type authRepo struct {
	redisClient redis.Client
}

func NewAuthRepo(redisClient redis.Client) AuthRepo {
	return &authRepo{redisClient}
}

func (r *authRepo) CreateAuth(authD map[string]interface{}, td *model.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	b, _ := json.Marshal(authD)
	atCreated, err := r.redisClient.Conn().Set(context.Background(), td.TokenUuid, b, at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := r.redisClient.Conn().Set(context.Background(), td.RefreshUuid, b, rt.Sub(now)).Result()
	if err != nil {
		return err
	}

	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

func (r *authRepo) FetchAuth(tokenUuid string) (map[string]interface{}, error) {
	authD, err := r.redisClient.Conn().Get(context.Background(), tokenUuid).Result()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}

	json.Unmarshal([]byte(authD), &data)

	return data, nil
}

func (r *authRepo) DeleteTokens(authD *model.AccessDetails) error {
	refreshUuid := fmt.Sprintf("%s++%v%s", authD.TokenUuid, authD.UserId, authD.Username)
	//delete access token
	_, err := r.redisClient.Conn().Del(context.Background(), authD.TokenUuid).Result()
	if err != nil {
		return err
	}

	_, err = r.redisClient.Conn().Del(context.Background(), refreshUuid).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepo) DeleteRefresh(refreshUuid string) error {
	deleted, err := r.redisClient.Conn().Del(context.Background(), refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
