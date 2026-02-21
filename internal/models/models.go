package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceLogs struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Mqtt      string             `json:"mqtt" bson:"mqqt"`
	Invid     string             `json:"invid" bson:"invid"`
	Guid      string             `json:"guid" bson:"guid"`
	MsgId     string             `json:"msg_id" bson:"msg_id"`
	Text      string             `json:"text" bson:"text"`
	Context   string             `json:"context" bson:"context"`
	ClassMsg  string             `json:"class_msg" bson:"class_msg"`
	Level     int                `json:"level" bson:"level"`
	Area      string             `json:"area" bson:"area"`
	Addr      string             `json:"addr" bson:"addr"`
	Block     string             `json:"block" bson:"block"`
	Type      string             `json:"type" bson:"type"`
	Bit       string             `json:"bit" bson:"bit"`
	InvertBit string             `json:"invert_bit" bson:"invert_bit"`
	CreateAt  time.Time          `json:"create_at" bson:"create_at"`
}

type ProcessedFile struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	ProcessedAt time.Time          `json:"processed_at" bson:"processed_at"`
	Status      string             `json:"status" bson:"status"`
	MsgError    string             `json:"msg_error" bson:"msg_error"`
}

type PaginationResult struct {
	Data  []*DeviceLogs `json:"data" bson:"data"`
	Total int64         `json:"total" bson:"total"`
	Page  int           `json:"page" bson:"page"`
	Limit int           `json:"limit" bson:"limit"`
}
