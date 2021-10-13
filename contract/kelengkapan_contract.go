package contract

import (
	"gorm.io/gorm"
)

type Kelengkapan struct {
	gorm.Model
	IdCust            uint    `gorm:"not null" json:"id_cust"`
	IdPengajuan       uint    `gorm:"not null" json:"id_pengajuan"`
	AlamatRumah       string  `gorm:"not null" json:"alamat_rumah" validate:"required"`
	LuasTanah         float64 `gorm:"not null" json:"luas_tanah" validate:"required"`
	HargaRumah        float64 `gorm:"not null" json:"harga_rumah" validate:"required"`
	JangkaPembayaran  uint    `gorm:"not null" json:"jangka_pembayaran" validate:"required"`
	DokumenPendukung  string  `gorm:"not null" json:"dokumen_pendukung" validate:"required"`
	StatusKelengkapan string  `gorm:"not null" json:"status_kelengkapan"`
}

type KelengkapanReturn struct {
	IdCust            uint    `gorm:"not null" json:"id_cust"`
	IdPengajuan       uint    `gorm:"not null" json:"id_pengajuan"`
	AlamatRumah       string  `gorm:"not null" json:"alamat_rumah" validate:"required"`
	LuasTanah         float64 `gorm:"not null" json:"luas_tanah" validate:"required"`
	HargaRumah        float64 `gorm:"not null" json:"harga_rumah" validate:"required"`
	JangkaPembayaran  uint    `gorm:"not null" json:"jangka_pembayaran" validate:"required"`
	DokumenPendukung  string  `gorm:"not null" json:"dokumen_pendukung" validate:"required"`
	StatusKelengkapan string  `gorm:"not null" json:"status_kelengkapan"`
}
