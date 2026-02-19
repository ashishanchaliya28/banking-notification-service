package repository

import (
	"context"
	"time"
	"github.com/banking-superapp/notification-service/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type NotificationRepo interface {
	Create(ctx context.Context, n *model.Notification) error
	FindByUserID(ctx context.Context, userID bson.ObjectID, page, limit int64) ([]model.Notification, int64, error)
}

type DeviceTokenRepo interface {
	Upsert(ctx context.Context, d *model.DeviceToken) error
	FindByUserID(ctx context.Context, userID bson.ObjectID) ([]model.DeviceToken, error)
}

type PreferencesRepo interface {
	Upsert(ctx context.Context, p *model.Preferences) error
	FindByUserID(ctx context.Context, userID bson.ObjectID) (*model.Preferences, error)
}

type notificationRepo struct{ col *mongo.Collection }
type deviceTokenRepo struct{ col *mongo.Collection }
type preferencesRepo struct{ col *mongo.Collection }

func NewNotificationRepo(db *mongo.Database) NotificationRepo { return &notificationRepo{col: db.Collection("notifications")} }
func NewDeviceTokenRepo(db *mongo.Database) DeviceTokenRepo   { return &deviceTokenRepo{col: db.Collection("device_tokens")} }
func NewPreferencesRepo(db *mongo.Database) PreferencesRepo   { return &preferencesRepo{col: db.Collection("preferences")} }

func (r *notificationRepo) Create(ctx context.Context, n *model.Notification) error {
	n.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, n)
	return err
}

func (r *notificationRepo) FindByUserID(ctx context.Context, userID bson.ObjectID, page, limit int64) ([]model.Notification, int64, error) {
	filter := bson.M{"user_id": userID}
	total, _ := r.col.CountDocuments(ctx, filter)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip((page-1)*limit).SetLimit(limit)
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var notifications []model.Notification
	cursor.All(ctx, &notifications)
	return notifications, total, nil
}

func (r *deviceTokenRepo) Upsert(ctx context.Context, d *model.DeviceToken) error {
	d.UpdatedAt = time.Now()
	d.IsActive = true
	_, err := r.col.UpdateOne(ctx,
		bson.M{"user_id": d.UserID, "token": d.Token},
		bson.M{"$set": d},
		&mongo.UpdateOptions{Upsert: boolPtr(true)},
	)
	return err
}

func (r *deviceTokenRepo) FindByUserID(ctx context.Context, userID bson.ObjectID) ([]model.DeviceToken, error) {
	cursor, err := r.col.Find(ctx, bson.M{"user_id": userID, "is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tokens []model.DeviceToken
	cursor.All(ctx, &tokens)
	return tokens, nil
}

func (r *preferencesRepo) Upsert(ctx context.Context, p *model.Preferences) error {
	p.UpdatedAt = time.Now()
	_, err := r.col.UpdateOne(ctx,
		bson.M{"user_id": p.UserID},
		bson.M{"$set": p},
		&mongo.UpdateOptions{Upsert: boolPtr(true)},
	)
	return err
}

func (r *preferencesRepo) FindByUserID(ctx context.Context, userID bson.ObjectID) (*model.Preferences, error) {
	var p model.Preferences
	err := r.col.FindOne(ctx, bson.M{"user_id": userID}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func boolPtr(b bool) *bool { return &b }
