package models
import "encoding/json"

// User is the current datatype for fenix users.
type User struct {
	ID string
	Token string
	Email string
	Salt []byte `json:"-"`
	Password []byte `json:"-"`
	Username string
	Discriminator string
	Servers []string
	Friends []string
	Activity Activity
	Settings UserSettings
}
// ToJSON converts the user to JSON
func (user *User) ToJSON() string {
	b, err := json.Marshal(user)
	if (err != nil) {
		print("Encountered an error while serializing a User object.  ")
		print(err)
		print("\n")
	}

	return string(b)
}