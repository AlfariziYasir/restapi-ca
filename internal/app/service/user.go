package service

import (
	"encoding/json"
	"restapi/internal/app/model"
	"restapi/internal/app/repository"
	"restapi/internal/config"
	"restapi/internal/constant"
	"restapi/internal/logger"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(req model.UserCreateRequest) (*model.UserResponse, error)
	Get(id uint) (*model.UserResponse, error)
	List(req model.RequestDataTable) (*model.UserList, error)
	Update(req model.UserUpdateRequest) (*model.UserResponse, error)
	UpdatePassword(req model.UserPasswordUpdateRequest) (*model.UserResponse, error)
	Delete(id uint) error
}

type userService struct {
	userRepo   repository.UserRepo
	customRepo repository.CustomRepo
}

func NewUserService(
	userRepo repository.UserRepo,
	CustomRepo repository.CustomRepo,
) UserService {
	return &userService{userRepo, CustomRepo}
}

func (s *userService) Create(req model.UserCreateRequest) (*model.UserResponse, error) {
	_, err := s.userRepo.GetByUsername(req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Log().Err(err).Msg("failed to get user by username")
		return nil, constant.ErrServer
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log().Err(err).Msg("failed to generate from password")
		return nil, constant.ErrServer
	}

	user := &model.User{
		Username: req.Username,
		Password: string(password),
		Role:     req.UserRole,
	}
	err = s.userRepo.Create(user)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create account")
		return nil, constant.ErrServer
	}

	user, err = s.userRepo.GetByUsername(req.Username)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create account")
		return nil, err
	}

	return model.NewUserResponse(user), nil
}

func (s *userService) Get(id uint) (*model.UserResponse, error) {
	user, err := s.userRepo.Get(id)
	if err != nil {
		logger.Log().Err(err).Msg("failed to get user by id")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, constant.ErrUserNotFound
		default:
			return nil, err
		}
	}

	res := model.NewUserResponse(user)

	return res, nil
}

func (s *userService) List(req model.RequestDataTable) (*model.UserList, error) {
	tableName := config.Cfg().DatabaseSchemaUser + ".users"
	data, err := s.customRepo.List(req, model.UserResponse{}, tableName)
	if err != nil {
		logger.Log().Err(err).Msg("failed to get list users")
		return nil, err
	}

	users := make([]*model.User, 0)
	b, _ := json.Marshal(data.Data)
	json.Unmarshal(b, &users)

	return model.NewUserListResponse(users, data.Count), nil
}

func (s *userService) Update(req model.UserUpdateRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Log().Err(err).Msg("failed to get user by username")
		return nil, constant.ErrServer
	} else if err == nil && user.ID != req.ID {
		return nil, constant.ErrEmailRegistered
	}

	user, err = s.userRepo.Get(req.ID)
	if err != nil {
		logger.Log().Err(err).Msg("failed to get user by id")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, constant.ErrUserNotFound
		default:
			return nil, constant.ErrServer
		}
	}

	user.Username = req.Username
	err = s.userRepo.Update(user)
	if err != nil {
		logger.Log().Err(err).Msg("failed to update user")
		return nil, constant.ErrServer
	}

	return model.NewUserResponse(user), nil
}

func (s *userService) UpdatePassword(req model.UserPasswordUpdateRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.Get(req.ID)
	if err != nil {
		logger.Log().Err(err).Msg("failed to get user by id")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, constant.ErrUserNotFound
		default:
			return nil, constant.ErrServer
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		logger.Log().Err(err).Msg("wrong password")
		return nil, constant.ErrWrongPassword
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log().Err(err).Msg("failed to generate from password")
		return nil, constant.ErrServer
	}

	user.Password = string(password)
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(user)
	if err != nil {
		logger.Log().Err(err).Msg("failed to update user password")
		return nil, constant.ErrServer
	}

	return model.NewUserResponse(user), nil
}

func (s *userService) Delete(id uint) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		logger.Log().Err(err).Msg("failed to delete user")
		return constant.ErrServer
	}

	return nil
}
