package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"restapi/internal/app/model"
	"restapi/internal/db/postgres"
	"restapi/internal/db/redis"
	"time"
)

type UserRepo interface {
	Create(user *model.User) error
	Get(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
}

type userRepo struct {
	pg  postgres.Client
	rds redis.Client
}

func NewUserRepo(pg postgres.Client, rds redis.Client) UserRepo {
	return &userRepo{pg, rds}
}

func (r *userRepo) Create(user *model.User) error {
	err := r.pg.Conn().Create(&user).Select("id").Scan(&user.ID).Error
	if err != nil {
		return err
	}

	temp, err := r.Get(user.ID)
	if err != nil {
		return err
	}

	*user = *temp
	return nil
}

func (r *userRepo) Get(id uint) (*model.User, error) {
	user := new(model.User)

	str, err := r.rds.Conn().Get(context.Background(), fmt.Sprintf("user_id:%v", id)).Result()
	if err == nil {
		json.Unmarshal([]byte(str), &user)
		return user, nil
	}

	err = r.pg.Conn().First(&user, id).Error
	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(user)
	_, err = r.rds.Conn().Set(context.Background(), fmt.Sprintf("user_id:%v", id), b, time.Duration(1*time.Hour)).Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) GetByUsername(username string) (*model.User, error) {
	user := new(model.User)

	err := r.pg.Conn().Where(&model.User{
		Username: username,
	}).First(&user).Error
	if err != nil {
		return nil, err
	}

	temp, err := r.Get(user.ID)
	if err != nil {
		return nil, err
	}

	*user = *temp
	return user, nil
}

func (r *userRepo) Update(user *model.User) error {
	err := r.pg.Conn().Model(&model.User{}).Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"is_login":   user.IsLogin,
			"token_uuid": user.TokenUuid,
			"password":   user.Password,
		}).Error
	if err != nil {
		return err
	}

	_, err = r.rds.Conn().Del(context.Background(), fmt.Sprintf("user_id:%v", user.ID)).Result()
	if err != nil {
		return err
	}

	temp, err := r.Get(user.ID)
	if err != nil {
		return err
	}

	*user = *temp
	return nil
}

func (r *userRepo) Delete(id uint) error {
	err := r.pg.Conn().Delete(&model.User{}, id).Error
	if err != nil {
		return err
	}

	_, err = r.rds.Conn().Del(context.Background(), fmt.Sprintf("user_id:%v", id)).Result()
	if err != nil {
		return err
	}
	return nil
}
