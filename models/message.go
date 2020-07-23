package models

type Message struct {
	ID        string
	UserID    string
	ChannelID string
	ServerID  string
	Comments  []Comment
	Reactions []Reaction
	Content   string
}
