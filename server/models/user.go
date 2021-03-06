package models

import (
	"encoding/json"
	"github.com/bloblet/fenix/server/utils"
	log "github.com/sirupsen/logrus"
)

// User is the current datatype for fenix users.
type User struct {
	ID            string
	Token         string
	Email         string
	Salt          []byte `json:"-"`
	Password      []byte `json:"-"`
	Username      string
	Discriminator string
	Servers       []string
	Friends       []string
	Activity      Activity
	Settings      UserSettings
}

// ToJSON converts the user to JSON
func (user *User) ToJSON() string {
	b, err := json.Marshal(user)
	if err != nil {
		utils.Log().WithFields(
			log.Fields{
				"byteLength": len(b),
				"error":      err,
			},
		).Error("Encountered an error while serializing a User object.")
	}

	return string(b)
}
