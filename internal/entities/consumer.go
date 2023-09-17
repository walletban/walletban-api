package entities

import "gorm.io/gorm"

type Consumer struct {
	gorm.Model
	ProjectID           uint
	Name                string
	Email               string
	IsFirstTime         bool
	IsWalletActivated   bool
	WalletGKey          string
	WalletEncryptedSKey string
}
