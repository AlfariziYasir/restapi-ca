package service

import (
	"errors"
	"restapi/internal/app/model"
	"restapi/internal/app/repository"
	"restapi/internal/constant"
	"restapi/internal/logger"
	"restapi/internal/security/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(req model.AuthRequest) (*model.AuthResponse, error)
	Refresh(req model.AccessDetails) (*model.AuthResponse, error)
	Logout(metaData *model.AccessDetails) error
}

func NewAuthService(
	userRepo repository.UserRepo,
	authRepo repository.AuthRepo,
	tk token.TokenInterface) AuthService {
	return &authService{userRepo, authRepo, tk}
}

type authService struct {
	userRepo repository.UserRepo
	authRepo repository.AuthRepo
	tk       token.TokenInterface
}

func (s *authService) Login(req model.AuthRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		logger.Log().Err(err).Msg("failed to get user by username")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, constant.ErrUserNameNotRegistered
		default:
			return nil, constant.ErrServer
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	if user.IsLogin && !req.ForceLogin {
		logger.Log().Err(errors.New("try to force login")).Msg("user is already logged in another device")
		return nil, errors.New("user is already logged in another device")
	} else if req.ForceLogin {
		metaData := &model.AccessDetails{
			TokenUuid: user.TokenUuid,
			Username:  user.Username,
			UserId:    user.ID,
		}

		err = s.authRepo.DeleteTokens(metaData)
		if err != nil {

			logger.Log().Err(err).Msg("failed to force login")
			return nil, err
		}
	}

	claims := map[string]interface{}{
		"user_id":   user.ID,
		"username":  user.Username,
		"user_role": user.Role,
	}
	ts, err := s.tk.CreateToken(claims)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create token")
		return nil, err
	}

	err = s.authRepo.CreateAuth(claims, ts)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create auth")
		return nil, err
	}

	user.IsLogin = true
	user.TokenUuid = ts.TokenUuid
	err = s.userRepo.Update(user)
	if err != nil {
		logger.Log().Err(err).Msg("failed to login")
		return nil, constant.ErrServer
	}

	res := &model.AuthResponse{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	}

	return res, nil
}

func (s *authService) Refresh(req model.AccessDetails) (*model.AuthResponse, error) {
	td, err := s.authRepo.FetchAuth(req.TokenUuid)
	if err != nil {
		logger.Log().Err(errors.New("token not valid")).Msg("error fetch auth with prev token")
		return nil, errors.New("refresh token not valid")
	}

	accDetail := &model.AccessDetails{
		TokenUuid: req.TokenUuid,
		UserId:    req.UserId,
		Username:  req.Username,
		Role:      req.Role,
	}
	err = s.authRepo.DeleteTokens(accDetail)
	if err != nil {
		logger.Log().Err(errors.New("token not valid")).Msg("refresh token not valid")
		return nil, errors.New("refresh token not valid")
	}

	data := map[string]interface{}{
		"user_id":   td["user_id"],
		"username":  td["username"],
		"user_role": req.Role,
	}
	ts, err := s.tk.CreateToken(data)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create new refresh token")
		return nil, err
	}

	err = s.authRepo.CreateAuth(data, ts)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create new refresh auth")
		return nil, err
	}

	user, err := s.userRepo.Get(uint(td["user_id"].(float64)))
	if err != nil {
		logger.Log().Err(err).Msg("failed to get user by id for new refresh auth")
		return nil, err
	}

	user.TokenUuid = ts.TokenUuid
	err = s.userRepo.Update(user)
	if err != nil {
		logger.Log().Err(err).Msg("failed update user for new refresh token")
		return nil, constant.ErrServer
	}

	res := &model.AuthResponse{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	}

	return res, nil
}

func (s *authService) Logout(metaData *model.AccessDetails) error {
	err := s.authRepo.DeleteTokens(metaData)
	if err != nil {
		logger.Log().Err(err).Msg("failed to logout")
		if err != nil {
			return err
		}
	}

	user, err := s.userRepo.GetByUsername(metaData.Username)
	if err != nil {
		logger.Log().Err(err).Msg("failed to logout")
		return err
	}

	user.IsLogin = false
	user.TokenUuid = ""
	err = s.userRepo.Update(user)
	if err != nil {
		logger.Log().Err(err).Msg("failed to logout")
		return err
	}

	return nil
}
