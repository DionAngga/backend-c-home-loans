package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
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

type customerService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type CustomerServiceInterface interface {
	SCGetCheckApply(idCust uint) string
	VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error)
	SCCreateIdentity(identity *contract.Identity, idCust uint) (*contract.IdentityReturn, error)
	SCCreateSubmission(submission *contract.Submission, idCust uint) *contract.SubmissionReturn
	SCGetSubmissionStatus(id uint) string
	SCGetSubmission(id uint) (*contract.Submission, error)
	SCUploadFileKTP(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) string
	SCUploadFileGaji(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) string
	SCUploadFilePendukung(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) string
}

func NewCustomerService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *customerService {
	return &customerService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *customerService) SCGetCheckApply(idCust uint) string {
	var identity contract.Identity

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("identities").First(&identity, "id_cust = ?", idCust).Error

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

func (s *customerService) SCCreateIdentity(identity *contract.Identity, idCust uint) (*contract.IdentityReturn, error) {
	identity.IdCust = idCust
	identity.Status = "Menunggu Verifikasi"

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Create(&identity).Error
	if err != nil {
		log.Error("error connect db, %q", err)
	}

	pReturn := contract.IdentityReturn{
		IdCust:             identity.IdCust,
		Nik:                identity.Nik,
		NamaLengkap:        identity.NamaLengkap,
		TempatLahir:        identity.TempatLahir,
		TanggalLahir:       identity.TanggalLahir,
		Alamat:             identity.Alamat,
		Pekerjaan:          identity.Pekerjaan,
		PendapatanPerbulan: identity.PendapatanPerbulan,
		BuktiKtp:           identity.BuktiKtp,
		Status:             identity.Status,
	}
	return &pReturn, nil
}

func (s *customerService) SCCreateSubmission(submission *contract.Submission, id uint) *contract.SubmissionReturn {

	submission.IdCust = id
	submission.IdPengajuan = id
	submission.StatusKelengkapan = "Menunggu Persetujuan"

	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	db.DbConnection.Create(&submission)

	kReturn := contract.SubmissionReturn{
		IdCust:            submission.IdCust,
		IdPengajuan:       submission.IdPengajuan,
		AlamatRumah:       submission.AlamatRumah,
		LuasTanah:         submission.LuasTanah,
		HargaRumah:        submission.HargaRumah,
		JangkaPembayaran:  submission.JangkaPembayaran,
		DokumenPendukung:  submission.DokumenPendukung,
		StatusKelengkapan: submission.StatusKelengkapan,
	}
	return &kReturn

}

func (s *customerService) SCGetSubmission(id uint) (*contract.Submission, error) {
	var getSubmission contract.Submission
	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	err := db.DbConnection.Table("submissions").Last(&getSubmission, "id_pengajuan = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &getSubmission, nil
}

func (s *customerService) SCGetSubmissionStatus(id uint) string {
	var getStatusKelengkapan contract.Submission
	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	err := db.DbConnection.Table("submissions").Last(&getStatusKelengkapan, "id_pengajuan = ?", id).Error
	if err != nil {
		return "Menu Submission invisible(Menu disable)"
	}
	return "Menu Submission visible(Menu able)"
}

func (s *customerService) SCUploadFile(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) string {
	idString := strconv.Itoa(int(resp.IdUser))
	fileLink := strings.Join([]string{"ktp-", idString, "-", resp.Username, ".pdf"}, "")
	fileName := strings.Join([]string{"ktp/", fileLink}, "")

	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()

	fileReader := io.Reader(*file)
	uploadInfo, err := mi.MinioClient.PutObject(ctx, mi.BucketName, fileName, fileReader, handler.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("Error in uploading the file #%s: %v.", fileName, err)
		return "Error in uploading the file"
	}

	log.Printf("Uploading the file #%s succeeded!", fileName)
	fmt.Println("UploadInfo:")
	fmt.Printf("%+v\n", uploadInfo)

	return fileLink
}

func (s *customerService) SCUploadFilePendukung(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) string {
	idString := strconv.Itoa(int(resp.IdUser))
	fileLink := strings.Join([]string{"bukti-pendukung-", idString, "-", resp.Username, ".pdf"}, "")
	fileName := strings.Join([]string{"bukti-pendukung/", fileLink}, "")

	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()

	fileReader := io.Reader(*file)
	uploadInfo, err := mi.MinioClient.PutObject(ctx, mi.BucketName, fileName, fileReader, handler.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("Error in uploading the file #%s: %v.", fileName, err)
		return "Error in uploading the file"
	}

	log.Printf("Uploading the file #%s succeeded!", fileName)
	fmt.Println("UploadInfo:")
	fmt.Printf("%+v\n", uploadInfo)

	return fileLink
}

func (s *customerService) SCUploadFileKTP(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) string {
	idString := strconv.Itoa(int(resp.IdUser))
	fileLink := strings.Join([]string{"ktp-", idString, "-", resp.Username, ".pdf"}, "")
	fileName := strings.Join([]string{"ktp/", fileLink}, "")

	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()

	fileReader := io.Reader(*file)
	uploadInfo, err := mi.MinioClient.PutObject(ctx, mi.BucketName, fileName, fileReader, handler.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("Error in uploading the file #%s: %v.", fileName, err)
		return "Error in uploading the file"
	}

	log.Printf("Uploading the file #%s succeeded!", fileName)
	fmt.Println("UploadInfo:")
	fmt.Printf("%+v\n", uploadInfo)

	return fileLink
}

func (s *customerService) SCUploadFileGaji(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) string {
	idString := strconv.Itoa(int(resp.IdUser))
	fileLink := strings.Join([]string{"slip-gaji-", idString, "-", resp.Username, ".pdf"}, "")
	fileName := strings.Join([]string{"slip-gaji/", fileLink}, "")

	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()

	fileReader := io.Reader(*file)
	uploadInfo, err := mi.MinioClient.PutObject(ctx, mi.BucketName, fileName, fileReader, handler.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("Error in uploading the file #%s: %v.", fileName, err)
		return "Error in uploading the file"
	}

	log.Printf("Uploading the file #%s succeeded!", fileName)
	fmt.Println("UploadInfo:")
	fmt.Printf("%+v\n", uploadInfo)

	return fileLink
}

// func (s *customerService) SCUploadFilePendukung(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) string {
// 	idString := strconv.Itoa(int(resp.IdUser))
// 	fileLink := strings.Join([]string{"bukti-pendukung-", idString, "-", resp.Username, ".pdf"}, "")
// 	fileName := strings.Join([]string{"bukti-endukung/", fileLink}, "")

// 	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

// 	ctx := context.Background()

// 	fileReader := io.Reader(*file)
// 	uploadInfo, err := mi.MinioClient.PutObject(ctx, mi.BucketName, fileName, fileReader, handler.Size, minio.PutObjectOptions{})
// 	if err != nil {
// 		log.Printf("Error in uploading the file #%s: %v.", fileName, err)
// 		return "Error in uploading the file"
// 	}

// 	log.Printf("Uploading the file #%s succeeded!", fileName)
// 	fmt.Println("UploadInfo:")
// 	fmt.Printf("%+v\n", uploadInfo)

// 	return fileLink
// }
