package models

type Privileges struct {
	CanSendMessages      bool
	CanViewAuditLog      bool
	CanGetMessageHistory bool
	CanManageChannel     bool
	CanChangeNicks       bool
	CanKick              bool
	CanBan               bool
	CanInvite            bool
	CanMute              bool
	CanMakeRoles         bool
	CanGiveRoles         bool
}
