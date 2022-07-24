package token

import (
	"errors"
	"fmt"
	"net/http"
	"restapi/internal/app/model"
	"restapi/internal/config"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

type tokenservice struct{}

func NewToken() *tokenservice {
	return &tokenservice{}
}

type TokenInterface interface {
	CreateToken(data map[string]interface{}) (*model.TokenDetails, error)
	ExtractTokenMetadata(*http.Request) (*model.AccessDetails, error)
}

//Token implements the TokenInterface
var _ TokenInterface = &tokenservice{}

func (t *tokenservice) CreateToken(data map[string]interface{}) (*model.TokenDetails, error) {
	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 4).Unix() //expires after 30 min
	tokenUuid, _ := uuid.NewV4()

	td.TokenUuid = tokenUuid.String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = fmt.Sprintf("%s++%v%v", td.TokenUuid, data["user_id"], data["username"])

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.TokenUuid
	atClaims["user_id"] = data["user_id"]
	atClaims["username"] = data["username"]
	atClaims["user_role"] = data["user_role"]
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(config.Cfg().JwtSecretKey))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = fmt.Sprintf("%s++%v%v", td.TokenUuid, data["user_id"], data["username"])

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = data["user_id"]
	rtClaims["username"] = data["username"]
	rtClaims["user_role"] = data["user_role"]
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(config.Cfg().JwtRefreshKey))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func TokenValid(r *http.Request) (*model.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("token invalid")
	}

	data, err := extract(token)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Cfg().JwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//get the token from the request body
func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func extract(token *jwt.Token) (*model.AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		userId, userOk := claims["user_id"].(float64)
		username, usernameOk := claims["username"].(string)
		role, roleOk := claims["user_role"].(string)

		if !ok && !userOk && !usernameOk && !roleOk {
			return nil, errors.New("unauthorized")
		} else {
			return &model.AccessDetails{
				TokenUuid: accessUuid,
				UserId:    uint(userId),
				Username:  username,
				Role:      role,
			}, nil
		}
	}
	return nil, errors.New("something went wrong")
}

func (t *tokenservice) ExtractTokenMetadata(r *http.Request) (*model.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	acc, err := extract(token)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
