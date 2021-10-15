package contract

import (
	"time"

	"gorm.io/gorm"
)

type Pengajuan struct {
	gorm.Model
	IdCust             uint    `gorm:"not null;unique" json:"id_cust"`
	Nik                string  `gorm:"not null;unique" json:"nik" validate:"required"`
	NamaLengkap        string  `gorm:"not null" json:"nama_lengkap" validate:"required"`
	TempatLahir        string  `gorm:"not null" json:"tempat_lahir" validate:"required"`
	TanggalLahir       string  `gorm:"not null" json:"tanggal_lahir" validate:"required"`
	Alamat             string  `gorm:"not null" json:"alamat" validate:"required"`
	Pekerjaan          string  `gorm:"not null" json:"pekerjaan" validate:"required"`
	PendapatanPerbulan float64 `gorm:"not null" json:"pendapatan_perbulan" validate:"required"`
	BuktiKtp           string  `gorm:"not null" json:"bukti_ktp" validate:"required"`
	Status             string  `gorm:"not null" json:"status"`
}

type ListPengajuan struct {
	TanggalPengajuan time.Time `json:"tanggal_pengajuan"`
	NamaLengkap      string    `json:"nama_lengkap"`
	Status           string    `json:"status"`
	Rekomendasi      string    `json:"rekomendasi"`
}

type PengajuanReturn struct {
	IdCust             uint    `gorm:"not null;unique" json:"id_cust"`
	Nik                string  `gorm:"not null;unique" json:"nik"`
	NamaLengkap        string  `gorm:"not null" json:"nama_lengkap"`
	TempatLahir        string  `gorm:"not null" json:"tempat_lahir"`
	TanggalLahir       string  `gorm:"not null" json:"tanggal_lahir"`
	Alamat             string  `gorm:"not null" json:"alamat"`
	Pekerjaan          string  `gorm:"not null" json:"pekerjaan"`
	PendapatanPerbulan float64 `gorm:"not null" json:"pendapatan_perbulan"`
	BuktiKtp           string  `gorm:"not null" json:"bukti_ktp"`
	Status             string  `gorm:"not null" json:"status"`
}

type PengajuanPage struct {
	CountPage int64 `json:"count_page"`
}
