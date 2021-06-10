package models

import "github.com/kamva/mgm/v3"

type InviteList struct {
	SyncModel        `bson:",inline"`
	mgm.DefaultModel `bson:",inline"`
	ChannelID        string
	Invites          []string
}

func (i *InviteList) CollectionName() string {
	return "invitelist"
}