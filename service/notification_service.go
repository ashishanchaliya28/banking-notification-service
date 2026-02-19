package service

import (
	"context"
	"log"
	"github.com/banking-superapp/notification-service/model"
	"github.com/banking-superapp/notification-service/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationService interface {
	Send(ctx context.Context, req *model.SendNotificationRequest) (*model.Notification, error)
	GetNotifications(ctx context.Context, userID string, page, limit int64) ([]model.Notification, int64, error)
	RegisterDevice(ctx context.Context, userID string, req *model.RegisterDeviceRequest) error
	UpdatePreferences(ctx context.Context, userID string, req *model.UpdatePreferencesRequest) error
}

type notificationService struct {
	notifRepo  repository.NotificationRepo
	deviceRepo repository.DeviceTokenRepo
	prefRepo   repository.PreferencesRepo
}

func NewNotificationService(nr repository.NotificationRepo, dr repository.DeviceTokenRepo, pr repository.PreferencesRepo) NotificationService {
	return &notificationService{nr, dr, pr}
}

func (s *notificationService) Send(ctx context.Context, req *model.SendNotificationRequest) (*model.Notification, error) {
	oid, err := bson.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, err
	}

	notif := &model.Notification{
		UserID:  oid,
		Title:   req.Title,
		Body:    req.Body,
		Type:    req.Type,
		Channel: req.Channel,
		IsRead:  false,
		Data:    req.Data,
	}

	if err := s.notifRepo.Create(ctx, notif); err != nil {
		return nil, err
	}

	// In production: send FCM push, SMS via MSG91, email via SendGrid
	log.Printf("[NOTIFICATION] To: %s | Title: %s | Body: %s", req.UserID, req.Title, req.Body)

	return notif, nil
}

func (s *notificationService) GetNotifications(ctx context.Context, userID string, page, limit int64) ([]model.Notification, int64, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.notifRepo.FindByUserID(ctx, oid, page, limit)
}

func (s *notificationService) RegisterDevice(ctx context.Context, userID string, req *model.RegisterDeviceRequest) error {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	device := &model.DeviceToken{
		UserID:   oid,
		Token:    req.Token,
		Platform: req.Platform,
	}
	return s.deviceRepo.Upsert(ctx, device)
}

func (s *notificationService) UpdatePreferences(ctx context.Context, userID string, req *model.UpdatePreferencesRequest) error {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	existing, err := s.prefRepo.FindByUserID(ctx, oid)
	if err != nil && !isNotFound(err) {
		return err
	}

	pref := &model.Preferences{
		UserID:            oid,
		PushEnabled:       req.PushEnabled,
		SMSEnabled:        req.SMSEnabled,
		EmailEnabled:      req.EmailEnabled,
		TransactionAlerts: req.TransactionAlerts,
		PromoAlerts:       req.PromoAlerts,
	}
	if existing != nil {
		pref.ID = existing.ID
	}

	return s.prefRepo.Upsert(ctx, pref)
}

func isNotFound(err error) bool {
	return err == mongo.ErrNoDocuments
}
