package entities

import "gorm.io/gorm"

type Consumer struct {
	gorm.Model
	ProjectID           uint
	Name                string
	Email               string
	IsFirstTime         bool `gorm:"default:true"`
	IsWalletActivated   bool `gorm:"default:false"`
	PfpUrl              string
	WalletGKey          string
	WalletEncryptedSKey string
	Token               string
}
