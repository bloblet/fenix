package models

import "github.com/kamva/mgm/v3"

type UserList struct {
	SyncModel        `bson:",inline"`
	mgm.DefaultModel `bson:",inline"`
	ChannelID        string
	Users            []ChannelUser
}

func (a *UserList) CollectionName() string {
	return "userlist"
}