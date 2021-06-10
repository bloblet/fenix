package models

import "github.com/kamva/mgm/v3"

type AuditList struct {
	SyncModel        `bson:",inline"`
	mgm.DefaultModel `bson:",inline"`
	ChannelID        string
	Actions          []Action
}

func (a *AuditList) CollectionName() string {
	return "auditlist"
}