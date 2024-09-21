package whatsapp

import (
	"context"
	"fmt"

	"github.com/mojganchakeri/whatsapp-manager/internal/consts"
	"github.com/mojganchakeri/whatsapp-manager/internal/init/logger"
	"github.com/mojganchakeri/whatsapp-manager/internal/repository/postgres/userConnections"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
)

type WhatsappService interface {
	CreateClient(ctx context.Context, userID string) (*whatsmeow.Client, error)
	Login(ctx context.Context, client *whatsmeow.Client, userID string) (string, error)
}

type whatsappService struct {
	log            logger.Logger
	userRepository userConnections.UserConnectionRepository
}

func New(
	log logger.Logger,
	userRepository userConnections.UserConnectionRepository,
) WhatsappService {
	return &whatsappService{
		log:            log,
		userRepository: userRepository,
	}
}

// Function to handle WhatsApp events
func handleEvent(
	ctx context.Context,
	log logger.Logger,
	userRepository userConnections.UserConnectionRepository,
	client *whatsmeow.Client,
) func(evt interface{}) {

	return func(evt interface{}) {
		switch evt.(type) {
		case *events.Disconnected:
			log.Debug(fmt.Sprintf("user %s disconnected", client.Store.ID.String()))
			userRepository.UpdateUserConnectionStatusByWaUserID(ctx, client.Store.ID.String(), consts.UserStatusDisconnect)

		case *events.LoggedOut:
			log.Debug(fmt.Sprintf("user %s logged out", client.Store.ID.String()))
			userRepository.UpdateUserConnectionStatusByWaUserID(ctx, client.Store.ID.String(), consts.UserStatusLogout)
		}
	}
}

func (ws *whatsappService) CreateClient(ctx context.Context, userID string) (*whatsmeow.Client, error) {
	dbPath := fmt.Sprintf("file:whatsapp_%s.db?_foreign_keys=on", userID)

	container, err := sqlstore.New("sqlite3", dbPath, nil)
	if err != nil {
		err := fmt.Errorf("failed to create database store: %v", err)
		ws.log.Error(err.Error())
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		err := fmt.Errorf("failed to get device for user %s: %w", userID, err)
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, nil)
	client.AddEventHandler(handleEvent(ctx, ws.log, ws.userRepository, client))
	return client, nil
}

func (ws *whatsappService) Login(ctx context.Context, client *whatsmeow.Client, userID string) (string, error) {
	if client.Store.ID == nil {
		// If the user is not logged in, generate QR code for login
		qrChan, _ := client.GetQRChannel(ctx)
		err := client.Connect()
		if err != nil {
			err = fmt.Errorf("failed to connect to WhatsApp: %v", err)
			ws.log.Error(err.Error())
			return "", err
		}

		ch := make(chan string)

		go func() {
			for evt := range qrChan {

				if evt.Event == "code" {
					ws.log.Debug(fmt.Sprintf("user %s waiting for login", userID))
					ws.userRepository.UpdateUserConnectionStatusByUserID(ctx, userID, "", consts.UserStatusWaitForLogin)
					ch <- evt.Code

				} else if evt.Event == "success" {
					ws.log.Debug(fmt.Sprintf("user %s connected - %s", userID, client.Store.ID.String()))
					ws.userRepository.UpdateUserConnectionStatusByUserID(ctx, userID, client.Store.ID.String(), consts.UserStatusLogin)
				}
			}
		}()

		val, ok := <-ch
		if ok {
			return val, nil
		}

	} else {
		// User is already logged in
		err := client.Connect()
		if err != nil {
			err = fmt.Errorf("failed to connect to WhatsApp: %v", err)
			ws.log.Error(err.Error())
			return "", err
		}
		fmt.Println("Logged in successfully!")
		return "", nil
	}

	return "", fmt.Errorf("unhandled error")
}
