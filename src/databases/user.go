package databases

import (
	"../models"
	"bytes"
	"encoding/gob"
	bolt "go.etcd.io/bbolt"
)

var userBucket = []byte("users")

// UserDatabase manages the user database
type UserDatabase struct {
	boltDB bolt.DB
}

// Setup sets up the database, should be called on startup only.
func (db *UserDatabase) Setup(boltDB *bolt.DB) error {
	db.boltDB = *boltDB
	err := db.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(userBucket)

		return err
	})

	return err
}
// Get will get a user from the User bucket and deserialize it.
func (db *UserDatabase) Get(UserID string) (*models.User, error) {
	var user *models.User
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		userBytes := tx.Bucket(userBucket).Get([]byte(UserID))

		if (userBytes == nil) {
			return nil
		}

		var buffer bytes.Buffer
		dec := gob.NewDecoder(&buffer)

		buffer.Write(userBytes)

		return dec.Decode(user)
	})

	return user, err
}
