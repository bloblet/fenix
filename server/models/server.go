package models

type Server struct {
	ServerID string
	OwnerID  string
	Members  []string
	// Spaces []Space
	Channels []Channel
}
