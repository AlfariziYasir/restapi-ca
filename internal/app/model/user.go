package model

import (
	"restapi/internal/config"
	"time"

	"gorm.io/gorm"
)

type User struct {
	CreatedAt time.Time      `gorm:"column:create_on"`
	UpdatedAt time.Time      `gorm:"column:change_on"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	No        int64          `json:"no" datatable:"-" gorm:"-"`
	ID        uint           `gorm:"primaryKey;index;NOT NULL;column:id;autoIncrement"`
	Username  string         `gorm:"type:varchar(20);NOT NULL;UNIQUE;index"`
	Password  string         `gorm:"type:varchar(255)"`
	Role      string         `gorm:"type:varchar(5)"`
	IsLogin   bool           `gorm:"column:is_login"`
	TokenUuid string         `gorm:"column:token_uuid"`
}

func (u *User) TableName() string {
	// custom table name, this is default
	return config.Cfg().DatabaseSchemaUser + ".users"
}

type UserCreateRequest struct {
	Username   string `json:"username" validate:"required,alpha,min=4,max=10"`
	Password   string `json:"password" validate:"required,min=8"`
	RePassword string `json:"repassword" validate:"required,max=20,min=8,eqfield=Password"`
	UserRole   string
}

type UserGetRequest struct {
	ID uint
}

type UserUpdateRequest struct {
	ID                uint   `json:"-"`
	Username          string `json:"username" validate:"required,alpha,min=4,max=10"`
	LoketID           string `json:"loket_id"`
	LoketPembayaranID string `json:"loket_pembayaran_id"`
	IsLogin           bool   `json:"is_login"`
}

type UserPasswordUpdateRequest struct {
	ID            uint   `json:"-"`
	OldPassword   string `json:"old_password" validate:"required,min=8"`
	NewPassword   string `json:"new_password" validate:"required,min=8"`
	ReNewPassword string `json:"renew_password" validate:"required,max=20,min=8,eqfield=NewPassword"`
}

type UserDeleteRequest struct {
	ID uint
}

type UserResponse struct {
	No       int64  `json:"no" datatable:"-"`
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"roles"`
}

type UserList struct {
	Count int             `json:"count"`
	Data  []*UserResponse `json:"data"`
}

func NewUserResponse(payload *User) *UserResponse {
	return &UserResponse{
		No:       payload.No,
		ID:       payload.ID,
		Username: payload.Username,
		Role:     payload.Role,
	}
}

func NewUserListResponse(payloads []*User, count int64) *UserList {
	res := make([]*UserResponse, len(payloads))
	for i, payload := range payloads {
		res[i] = NewUserResponse(payload)
	}

	return &UserList{
		Count: int(count),
		Data:  res,
	}
}
