package models

import "github.com/kamva/mgm/v3"

type RoleList struct {
	SyncModel        `bson:",inline"`
	mgm.DefaultModel `bson:",inline"`
	ChannelID        string
	Roles            map[string]Role
}

func (c *RoleList) CollectionName() string {
	return "rolelist"
}