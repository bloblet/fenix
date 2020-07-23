package models

import "time"

// Activity represents an activity a user is displaying on their profile
type Activity struct {
	StopAt time.Time
	Name   string
	Emoji  RawEmoji
}
