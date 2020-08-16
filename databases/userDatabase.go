package databases

import (
	"context"
	crypto_rand "crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	models "fenix/models/database"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"

	// THANK YOU ETCD
	"github.com/coreos/etcd/clientv3/concurrency"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
)

var authDB = "/auth/"
var userDB = "/user/"
var nameDB = "/name/"

// NewUserDatabase makes a new user database
func NewUserDatabase(username, password string, testing bool, prefix string) UserDatabase {
	if testing {
		authDB = "/testing/" + prefix + "/auth/"
		userDB = "/testing/" + prefix + "/user/"
		nameDB = "/testing/" + prefix + "/name/"
	}

	db := UserDatabase{}
	db.username = username
	db.password = password
	return db
}

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

func fatal(err error) {
	now := time.Time{}
	fmt.Print(now)
	fmt.Print("\u001b[31m FATAL: ")
	fmt.Print(err)
	fmt.Print("\u001b[0m\n")
	panic(err)
}

// UserDatabase manages the user database
type UserDatabase struct {
	username string
	password string
}

// Database opens a Database connection.  DO NOT FORGET TO defer cli.Close()
func (db *UserDatabase) Database() (*clientv3.Client, error) {
	var options []grpc.DialOption
	options = append(options, grpc.WithBlock(), grpc.WithTimeout(dialTimeout))
	c := zap.NewProductionConfig()
	c.Level = zap.NewAtomicLevel()
	c.Level.SetLevel(zap.FatalLevel)

	return clientv3.New(clientv3.Config{
		LogConfig:   &c,
		DialOptions: options,
		DialTimeout: dialTimeout,
		Username:    db.username,
		Password:    db.password,
		Endpoints:   []string{"98.212.66.76:2379"},
	})
}

func (db *UserDatabase) padDiscriminator(d string) string {
	// Pad d so that its in the form of 0001
	if len(d) != 4 {
		for len(d) != 4 {
			d = "0" + d
		}
	}
	return d
}

func (db *UserDatabase) sanitize(target string) string {
	// Uses base64url encoding to sanitize client data, so it doesn't mess with paths.
	return base64.URLEncoding.EncodeToString([]byte(target))
}

func (db *UserDatabase) unsanitize(target string) (string, error) {
	r, err := base64.URLEncoding.DecodeString(target)
	return string(r), err
}

// UserExists checks the userDB to see if a user exists.
func (db *UserDatabase) UserExists(email string, cli *clientv3.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	kv := clientv3.NewKV(cli)

	res, err := kv.Get(ctx, authDB+email)
	cancel()

	if err != nil {
		fmt.Print("UserExists error, ")
		fmt.Println(err)
		return false
	}

	return res.Count == 1
}

// CreateUser will create a user, and pack it into a User object
func (db *UserDatabase) CreateUser(email, password, username string) (models.User, error) {
	// Get database client
	cli, dberr := db.Database()
	defer cli.Close()

	// Make sure there wasn't any problems
	if dberr != nil {
		fatal(dberr)
		return models.User{}, dberr
	}

	// Sanitize user provided info
	email = db.sanitize(email)
	username = db.sanitize(username)

	// Open our concurrency session
	s, lockErr := concurrency.NewSession(cli, concurrency.WithContext(context.Background()))
	defer s.Close()

	// Make sure there wasn't any problems
	if lockErr != nil {
		fatal(lockErr)
		return models.User{}, lockErr
	}

	// Aquire locks
	nameLock := concurrency.NewMutex(s, nameDB+username)
	nameLock.Lock(context.Background())

	// The authDB lock is to prevent race conditions with 2 emails being registered at once
	authLock := concurrency.NewMutex(s, authDB+email)
	authLock.Lock(context.Background())

	// Check if the email is already taken
	if db.UserExists(email, cli) {
		go nameLock.Unlock(context.Background())
		go authLock.Unlock(context.Background())
		return models.User{}, UserExists{}
	}
	// Check the username db for enteries
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := cli.Get(ctx, nameDB+username)
	// cancel()

	// Make sure there wasn't any problems
	if err != nil {
		fatal(err)
	}

	var discrims []int

	// This username has never been used before.
	if res.Count == 0 {
		discrims = make([]int, 10000)

		for i := 0; i < 10000; i++ {
			discrims[i] = i
		}
	} else {
		// Unmarshal the JSON into a map
		json.Unmarshal(res.Kvs[0].Value, &discrims)
	}
	// There *could* be no limit here, but say you were impersonating a popular account, this is good.
	// Besides, Fenix probably won't ever get over 9999 users...
	// -1 is the special key for not having any more discriminators.
	if discrims[0] == -1 {
		return models.User{}, NoMoreDiscriminators{}
	}

	discrimKeys := make([]int, len(discrims))

	i := 0
	for k := range discrims {
		discrimKeys[i] = k
		i++
	}

	pos := rand.Intn(len(discrims) - 1)
	d := discrimKeys[pos]
	discrims = append(discrims[:pos], discrims[pos+1:]...)

	b, _ := json.Marshal(discrims)

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Put(ctx, nameDB+username, string(b))

	if err != nil {
		fatal(err)
	}

	cancel()
	err = nameLock.Unlock(context.Background())

	if err != nil {
		fatal(err)
	}

	user := models.User{}
	user.Discriminator = db.padDiscriminator(strconv.Itoa(d))
	uid, _ := uuid.NewRandom()
	user.Activity = models.Activity{}
	user.ID = uid.String()
	user.Salt = []byte(generateToken())
	user.Token = generateToken()
	user.Password = pbkdf2.Key([]byte(password), user.Salt, 200000, 64, sha512.New)
	user.Settings = models.UserSettings{}
	user.Username, _ = db.unsanitize(username)
	user.Email, _ = db.unsanitize(email)

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Put(ctx, authDB+email, user.ID)
	cancel()

	if err != nil {
		fatal(err)
	}

	err = authLock.Unlock(context.Background())

	if err != nil {
		fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	cli.Put(ctx, userDB+user.ID, user.ToJSON())
	cancel()

	return user, nil
}

type NoMoreDiscriminators struct{}

func (e NoMoreDiscriminators) Error() string {
	return "NoMoreDiscriminators"
}

type UserExists struct{}

func (e UserExists) Error() string {
	return "UserExists"
}
