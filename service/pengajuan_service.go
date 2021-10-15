package service

import (
	"net/http"
	"strconv"

	"github.com/rysmaadit/go-template/config"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/external/mysql"
)

type pengajuanService struct{}

type PengajuanServiceInterface interface {
	GetListPengajuan(w http.ResponseWriter, r *http.Request) interface{}
	GetSebelumPengajuan(w http.ResponseWriter, r *http.Request) interface{}
}

func NewPengajuanService(appConfig *config.Config) *pengajuanService {
	return &pengajuanService{}
}

func (s *pengajuanService) GetListPengajuan(w http.ResponseWriter, r *http.Request) interface{} {
	sPengajuan := []contract.Pengajuan{}

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Find(&sPengajuan)

	var listPengajuan []contract.ListPengajuan

	for _, v := range sPengajuan {
		lpengajuan := contract.ListPengajuan{
			Tanggal_pengajuan: v.CreatedAt,
			Nama_lengkap:      v.Nama_lengkap,
			Status:            v.Status,
			Rekomendasi:       "belum ada logicnya",
		}
		listPengajuan = append(listPengajuan, lpengajuan)
	}

	return listPengajuan
}

func (s *pengajuanService) GetSebelumPengajuan(w http.ResponseWriter, r *http.Request) interface{} {
	sPengajuan := []contract.Pengajuan{}

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Find(&sPengajuan, "status = ?", 1)

	return strconv.Itoa(len(sPengajuan))
}
