package service

import (
	"fmt"
	"math"
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
	SPGetListPengajuan(page string) (*[]contract.ListPengajuan, error)
	SPGetListByName(name string) *[]contract.ListPengajuan
	SPGetCountPage() *contract.PengajuanPage
	VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error)
	SPGetSubmission(id uint) (*contract.Kelengkapan, error)
	SPPostSubmissionStatus(statusKelengkapan *contract.Kelengkapan, id uint) (*string, error)
	SPPostIdentityStatus(statusPengajuan *contract.Pengajuan, id uint) (*string, error)
}

func NewPetugasService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *petugasService {
	return &petugasService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *petugasService) SPGetListPengajuan(page string) (*[]contract.ListPengajuan, error) {
	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	pages, err := strconv.Atoi(page)

	if err != nil {
		return nil, err
	}

	var listPengajuan []contract.ListPengajuan

	for i := ((5 * pages) - 2); i < ((5 * pages) + 3); i++ {
		var pengajuan contract.Pengajuan
		var kelengkapan contract.Kelengkapan

		err := db.DbConnection.Table("pengajuans").Where("id_cust = ?", i).Find(&pengajuan).Error
		if pengajuan.IdCust == 0 {
			break
		}
		if err != nil {
			break
		}

		er := db.DbConnection.Table("kelengkapans").Last(&kelengkapan, "id_cust = ?", pengajuan.IdCust).Error
		if er != nil {
			lpengajuan := contract.ListPengajuan{
				TanggalPengajuan: pengajuan.UpdatedAt,
				NamaLengkap:      pengajuan.NamaLengkap,
				Status:           pengajuan.Status,
				Rekomendasi:      "-",
			}
			listPengajuan = append(listPengajuan, lpengajuan)
		} else {
			lpengajuan := contract.ListPengajuan{
				TanggalPengajuan: kelengkapan.UpdatedAt,
				NamaLengkap:      pengajuan.NamaLengkap,
				Status:           pengajuan.Status,
				Rekomendasi:      Recommendation(&pengajuan, &kelengkapan),
			}
			listPengajuan = append(listPengajuan, lpengajuan)
		}
	}
	return &listPengajuan, nil
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

func (s *petugasService) SPGetCountPage() *contract.PengajuanPage {
	var countPage int64

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Table("pengajuans").Count(&countPage)

	count := contract.PengajuanPage{
		CountPage: int64(math.Ceil(float64(countPage) / 5.00)),
	}
	return &count
}

func (s *petugasService) SPGetListByName(name string) *[]contract.ListPengajuan {
	var pengajuan []contract.Pengajuan
	var listPengajuan []contract.ListPengajuan
	var kelengkapan contract.Kelengkapan

	namePersen := fmt.Sprint("%" + name + "%")
	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	db.DbConnection.Table("pengajuans").Where("nama_lengkap LIKE ?", namePersen).Find(&pengajuan)

	for _, v := range pengajuan {
		er := db.DbConnection.Table("kelengkapans").Last(&kelengkapan, "id_cust = ?", v.IdCust).Error
		if er != nil {
			lpengajuan := contract.ListPengajuan{
				TanggalPengajuan: v.UpdatedAt,
				NamaLengkap:      v.NamaLengkap,
				Status:           v.Status,
				Rekomendasi:      "-",
			}
			listPengajuan = append(listPengajuan, lpengajuan)
		} else {
			lpengajuan := contract.ListPengajuan{
				TanggalPengajuan: kelengkapan.UpdatedAt,
				NamaLengkap:      v.NamaLengkap,
				Status:           v.Status,
				Rekomendasi:      Recommendation(&v, &kelengkapan),
			}
			listPengajuan = append(listPengajuan, lpengajuan)
		}
	}
	return &listPengajuan
}

func (s *petugasService) SPGetSubmission(id uint) (*contract.Kelengkapan, error) {

	var getkelengkapan contract.Kelengkapan

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("kelengkapans").Last(&getkelengkapan, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &getkelengkapan, nil
}

func (s *petugasService) SPPostSubmissionStatus(statusKelengkapan *contract.Kelengkapan, id uint) (*string, error) {
	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	var kelengkapanUpdates contract.Kelengkapan
	kelengkapanUpdates.StatusKelengkapan = statusKelengkapan.StatusKelengkapan
	var kelengkapan contract.Kelengkapan
	err := db.DbConnection.Table("kelengkapans").Last(&kelengkapan, "id_cust = ?", id).Error
	if err != nil {
		return nil, err
	}
	err = db.DbConnection.Model(&kelengkapan).Updates(kelengkapanUpdates).Error
	if err != nil {
		return nil, err
	}
	return &kelengkapanUpdates.StatusKelengkapan, nil
}

func (s *petugasService) SPPostIdentityStatus(statusPengajuan *contract.Pengajuan, id uint) (*string, error) {
	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	var pengajuanUpdates contract.Pengajuan          //rbody
	pengajuanUpdates.Status = statusPengajuan.Status //bridge rbody-db
	var pengajuan contract.Pengajuan                 //db
	err := db.DbConnection.Table("pengajuans").Last(&pengajuan, "id_cust = ?", id).Error
	if err != nil {
		return nil, err
	}
	err = db.DbConnection.Model(&pengajuan).Updates(pengajuanUpdates).Error
	if err != nil {
		return nil, err
	}
	return &pengajuanUpdates.Status, nil
}
