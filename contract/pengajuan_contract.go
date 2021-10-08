package contract

import (
	"time"

	"gorm.io/gorm"
)

type Pengajuan struct {
	gorm.Model
	Id_cust             uint    `gorm:"not null;unique" json:"id_cust"`
	Nik                 string  `gorm:"not null;unique" json:"nik"`
	Nama_lengkap        string  `gorm:"not null" json:"nama_lengkap"`
	Tempat_lahir        string  `gorm:"not null" json:"tempat_lahir"`
	Tanggal_lahir       string  `gorm:"not null" json:"tanggal_lahir"`
	Alamat              string  `gorm:"not null" json:"alamat"`
	Pekerjaan           string  `gorm:"not null" json:"pekerjaan"`
	Pendapatan_perbulan float64 `gorm:"not null" json:"pendapatan_perbulan"`
	Bukti_ktp           string  `gorm:"not null" json:"bukti_ktp"`
	Status              uint    `gorm:"not null" json:"status"`
}

type ListPengajuan struct {
	Tanggal_pengajuan time.Time `json:"tanggal_pengajuan"`
	Nama_lengkap      string    `json:"nama_lengkap"`
	Status            uint      `json:"status"`
	Rekomendasi       string    `json:"rekomendasi"`
}
