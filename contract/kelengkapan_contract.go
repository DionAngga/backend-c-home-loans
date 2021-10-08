package contract

import "gorm.io/gorm"

type Kelengkapan struct {
	gorm.Model
	Id_kelengkapan     uint   `gorm:"not null" json:"id_kelengkapan"`
	Alamat_rumah       string `gorm:"not null" json:"alamat_rumah"`
	Luas_tanah         uint   `gorm:"not null" json:"luas_tanah"`
	Harga_rumah        uint   `gorm:"not null" json:"harga_rumah"`
	Jangka_pembayaran  uint   `gorm:"not null" json:"jangka_pembayaran"`
	Dokumen_pendukung  string `gorm:"not null" json:"dokumen_pendukung"`
	Status_kelengkapan uint   `gorm:"not null" json:"status_completion"`
}
