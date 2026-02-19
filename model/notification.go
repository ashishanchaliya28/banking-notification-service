package model

import (
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Notification struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id" json:"user_id"`
	Title     string        `bson:"title" json:"title"`
	Body      string        `bson:"body" json:"body"`
	Type      string        `bson:"type" json:"type"` // transaction | alert | promo | system
	Channel   string        `bson:"channel" json:"channel"` // push | sms | email | all
	IsRead    bool          `bson:"is_read" json:"is_read"`
	Data      interface{}   `bson:"data,omitempty" json:"data,omitempty"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

type DeviceToken struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id" json:"user_id"`
	Token     string        `bson:"token" json:"token"`
	Platform  string        `bson:"platform" json:"platform"` // android | ios
	IsActive  bool          `bson:"is_active" json:"is_active"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type Preferences struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           bson.ObjectID `bson:"user_id" json:"user_id"`
	PushEnabled      bool          `bson:"push_enabled" json:"push_enabled"`
	SMSEnabled       bool          `bson:"sms_enabled" json:"sms_enabled"`
	EmailEnabled     bool          `bson:"email_enabled" json:"email_enabled"`
	TransactionAlerts bool         `bson:"transaction_alerts" json:"transaction_alerts"`
	PromoAlerts      bool          `bson:"promo_alerts" json:"promo_alerts"`
	UpdatedAt        time.Time     `bson:"updated_at" json:"updated_at"`
}

type SendNotificationRequest struct {
	UserID  string      `json:"user_id"`
	Title   string      `json:"title"`
	Body    string      `json:"body"`
	Type    string      `json:"type"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data,omitempty"`
}

type RegisterDeviceRequest struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

type UpdatePreferencesRequest struct {
	PushEnabled       bool `json:"push_enabled"`
	SMSEnabled        bool `json:"sms_enabled"`
	EmailEnabled      bool `json:"email_enabled"`
	TransactionAlerts bool `json:"transaction_alerts"`
	PromoAlerts       bool `json:"promo_alerts"`
}
