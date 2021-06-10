package models

type ChannelUser struct {
	UserID   string
	Privs    Privileges
	Nickname string
	Roles    []string
}