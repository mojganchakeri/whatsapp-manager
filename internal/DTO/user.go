package dto

type WhatsappConnection struct {
	UserID         string `gorm:"user_id"`
	Status         string `gorm:"status"`
	WhatsappUserID string `gorm:"whatsapp_user_id"`
}
