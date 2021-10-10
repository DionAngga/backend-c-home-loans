package service

import (
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/rysmaadit/go-template/common/errors"
	"github.com/rysmaadit/go-template/config"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/external/jwt_client"
	"github.com/rysmaadit/go-template/external/mysql"
	log "github.com/sirupsen/logrus"
)

type customerService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type CustomerServiceInterface interface {
	SCGetCekPengajuan(id_cust uint) interface{}
	VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error)
	SCCreatePengajuan(pengajuan *contract.Pengajuan, id_cust uint) *contract.Pengajuan
	SCCreateKelengkapan(kelengkapan *contract.Kelengkapan, id_cust uint) *contract.Kelengkapan
	SCGetByIdKelengkapan(id uint) interface{}
}

func NewCustomerService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *customerService {
	return &customerService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *customerService) SCGetCekPengajuan(id_cust uint) interface{} {
	var pengajuan contract.Pengajuan

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Table("pengajuans").First(&pengajuan, id_cust)

	return pengajuan
}

func (s *customerService) VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error) {
	claims := jwt.MapClaims{}

	err := s.jwtClient.ParseTokenWithClaims(req.Token, claims, s.appConfig.JWTSecret)

	if err != nil {
		log.Errorln(err)
		return nil, errors.NewUnauthorizedError("invalid parse token with claims")
	}

	authorized := fmt.Sprintf("%v", claims["authorized"])
	requestID := fmt.Sprintf("%v", claims["requestID"])

	if authorized == "" || requestID == "" {
		return nil, errors.NewUnauthorizedError("invalid payload")
	}

	ok, err := strconv.ParseBool(authorized)

	if err != nil || !ok {
		log.Errorln(err)
		return nil, errors.NewUnauthorizedError("invalid payload")
	}

	id_user_uint := uint(claims["id_user"].(float64))
	login_as_uint := uint(claims["login_as"].(float64))

	resp := &contract.JWTMapClaim{
		Authorized:     claims["authorized"].(bool),
		RequestID:      claims["requestID"].(string),
		IdUser:         id_user_uint,
		LoginAs:        login_as_uint,
		StandardClaims: jwt.StandardClaims{},
	}

	return resp, nil
}

func (s *customerService) SCCreatePengajuan(pengajuan *contract.Pengajuan, id_cust uint) *contract.Pengajuan {
	pengajuan.IdCust = id_cust
	status := 1
	pengajuan.Status = uint(status)

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Create(&pengajuan)

	return pengajuan
}

func (s *customerService) SCCreateKelengkapan(kelengkapan *contract.Kelengkapan, id uint) *contract.Kelengkapan {

	kelengkapan.IdPengajuan = id
	kelengkapan.StatusKelengkapan = 1

	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	db.DbConnection.Create(&kelengkapan)

	return kelengkapan
}

func (s *customerService) SCGetByIdKelengkapan(id uint) interface{} {

	var kelengkapan contract.Kelengkapan

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Table("kelengkapans").First(&kelengkapan, id)

	return kelengkapan
}
