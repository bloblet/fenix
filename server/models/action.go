package models

import "time"

type Action struct {
	UserID    string
	Method    string
	Timestamp time.Time
}