package user_session

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	SessionToken string `json:"session_token" bson:"session_token"`
	Handle string `json:"handle" bson:"handle"`
	DeviceId string `json:"device_id" bson:"device_id"`
	BrowserType string `json:"browser_type" bson:"browser_type"`
	ActiveTill int64 `json:"active_till" bson:"active_till"`
	CreatedAt int64 `json:"started_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at, omitempty"`
}

type Device struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	DeviceName string `json:"device_name" bson:"device_name"`
	DeviceIP string `json:"device_ip" bson:"device_id"`
	CreatedAt int64 `json:"started_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at, omitempty"`
}