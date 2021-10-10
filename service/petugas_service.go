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

type petugasService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type PetugasServiceInterface interface {
	SPGetListPengajuan(page string) (interface{}, error)
	VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error)
}

func NewPetugasService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *petugasService {
	return &petugasService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *petugasService) SPGetListPengajuan(page string) (interface{}, error) {
	var pengajuan contract.Pengajuan
	var kelengkapan contract.Kelengkapan

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	pages, err := strconv.Atoi(page)

	if err != nil {
		return nil, err
	}

	var listPengajuan []contract.ListPengajuan

	for i := ((5 * pages) - 4); i < ((5 * pages) + 1); i++ {
		err := db.DbConnection.Table("pengajuans").First(&pengajuan, "ID = ?", i).Error
		if err != nil {
			break
		}

		er := db.DbConnection.Table("kelengkapans").Last(&kelengkapan, "id_cust = ?", pengajuan.IdCust).Error
		if er != nil {
			break
		}
		lpengajuan := contract.ListPengajuan{
			TanggalPengajuan: pengajuan.UpdatedAt,
			NamaLengkap:      pengajuan.NamaLengkap,
			Status:           pengajuan.Status,
			Rekomendasi:      Recommendation(&pengajuan, &kelengkapan),
		}
		listPengajuan = append(listPengajuan, lpengajuan)
	}
	return listPengajuan, nil
}

func Recommendation(pengajuan *contract.Pengajuan, kelengkapan *contract.Kelengkapan) string {
	var kemampuanCicilanPerbulan float64
	var kenyataanCicilanPerbulan float64
	kemampuanCicilanPerbulan = (pengajuan.PendapatanPerbulan / 3)
	kenyataanCicilanPerbulan = (kelengkapan.HargaRumah / float64(kelengkapan.JangkaPembayaran)) / 12
	if kemampuanCicilanPerbulan > kenyataanCicilanPerbulan {
		return "Boleh"
	}
	return "Tidak Boleh"
}

func (s *petugasService) VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error) {
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
		Username:       claims["username"].(string),
		LoginAs:        login_as_uint,
		StandardClaims: jwt.StandardClaims{},
	}

	return resp, nil
}
