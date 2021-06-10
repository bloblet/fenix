package models

import (
	"github.com/kamva/mgm/v3"
)

type Channel struct {
	SyncModel        `bson:",inline"`
	mgm.DefaultModel `bson:",inline"`
	AuditLogID       string
	UserListID       string
	InviteListID     string
	RolesListID      string
	Bans             map[string]string
	Settings         Settings
	ChannelName      string
}

func (c *Channel) CollectionName() string {
	return "channels"
}
