package contract

import "gorm.io/gorm"

type Kelengkapan struct {
	gorm.Model
	IdCust            uint    `gorm:"not null" json:"id_cust"`
	IdPengajuan       uint    `gorm:"not null" json:"id_pengajuan"`
	AlamatRumah       string  `gorm:"not null" json:"alamat_rumah"`
	LuasTanah         float64 `gorm:"not null" json:"luas_tanah"`
	HargaRumah        float64 `gorm:"not null" json:"harga_rumah"`
	JangkaPembayaran  uint    `gorm:"not null" json:"jangka_pembayaran"`
	DokumenPendukung  string  `gorm:"not null" json:"dokumen_pendukung"`
	StatusKelengkapan uint    `gorm:"not null" json:"status_kelengkapan"`
}
