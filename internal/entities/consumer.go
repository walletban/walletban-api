package entities

import "gorm.io/gorm"

type Consumer struct {
	gorm.Model
	ProjectID           uint
	Name                string
	Email               string
	IsWalletActivated   bool
	WalletGKey          string
	WalletEncryptedSKey string
}
