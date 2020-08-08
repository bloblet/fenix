package models

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
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
	PublicKey     ecdsa.PublicKey
}

// ToJSON converts the user to JSON
func (user *User) ToJSON() string {
	b, err := json.Marshal(user)
	if err != nil {
		fmt.Print("Encountered an error while serializing a User object.  ")
		fmt.Print(err)
		fmt.Print("\n")
		panic(err)
	}

	return string(b)
}

// Verify a message was signed with a user's key
func (user *User) Verify(hash []byte, r, s *big.Int) bool {
	return ecdsa.Verify(&user.PublicKey, hash, r, s)
}
