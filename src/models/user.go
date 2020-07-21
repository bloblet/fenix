package models

// User is the current datatype for fenix users.
type User struct {
	ID string
	Username string
	Servers []Server
	Friends []Friend
	Activity Activity
	Settings UserSettings
}