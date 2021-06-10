package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type Invite struct {
	SyncModel        `bson:",inline"`
	mgm.DefaultModel `bson:",inline"`
	ChannelID        string
	Timestamp        time.Time
	TimesUsed        int64
	CreatedBy        string
}

func (i *Invite) CollectionName() string {
	return "invites"
}
