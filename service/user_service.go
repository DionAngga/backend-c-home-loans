package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/rysmaadit/go-template/common/errors"
	"github.com/rysmaadit/go-template/config"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/external/jwt_client"
	"github.com/rysmaadit/go-template/external/mysql"
	log "github.com/sirupsen/logrus"
)

type userService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type UserServiceInterface interface {
	Login(user *contract.User) (interface{}, error)
	Create(w http.ResponseWriter, r *http.Request) interface{}
	GetToken(*contract.User) (*contract.GetTokenResponseContract, error)
}

func NewUserService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *userService {
	return &userService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *userService) Create(w http.ResponseWriter, r *http.Request) interface{} {
	payloads, _ := ioutil.ReadAll(r.Body)

	var user contract.User
	json.Unmarshal(payloads, &user)

	user.Login_as = 1

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Create(&user)

	return user
}

func (s *userService) Login(user *contract.User) (interface{}, error) {
	var registeredUser *contract.User

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Table("users").First(&registeredUser, "username = ?", user.Username)

	if user.Username != registeredUser.Username {
		return nil, errors.NewUnauthorizedError("combination of username and password not match, username tidak ada")
	}

	if user.Password != registeredUser.Password {
		return nil, errors.NewUnauthorizedError("combination of username and password not match, password salah")
	}

	tk, err := s.GetToken(registeredUser) //

	if err != nil {
		errMsg := fmt.Sprintf("error di get token: %v", err)
		log.Errorf(errMsg)
		return nil, errors.NewInternalError(err, errMsg)
	}

	return tk, nil
}

func (s *userService) GetToken(user *contract.User) (*contract.GetTokenResponseContract, error) {
	atClaims := contract.JWTMapClaim{
		Authorized: true,
		RequestID:  uuid.New().String(),
		Id_user:    user.ID,
		Login_as:   user.Login_as,
	}

	token, err := s.jwtClient.GenerateTokenStringWithClaims(atClaims, s.appConfig.JWTSecret) //

	if err != nil {
		errMsg := fmt.Sprintf("error signed JWT credentials: %v", err)
		log.Errorf(errMsg)
		return nil, errors.NewInternalError(err, errMsg)
	}

	return &contract.GetTokenResponseContract{Token: token}, err
}
