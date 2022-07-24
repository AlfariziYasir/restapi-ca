package handler

import (
	"fmt"
	"restapi/internal/app/model"
	"restapi/internal/app/service"
	"restapi/internal/constant"
	"restapi/internal/validation"
	"restapi/internal/web"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetByToken(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	UpdatePassword(c *gin.Context)
	Delete(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService}
}

func (h *userHandler) Create(c *gin.Context) {

	var req model.UserCreateRequest

	role := web.GetUrlQueryString(c, "role")

	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.MarshalError(c, http.StatusBadRequest, constant.ErrUnauthorized.Error(), nil)
		c.Abort()
		return
	}

	err = validation.Struct(req)
	if err != nil {
		web.MarshalError(c, http.StatusConflict, err.Error(), nil)
		c.Abort()
		return
	}

	req.UserRole = role
	res, err := h.userService.Create(req)
	if err != nil {
		switch err {
		case constant.ErrEmailRegistered:
			web.MarshalError(c, http.StatusConflict, err.Error(), nil)
			c.Abort()
			return
		default:
			web.MarshalError(c, http.StatusInternalServerError, err.Error(), nil)
			c.Abort()
			return
		}
	}

	web.MarshalPayload(c, http.StatusOK, "create user is successfully", res)
}

func (h *userHandler) Get(c *gin.Context) {
	id, err := web.GetUrlQueryInt64(c, "id")
	if err != nil {
		web.MarshalError(c, http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	res, err := h.userService.Get(uint(id))
	if err != nil {
		switch err {
		case constant.ErrUserNotFound:
			web.MarshalError(c, http.StatusNotFound, err.Error(), nil)
			c.Abort()
			return
		default:
			web.MarshalError(c, http.StatusInternalServerError, err.Error(), nil)
			c.Abort()
			return
		}
	}

	web.MarshalPayload(c, http.StatusOK, "get data is success", res)
}

func (h *userHandler) GetByToken(c *gin.Context) {
	id := c.MustGet("user_id").(uint)

	res, err := h.userService.Get(id)
	fmt.Println("error get user:", err)
	if err != nil {
		switch err {
		case constant.ErrUserNotFound:
			web.MarshalError(c, http.StatusNotFound, err.Error(), nil)
			c.Abort()
			return
		default:
			web.MarshalError(c, http.StatusInternalServerError, err.Error(), nil)
			c.Abort()
			return
		}
	}

	web.MarshalPayload(c, http.StatusOK, "get data is success", res)
}

func (h *userHandler) List(c *gin.Context) {
	req := model.RequestDataTable{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		web.MarshalError(c, http.StatusBadRequest, "please check your data", nil)
		c.Abort()
		return
	}

	res, err := h.userService.List(req)
	if err != nil {
		web.MarshalError(c, http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return
	}

	web.MarshalPayload(c, http.StatusOK, "success get list users", res)
}

func (h *userHandler) Update(c *gin.Context) {
	id, err := web.GetUrlQueryInt64(c, "id")
	if err != nil {
		web.MarshalError(c, http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	req := model.UserUpdateRequest{ID: uint(id)}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		web.MarshalError(c, http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	err = validation.Struct(req)
	if err != nil {
		web.MarshalError(c, http.StatusConflict, err.Error(), nil)
		c.Abort()
		return
	}

	res, err := h.userService.Update(req)
	if err != nil {
		switch err {
		case constant.ErrUnauthorized:
			web.MarshalError(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		case constant.ErrEmailRegistered:
			web.MarshalError(c, http.StatusConflict, err.Error(), nil)
			c.Abort()
			return
		case constant.ErrUserNotFound:
			web.MarshalError(c, http.StatusNotFound, err.Error(), nil)
			c.Abort()
			return
		default:
			web.MarshalError(c, http.StatusInternalServerError, err.Error(), nil)
			c.Abort()
			return
		}
	}

	web.MarshalPayload(c, http.StatusOK, "update user is success", res)
}

func (h *userHandler) UpdatePassword(c *gin.Context) {
	id, err := web.GetUrlQueryInt64(c, "id")
	if err != nil {
		web.MarshalError(c, http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	req := model.UserPasswordUpdateRequest{ID: uint(id)}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		web.MarshalError(c, http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	err = validation.Struct(req)
	if err != nil {
		web.MarshalError(c, http.StatusConflict, err.Error(), nil)
		c.Abort()
		return
	}

	res, err := h.userService.UpdatePassword(req)
	if err != nil {
		switch err {
		case constant.ErrUnauthorized, constant.ErrWrongPassword:
			web.MarshalError(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		case constant.ErrUserNotFound:
			web.MarshalError(c, http.StatusNotFound, err.Error(), nil)
			c.Abort()
			return
		default:
			web.MarshalError(c, http.StatusInternalServerError, err.Error(), nil)
			c.Abort()
			return
		}
	}

	web.MarshalPayload(c, http.StatusOK, "update password user is success", res)
}

func (h *userHandler) Delete(c *gin.Context) {
	id, err := web.GetUrlQueryInt64(c, "id")
	if err != nil {
		web.MarshalError(c, http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	err = h.userService.Delete(uint(id))
	if err != nil {
		switch err {
		case constant.ErrUnauthorized:
			web.MarshalError(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		case constant.ErrUserNotFound:
			web.MarshalError(c, http.StatusNotFound, err.Error(), nil)
			c.Abort()
			return
		default:
			web.MarshalError(c, http.StatusInternalServerError, err.Error(), nil)
			c.Abort()
			return
		}
	}

	web.MarshalPayload(c, http.StatusOK, "delete user is success", nil)
}
