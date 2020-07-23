package databases

import (
	"context"
	crypto_rand "crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fenix/models"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/crypto/pbkdf2"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

const authDB = "auth/"
const userDB = "user/"
const nameDB = "name/"

func generateToken() string {
	// Get 128 random bytes (1024 bits).
	raw := make([]byte, 128)
	_, err := crypto_rand.Read(raw)
	if err != nil {
		panic(err)
	}

	// Encode it with base64.RawURLEncoding
	var builder strings.Builder
	encoder := base64.NewEncoder(base64.RawURLEncoding, &builder)
	encoder.Write(raw)
	encoder.Close()

	// Return the string
	return builder.String()
}

func makeDiscriminator() string {
	d := string(rand.Intn(9998) + 1)

	// Pad d so that its in the form of 0001
	if len(d) != 4 {
		for len(d) != 4 {
			d = "0" + d
		}
	}
	return d
}

func userTxn(user *models.User, txn clientv3.Txn) (*clientv3.TxnResponse, error) {
	txn.If(clientv3.Compare(clientv3.Value(nameDB+user.Username+"#"+user.Discriminator), "=", ""))
	txn.Then(clientv3.OpPut(authDB+user.Email, user.ID), clientv3.OpPut(userDB+user.ID, user.ToJSON()))
	return txn.Commit()

}

func fatal(err error) {
	now := time.Time{}
	print(now)
	print("\u001b[31m FATAL: ")
	print(err)
	print("\u001b[0m\n")
}

// UserDatabase manages the user database
type UserDatabase struct {
}

// database opens a database connection.  DO NOT FORGET TO defer cli.Close()
func (db *UserDatabase) database() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"127.0.0.1:2379"},
	})
}

// UserExists checks the userDB to see if a user exists.
func (db *UserDatabase) UserExists(email string) bool {
	cli, err := db.database()
	defer cli.Close()

	if err != nil {
		fatal(err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	kv := clientv3.NewKV(cli)

	res, err := kv.Get(ctx, authDB+email)
	cancel()

	if err != nil {
		fatal(err)
		return false
	}
	print(res.Count)
	return res.Count == 1
}

// CreateUser will create a user, and pack it into a User object
func (db *UserDatabase) CreateUser(email, password, username string) (models.User, error) {
	if db.UserExists(email) {
		return models.User{}, UserExists{}
	}

	cli, err := db.database()
	defer cli.Close()

	if err != nil {
		fatal(err)
		return models.User{}, err
	}

	uid, _ := uuid.NewRandom()

	

	user := models.User{}
	user.Activity = models.Activity{}
	user.Discriminator = makeDiscriminator()
	user.ID = uid.String()
	user.Salt = []byte(generateToken())
	user.Token = generateToken()
	user.Password = pbkdf2.Key([]byte(password), user.Salt, 200000, 64, sha512.New)
	user.Settings = models.UserSettings{}
	user.Username = username
	user.Email = email

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	kv := clientv3.NewKV(cli)
	
	for i := 0; i <= 9999; i++ {
		res, _ := userTxn(&user, kv.Txn(ctx))
		if (res.Succeeded) {
			break
		}
		if (i == 9999) {
			cancel()
			return models.User{}, NoMoreDiscriminators{}
		}

		user.Discriminator = makeDiscriminator()
	}

	cancel()

	return user, nil
}

type NoMoreDiscriminators struct {}
func (e NoMoreDiscriminators) Error() string {
	return "NoMoreDiscriminators"
}

type UserExists struct {}
func (e UserExists) Error() string {
	return "UserExists"
}