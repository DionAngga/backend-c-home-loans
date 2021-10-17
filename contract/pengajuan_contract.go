package contract

import (
	"time"

	"gorm.io/gorm"
)

type Identity struct {
	gorm.Model
	IdCust             uint    `gorm:"not null;unique" json:"id_cust"`
	Nik                string  `gorm:"not null;unique" json:"nik" validate:"required,len=16,numeric"`
	NamaLengkap        string  `gorm:"not null" json:"nama_lengkap" validate:"required"`
	TempatLahir        string  `gorm:"not null" json:"tempat_lahir" validate:"required"`
	TanggalLahir       string  `gorm:"not null" json:"tanggal_lahir" validate:"required"`
	Alamat             string  `gorm:"not null" json:"alamat" validate:"required"`
	Pekerjaan          string  `gorm:"not null" json:"pekerjaan" validate:"required"`
	PendapatanPerbulan float64 `gorm:"not null" json:"pendapatan_perbulan" validate:"required"`
	BuktiKtp           string  `gorm:"not null" json:"bukti_ktp" validate:"required"`
	BuktiGaji          string  `gorm:"not null" json:"bukti_gaji" validate:"required"`
	Status             string  `gorm:"not null" json:"status"`
}

type ListSubmission struct {
	TanggalPengajuan time.Time `json:"tanggal_pengajuan"`
	NamaLengkap      string    `json:"nama_lengkap"`
	Status           string    `json:"status"`
	Rekomendasi      string    `json:"rekomendasi"`
}

type IdentityReturn struct {
	IdCust             uint    `gorm:"not null;unique" json:"id_cust"`
	Nik                string  `gorm:"not null;unique" json:"nik"`
	NamaLengkap        string  `gorm:"not null" json:"nama_lengkap"`
	TempatLahir        string  `gorm:"not null" json:"tempat_lahir"`
	TanggalLahir       string  `gorm:"not null" json:"tanggal_lahir"`
	Alamat             string  `gorm:"not null" json:"alamat"`
	Pekerjaan          string  `gorm:"not null" json:"pekerjaan"`
	PendapatanPerbulan float64 `gorm:"not null" json:"pendapatan_perbulan"`
	BuktiKtp           string  `gorm:"not null" json:"bukti_ktp"`
	BuktiGaji          string  `gorm:"not null" json:"bukti_gaji"`
	Status             string  `gorm:"not null" json:"status"`
}

type NumberOfPage struct {
	NumberOfPage int64 `json:"number_of_page"`
}
