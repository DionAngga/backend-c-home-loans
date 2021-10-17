package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/minio/minio-go/v7"
	"github.com/rysmaadit/go-template/common/errors"
	"github.com/rysmaadit/go-template/config"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/external/jwt_client"
	miniopkg "github.com/rysmaadit/go-template/external/minio"
	"github.com/rysmaadit/go-template/external/mysql"
	log "github.com/sirupsen/logrus"
)

type employeeService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type EmployeeServiceInterface interface {
	SPGetListSubmission(page string) (*[]contract.ListSubmission, error)
	SPGetListByName(name string) *[]contract.ListSubmission
	SPGetNumberOfPage() *contract.NumberOfPage
	VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error)
	SPGetSubmission(id uint) (*contract.Submission, error)
	SPPostSubmissionStatus(submissionStatus *contract.Submission, id uint) (*string, error)
	SPPostIdentityStatus(statusPengajuan *contract.Identity, id uint) (*string, error)
	SPGetFileKtp(buktiKtp string) *minio.Object
	SPGetFileBuktiGaji(buktiGaji string) *minio.Object
	SPGetFileBuktiPendukung(buktiPendukung string) *minio.Object
}

func NewEmployeeService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *employeeService {
	return &employeeService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *employeeService) SPGetListSubmission(page string) (*[]contract.ListSubmission, error) {
	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	pages, err := strconv.Atoi(page)

	if err != nil {
		return nil, err
	}

	var ListSubmission []contract.ListSubmission

	for i := ((5 * pages) - 2); i < ((5 * pages) + 3); i++ {
		var identity contract.Identity
		var submission contract.Submission

		err := db.DbConnection.Table("identities").Where("id_cust = ?", i).Find(&identity).Error
		if identity.IdCust == 0 {
			break
		}
		if err != nil {
			break
		}

		er := db.DbConnection.Table("submissions").Last(&submission, "id_cust = ?", identity.IdCust).Error
		if er != nil {
			lsubmission := contract.ListSubmission{
				TanggalPengajuan: identity.UpdatedAt,
				NamaLengkap:      identity.NamaLengkap,
				Status:           identity.Status,
				Rekomendasi:      "-",
			}
			ListSubmission = append(ListSubmission, lsubmission)
		} else {
			lsubmission := contract.ListSubmission{
				TanggalPengajuan: submission.UpdatedAt,
				NamaLengkap:      identity.NamaLengkap,
				Status:           identity.Status,
				Rekomendasi:      Recommendation(&identity, &submission),
			}
			ListSubmission = append(ListSubmission, lsubmission)
		}
	}
	return &ListSubmission, nil
}

func Recommendation(identity *contract.Identity, submission *contract.Submission) string {
	var kemampuanCicilanPerbulan float64
	var kenyataanCicilanPerbulan float64
	kemampuanCicilanPerbulan = (identity.PendapatanPerbulan / 3)
	kenyataanCicilanPerbulan = (submission.HargaRumah / float64(submission.JangkaPembayaran)) / 12
	if kemampuanCicilanPerbulan > kenyataanCicilanPerbulan {
		return "Boleh"
	}
	return "Tidak Boleh"
}

func (s *employeeService) VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error) {
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

func (s *employeeService) SPGetNumberOfPage() *contract.NumberOfPage {
	var NumberOfPage int64

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Table("identities").Count(&NumberOfPage)

	count := contract.NumberOfPage{
		NumberOfPage: int64(math.Ceil(float64(NumberOfPage) / 5.00)),
	}
	return &count
}

func (s *employeeService) SPGetListByName(name string) *[]contract.ListSubmission {
	var identity []contract.Identity
	var ListSubmission []contract.ListSubmission
	var submission contract.Submission

	namePersen := fmt.Sprint("%" + name + "%")
	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	db.DbConnection.Table("identities").Where("nama_lengkap LIKE ?", namePersen).Find(&identity)

	for _, v := range identity {
		er := db.DbConnection.Table("submissions").Last(&submission, "id_cust = ?", v.IdCust).Error
		if er != nil {
			lsubmission := contract.ListSubmission{
				TanggalPengajuan: v.UpdatedAt,
				NamaLengkap:      v.NamaLengkap,
				Status:           v.Status,
				Rekomendasi:      "-",
			}
			ListSubmission = append(ListSubmission, lsubmission)
		} else {
			lsubmission := contract.ListSubmission{
				TanggalPengajuan: submission.UpdatedAt,
				NamaLengkap:      v.NamaLengkap,
				Status:           v.Status,
				Rekomendasi:      Recommendation(&v, &submission),
			}
			ListSubmission = append(ListSubmission, lsubmission)
		}
	}
	return &ListSubmission
}

func (s *employeeService) SPGetSubmission(id uint) (*contract.Submission, error) {

	var getSubmission contract.Submission

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("submissions").Last(&getSubmission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &getSubmission, nil
}

func (s *employeeService) SPPostSubmissionStatus(submissionStatus *contract.Submission, id uint) (*string, error) {
	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	var submissionUpdates contract.Submission
	submissionUpdates.StatusKelengkapan = submissionStatus.StatusKelengkapan
	var submission contract.Submission
	err := db.DbConnection.Table("submissions").Last(&submission, "id_cust = ?", id).Error
	if err != nil {
		return nil, err
	}
	err = db.DbConnection.Model(&submission).Updates(submissionUpdates).Error
	if err != nil {
		return nil, err
	}
	return &submissionUpdates.StatusKelengkapan, nil
}

func (s *employeeService) SPPostIdentityStatus(statusPengajuan *contract.Identity, id uint) (*string, error) {
	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	var pengajuanUpdates contract.Identity
	pengajuanUpdates.Status = statusPengajuan.Status
	var pengajuan contract.Identity

  err := db.DbConnection.Table("identities").Last(&pengajuan, "id_cust = ?", id).Error

	if err != nil {
		return nil, err
	}
	err = db.DbConnection.Model(&pengajuan).Updates(pengajuanUpdates).Error
	if err != nil {
		return nil, err
	}
	return &pengajuanUpdates.Status, nil
}

func (s *employeeService) SPGetFileKtp(buktiKtp string) *minio.Object {
	fileName := strings.Join([]string{"ktp/", buktiKtp}, "")
	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()
	obj, err := mi.MinioClient.GetObject(ctx, mi.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error in getting the object: %v.", err)
		return nil
	}
	return obj
}

func (s *employeeService) SPGetFileBuktiGaji(buktiGaji string) *minio.Object {
	fileName := strings.Join([]string{"slip-gaji/", buktiGaji}, "")
	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()
	obj, err := mi.MinioClient.GetObject(ctx, mi.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error in getting the object: %v.", err)
		return nil
	}
	return obj
}

func (s *employeeService) SPGetFileBuktiPendukung(buktiPendukung string) *minio.Object {
	fileName := strings.Join([]string{"bukti-pendukung/", buktiPendukung}, "")
	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()
	obj, err := mi.MinioClient.GetObject(ctx, mi.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error in getting the object: %v.", err)
		return nil
	}
	return obj
}
