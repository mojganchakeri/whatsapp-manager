package userConnections

import (
	"context"

	dto "github.com/mojganchakeri/whatsapp-manager/internal/DTO"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/config"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/database"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/logger"

	"gorm.io/gorm/clause"
)

type UserConnectionRepository interface {
	UpdateUserConnectionStatusByUserID(ctx context.Context, userID string, waUserID string, status string)
	UpdateUserConnectionStatusByWaUserID(ctx context.Context, waUserID string, status string)
	GetUserStatusByUserID(ctx context.Context, userID string) (string, error)
}

type userConnections struct {
	log logger.Logger
	db  database.Database
}

func New(
	log logger.Logger,
	cfg config.Config,
	db database.Database,
) UserConnectionRepository {
	return &userConnections{
		log: log,
		db:  db,
	}
}

func (ur *userConnections) UpdateUserConnectionStatusByUserID(ctx context.Context, userID string, waUserID string, status string) {
	var user dto.WhatsappConnection
	user.UserID = userID
	user.Status = status
	user.WhatsappUserID = waUserID

	err := ur.db.GetConnection(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"status", "whatsapp_user_id"}),
		}).
		Create(&user).
		Error
	if err != nil {
		ur.log.Error(err.Error())
	}
}

func (ur *userConnections) UpdateUserConnectionStatusByWaUserID(ctx context.Context, waUserID string, status string) {
	err := ur.db.GetConnection(ctx).
		Model(dto.WhatsappConnection{}).
		Where("whatsapp_user_id = ?", waUserID).
		Update("status", status).
		Error
	if err != nil {
		ur.log.Error(err.Error())
	}
}

func (ur *userConnections) GetUserStatusByUserID(ctx context.Context, userID string) (string, error) {
	var result dto.WhatsappConnection
	err := ur.db.GetConnection(ctx).
		Model(dto.WhatsappConnection{UserID: userID}).
		First(&result).
		Error
	if err != nil {
		ur.log.Error(err.Error())
		return "", err
	}
	return result.Status, nil
}
