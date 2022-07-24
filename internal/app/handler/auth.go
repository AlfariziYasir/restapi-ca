package handler

import (
	"net/http"
	"restapi/internal/app/model"
	"restapi/internal/app/service"
	"restapi/internal/constant"
	"restapi/internal/security/token"
	"restapi/internal/validation"
	"restapi/internal/web"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Refresh(c *gin.Context)
	Logout(c *gin.Context)
}

type authHandler struct {
	authService service.AuthService
	tk          token.TokenInterface
}

func NewAuthHandler(authService service.AuthService, tk token.TokenInterface) AuthHandler {
	return &authHandler{authService, tk}
}

func (h *authHandler) Login(c *gin.Context) {
	var req model.AuthRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.MarshalError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = validation.Struct(req)
	if err != nil {
		web.MarshalError(c, http.StatusConflict, err.Error(), nil)
		return
	}

	if ok, m := validation.CheckCaptchaSolver(req.ValueSolution, sessions.Default(c)); !ok {
		web.MarshalError(c, http.StatusBadRequest, m, "captcha")
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		switch err {
		case constant.ErrUserNameNotRegistered, constant.ErrWrongPassword:
			web.MarshalError(c, http.StatusUnauthorized, err.Error(), nil)
		default:
			web.MarshalError(c, http.StatusInternalServerError, err.Error(), nil)
		}

		c.Abort()
		return
	}

	web.MarshalPayload(c, http.StatusOK, "login successfully", res)
}

func (h *authHandler) Refresh(c *gin.Context) {
	t, err := h.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		web.MarshalError(c, http.StatusInternalServerError, "failed to refresh token", nil)
		c.Abort()
		return
	}

	req := model.AccessDetails{
		TokenUuid: t.TokenUuid,
		Username:  t.Username,
		UserId:    t.UserId,
		Role:      t.Role,
	}
	res, err := h.authService.Refresh(req)
	if err != nil {
		web.MarshalError(c, http.StatusUnauthorized, err.Error(), nil)
		c.Abort()
		return
	}

	web.MarshalPayload(c, http.StatusOK, "refresh token successfully", res)
}

func (h *authHandler) Logout(c *gin.Context) {
	metadata, err := h.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		web.MarshalError(c, http.StatusUnauthorized, err.Error(), nil)
		c.Abort()
		return
	} else if metadata != nil {
		err = h.authService.Logout(metadata)
		if err != nil {
			web.MarshalError(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		}
	}

	c.Writer.Header().Del("username")
	c.Writer.Header().Del("user_id")

	web.MarshalPayload(c, http.StatusOK, "success logout", nil)
}
