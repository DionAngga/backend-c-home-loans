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
	SCGetCekPengajuan(idCust uint) string
	VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error)
	SCCreatePengajuan(pengajuan *contract.Pengajuan, idCust uint) (*contract.PengajuanReturn, error)
	SCCreateKelengkapan(kelengkapan *contract.Kelengkapan, id uint) *contract.KelengkapanReturn
	SCGetByIdKelengkapan(id uint) (*contract.Kelengkapan, error)
	SCGetStatusByIdKelengkapan(id uint) string
}

func NewCustomerService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *customerService {
	return &customerService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *customerService) SCGetCekPengajuan(idCust uint) string {
	var pengajuan contract.Pengajuan

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("pengajuans").First(&pengajuan, "id_cust = ?", idCust).Error

	if err != nil {
		return "Anda sedang tidak mengajukan KPR saat ini"
	}

	return "Anda sedang mengajukan KPR saat ini"
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
		Username:       claims["username"].(string),
		LoginAs:        login_as_uint,
		StandardClaims: jwt.StandardClaims{},
	}

	return resp, nil
}

func (s *customerService) SCCreatePengajuan(pengajuan *contract.Pengajuan, idCust uint) (*contract.PengajuanReturn, error) {
	pengajuan.IdCust = idCust
	// status := 1
	pengajuan.Status = "Waiting"

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Create(&pengajuan).Error
	if err != nil {
		log.Error("error connect db, %q", err)
	}

	pReturn := contract.PengajuanReturn{
		IdCust:             pengajuan.IdCust,
		Nik:                pengajuan.Nik,
		NamaLengkap:        pengajuan.NamaLengkap,
		TempatLahir:        pengajuan.TempatLahir,
		TanggalLahir:       pengajuan.TanggalLahir,
		Alamat:             pengajuan.Alamat,
		Pekerjaan:          pengajuan.Pekerjaan,
		PendapatanPerbulan: pengajuan.PendapatanPerbulan,
		BuktiKtp:           pengajuan.BuktiKtp,
		Status:             pengajuan.Status,
	}
	return &pReturn, nil
}

func (s *customerService) SCCreateKelengkapan(kelengkapan *contract.Kelengkapan, id uint) *contract.KelengkapanReturn {

	kelengkapan.IdCust = id
	kelengkapan.IdPengajuan = id
	kelengkapan.StatusKelengkapan = "Waiting"

	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	db.DbConnection.Create(&kelengkapan)

	kReturn := contract.KelengkapanReturn{
		IdCust:            kelengkapan.IdCust,
		IdPengajuan:       kelengkapan.IdPengajuan,
		AlamatRumah:       kelengkapan.AlamatRumah,
		LuasTanah:         kelengkapan.LuasTanah,
		HargaRumah:        kelengkapan.HargaRumah,
		JangkaPembayaran:  kelengkapan.JangkaPembayaran,
		DokumenPendukung:  kelengkapan.DokumenPendukung,
		StatusKelengkapan: kelengkapan.StatusKelengkapan,
	}
	return &kReturn

}

func (s *customerService) SCGetByIdKelengkapan(id uint) (*contract.Kelengkapan, error) {

	var getkelengkapan contract.Kelengkapan

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("kelengkapans").Last(&getkelengkapan, "id_pengajuan = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &getkelengkapan, nil
}

func (s *customerService) SCGetStatusByIdKelengkapan(id uint) string {

	var getStatusKelengkapan contract.Kelengkapan //db

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("kelengkapans").Last(&getStatusKelengkapan, "id_pengajuan = ?", id).Error

	if err != nil {

		return "Menu Submission invisible(Menu disable)"
	}

	return "Menu Submission visible(Menu able)"
}
