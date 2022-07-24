package model

import "github.com/gin-contrib/sessions"

type AuthRequest struct {
	Username      string           `json:"username" validate:"required,alpha,min=4,max=10"`
	Password      string           `json:"password" validate:"required,min=8"`
	ValueSolution string           `json:"value_solution"`
	ForceLogin    bool             `json:"force_login"`
	Session       sessions.Session `json:"-"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	TokenUuid string
	UserId    uint
	Username  string
	Role      string
}
