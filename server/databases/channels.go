package databases

import (
	"context"
	"fmt"
	"github.com/bloblet/fenix/server/models"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ChannelDatabase struct {
}

func (c *ChannelDatabase) JoinChannel(inviteID string, user *models.User) (*models.Channel, error) {
	invite := &models.Invite{}
	err := mgm.Coll(invite).FindByID(inviteID, invite)

	if err != nil {
		return nil, err
	}

	channel := &models.Channel{}
	mgm.Coll(channel).FindByID(invite.ChannelID, channel)

	err = mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		_, err := mgm.Coll(invite).UpdateOne(sc, bson.M{
			"_id": bson.M{
				operator.Eq: invite.ID.Hex(),
			},
		},
			bson.M{
				operator.Add: bson.M{
					"timesUsed": 1,
				},
			})

		if err != nil {
			return err
		}

		_, err = mgm.Coll(&models.UserList{}).UpdateOne(sc,
			bson.M{
				"_id": bson.M{
					operator.Eq: channel.UserListID,
				},
			},
			bson.M{
				operator.Push: bson.M{
					"users": user.ID.Hex(),
				},
			},
		)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})

	return channel, err
}

func (c *ChannelDatabase) CreateChannel(name string, user *models.User) (*models.Channel, error) {
	channel := &models.Channel{
		Bans: map[string]string{},
		Settings: models.Settings{
			Owner: user.ID.Hex(),
		},
		ChannelName: name,
	}
	err := mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(channel).CreateWithCtx(sc, channel)

		if err != nil {
			return err
		}
		cID := channel.ID.Hex()

		userList := &models.UserList{
			ChannelID: cID,
			Users: []models.ChannelUser{
				{
					UserID:   user.ID.Hex(),
					Privs:    models.Privileges{},
					Nickname: "",
					Roles: []string{
						fmt.Sprintf("%v.%v", cID, 0),
					},
				},
			},
		}

		err = mgm.Coll(userList).CreateWithCtx(sc, userList)
		if err != nil {
			return err
		}

		auditList := &models.AuditList{
			ChannelID: cID,
			Actions: []models.Action{
				{
					UserID:    user.ID.Hex(),
					Method:    "CreateChannel",
					Timestamp: time.Now(),
				},
			},
		}

		err = mgm.Coll(auditList).CreateWithCtx(sc, auditList)
		if err != nil {
			return err
		}

		invite := &models.Invite{
			ChannelID: cID,
			Timestamp: time.Now(),
			TimesUsed: 0,
			CreatedBy: user.ID.Hex(),
		}

		err = mgm.Coll(invite).CreateWithCtx(sc, invite)
		if err != nil {
			return err
		}

		inviteList := &models.InviteList{
			ChannelID: cID,
			Invites: []string{
				invite.ID.Hex(),
			},
		}

		err = mgm.Coll(inviteList).CreateWithCtx(sc, inviteList)
		if err != nil {
			return err
		}

		roleList := &models.RoleList{
			ChannelID: cID,
			Roles: map[string]models.Role{
				fmt.Sprintf("%v.%v", cID, 0): {
					Privs: models.Privileges{
						CanSendMessages:      true,
						CanGetMessageHistory: true,
						CanChangeNicks:       true,
						CanInvite:            true,
					},
					Color:  "808080",
					RoleID: fmt.Sprintf("%v.%v", cID, 0),
					Name:   "everyone",
				},
			},
		}

		err = mgm.Coll(roleList).CreateWithCtx(sc, roleList)
		if err != nil {
			return err
		}

		channel.UserListID = userList.ID.Hex()
		channel.AuditLogID = auditList.ID.Hex()
		channel.RolesListID = roleList.ID.Hex()
		channel.InviteListID = inviteList.ID.Hex()
		err = mgm.Coll(channel).UpdateWithCtx(sc, channel)

		return err
	})
	return channel, err
}

func (c *ChannelDatabase) GetPrivsForUser(ctx context.Context, channelID string, user *models.User) (*models.Privileges, error) {
	channel := &models.Channel{}
	mgm.Coll(channel).FindByID(channelID, channel)

	cUser := &models.ChannelUser{}
	res := mgm.Coll(&models.UserList{}).FindOne(ctx,
		bson.M{
			"_id": bson.M{
				operator.Eq: channel.UserListID,
			},
			"users": bson.M{
				operator.Match: bson.M{
					"userID": bson.M{
						operator.Eq: user.ID.Hex(),
					},
				},
			},
		},
	)
	err := res.Decode(cUser)
	if err != nil {
		return nil, err
	}
	rolesQuery := bson.A{}

	for _, roleID := range cUser.Roles {
		rolesQuery = append(rolesQuery, roleID)
	}

	cursor, err := mgm.Coll(&models.RoleList{}).
		Find(ctx,
			bson.M{
				"_id": bson.M{
					operator.Eq: channel.RolesListID,
				},
				"roles": bson.M{
					"roleID": bson.M{
						operator.In: rolesQuery,
					},
				},
			},
			options.Find().SetSort(bson.M{"RoleID": 1}),
		)

	if err != nil {
		return nil, err
	}

	roles := make([]models.Role, 0)

	err = cursor.Decode(roles)
	if err != nil {
		return nil, err
	}

	userPrivs := &models.Privileges{}

	// Apply all role privs
	for _, role := range roles {
		c.processPrivs(userPrivs, &role.Privs)
	}

	// Apply all user specific privs.  These override role privs
	c.processPrivs(userPrivs, &cUser.Privs)
	return userPrivs, nil
}

func (c *ChannelDatabase) processPrivs(existing, new *models.Privileges) {
	if new.CanSendMessages {
		existing.CanSendMessages = true
	}
	if new.CanViewAuditLog {
		existing.CanViewAuditLog = true
	}
	if new.CanGetMessageHistory {
		existing.CanGetMessageHistory = true
	}
	if new.CanManageChannel {
		existing.CanManageChannel = true
	}
	if new.CanChangeNicks {
		existing.CanChangeNicks = true
	}
	if new.CanKick {
		existing.CanKick = true
	}
	if new.CanBan {
		existing.CanBan = true
	}
	if new.CanInvite {
		existing.CanInvite = true
	}
	if new.CanMute {
		existing.CanMute = true
	}
	if new.CanMakeRoles {
		existing.CanMakeRoles = true
	}
	if new.CanGiveRoles {
		existing.CanGiveRoles = true
	}
}

func (c *ChannelDatabase) CreateChannelInvite(channelID string, user *models.User) (*models.Invite, error) {
	invite := &models.Invite{
		ChannelID: channelID,
		Timestamp: time.Now(),
		TimesUsed: 0,
		CreatedBy: user.ID.Hex(),
	}

	channel := &models.Channel{}
	mgm.Coll(channel).FindByID(channelID, channel)

	privs, err := c.GetPrivsForUser(mgm.Ctx(), channelID, user)

	if err != nil {
		return nil, err
	}

	if !privs.CanInvite {
		return nil, NotAuthorized{Message: "You cannot perform this action!"}
	}
	err = mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		err := mgm.Coll(invite).CreateWithCtx(sc, invite)

		if err != nil {
			return err
		}

		_, err = mgm.Coll(&models.InviteList{}).UpdateOne(sc,
			bson.M{
				"_id": bson.M{
					operator.Eq: channel.InviteListID,
				},
			},
			bson.M{
				operator.Push: invite.ID.Hex(),
			},
		)
		return err
	})

	return invite, err
}

func (c *ChannelDatabase) DeleteChannelInvite() {

}

func (c *ChannelDatabase) GetChannel() {

}

func (c *ChannelDatabase) GetChannelAuditLog() {

}

func (c *ChannelDatabase) GetChannelRoles() {

}

func (c *ChannelDatabase) GetUserInfo() {

}

func (c *ChannelDatabase) GetRoles() {

}

func (c *ChannelDatabase) GetChannelInfo() {

}

func (c *ChannelDatabase) UpdateChannelInfo() {

}

func (c *ChannelDatabase) KickUser() {

}

func (c *ChannelDatabase) BanUser() {

}

func (c *ChannelDatabase) UpdateUserPrivs() {

}

func (c *ChannelDatabase) GetUsers() {

}

type NotAuthorized struct {
	Message string
}

func (e NotAuthorized) Error() string {
	return fmt.Sprintf("Not Authorized: %v", e.Message)
}
