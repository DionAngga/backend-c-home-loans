package contract

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	LoginAs  uint   `gorm:"not null" json:"login_as"`
}
