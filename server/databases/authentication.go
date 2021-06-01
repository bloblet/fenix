package databases

import (
	"crypto"
	"crypto/rand"
	"crypto/subtle"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/bloblet/fenix/server/models"
	"github.com/bloblet/fenix/server/utils"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"github.com/piecubed/twofactor"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/gomail.v2"
	"gopkg.in/mgo.v2/bson"
	mrand "math/rand"
	"time"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func NewAuthenticationManager() *AuthenticationManager {
	a := &AuthenticationManager{}
	a.VerificationListeners = map[string]chan bool{}
	a.LastVerificationEmailSentAt = map[string]time.Time{}
	return a
}

type AuthenticationManager struct {
	VerificationListeners       map[string]chan bool
	LastVerificationEmailSentAt map[string]time.Time
}

func (a *AuthenticationManager) hashPassword(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, 400_000, 64, crypto.SHA512.New)
}

func (a *AuthenticationManager) checkPassword(rawPassword, hashedPassword, salt []byte) bool {
	res := subtle.ConstantTimeCompare(a.hashPassword(rawPassword, salt), hashedPassword)
	return res == 1
}

func (a *AuthenticationManager) NewUser(email, username, password string, sync ...bool) (*models.User, error) {
	res := mgm.Coll(&models.User{}).FindOne(mgm.Ctx(), bson.M{
		"email": bson.M{operator.Eq: email},
	})

	if res.Err() == nil {
		return nil, UserExistsError{}
	}

	_sync := false

	if len(sync) != 0 {
		_sync = sync[0]
	}

	salt, err := GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}

	ott, err := GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}

	mfaSecret, err := GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}

	t, err := a.CreateToken()
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Tokens: map[string]models.Token{
			t.TokenID: *t,
		},
		Email:         email,
		Salt:          salt,
		Password:      a.hashPassword([]byte(password), salt),
		OTT:           base64.URLEncoding.EncodeToString(ott),
		EmailVerified: false,
		Username:      username,
		Discriminator: fmt.Sprint(mrand.Intn(9999)),
		MFAEnabled:    true,
		AuthSecret:    base64.URLEncoding.EncodeToString(mfaSecret),
	}

	u.New()
	err = mgm.Coll(u).Create(u)

	if _sync {
		u.WaitForSave()
	}

	if err != nil {
		return nil, err
	}

	err = a.SendVerificationEmail(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (a *AuthenticationManager) GetUser(userID string) (*models.User, bool) {
	u := &models.User{}
	err := mgm.Coll(u).FindByID(userID, u)

	if err != nil {
		return nil, false
	}

	return u, true
}

func (a *AuthenticationManager) PasswordAuthenticateUser(email, password string) (*models.User, bool) {
	u := &models.User{}

	res := mgm.Coll(u).FindOne(mgm.Ctx(), bson.M{
		"email": bson.M{operator.Eq: email},
	})

	if res.Err() != nil {
		return nil, false
	}

	err := res.Decode(u)

	if err != nil {
		return nil, false
	}

	return u, a.checkPassword([]byte(password), u.Password, u.Salt)
}

func (a *AuthenticationManager) TokenAuthenticateUser(userID, token, tokenID string) (*models.User, bool) {
	u := &models.User{}

	err := mgm.Coll(u).FindByID(userID, u)

	if err != nil {
		return nil, false
	}

	t, ok := u.Tokens[tokenID]

	if !ok {
		return nil, false
	}

	if t.Expires.Before(time.Now()) {
		delete(u.Tokens, tokenID)
		mgm.Coll(u).Update(u)
		return nil, false
	}

	res := subtle.ConstantTimeCompare([]byte(token), []byte(t.Token))

	if res == 1 {
		return u, true
	}

	delete(u.Tokens, tokenID)
	mgm.Coll(u).Update(u)
	return nil, false
}

func (a *AuthenticationManager) SendVerificationEmail(u *models.User) error {
	if t, ok := a.LastVerificationEmailSentAt[u.ID.Hex()]; ok && t.Before(time.Now().Add(time.Minute)) {
		return VerificationEmailCooldown{}
	}

	server := config.Authenticator.MailServer
	port := config.Authenticator.MailServerPort
	user := config.Authenticator.MailServerUser
	password := config.Authenticator.MailServerPassword
	verificationEndpoint := config.Authenticator.VerificationEndpoint

	if server == "" || port == 0 || user == "" || password == "" || verificationEndpoint == "" {
		utils.Log().Error("Mail server not configured")
	}
	body := fmt.Sprintf("Verify your fenix account by clicking the link %v?u=%v&t=%v", verificationEndpoint, u.ID.Hex(), u.OTT)

	m := gomail.NewMessage()
	m.SetHeader("To", u.Email)
	m.SetHeader("Subject", "Fenix Verification")
	m.SetBody("text/html", body)
	m.SetHeader("From", user)

	d := gomail.NewDialer(server, port, user, password)
	d.SSL = false
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         server,
	}

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}
	a.LastVerificationEmailSentAt[u.ID.Hex()] = time.Now()
	return nil
}

func (a *AuthenticationManager) Verify(ott, userID string) bool {
	u := &models.User{}
	err := mgm.Coll(u).FindByID(userID, u)
	if err != nil {
		return false
	}

	if u.EmailVerified {
		return true
	}

	if u.OTT == "" {
		return false
	}

	res := subtle.ConstantTimeCompare([]byte(ott), []byte(u.OTT))

	if res == 1 {
		u.EmailVerified = true
		res, err := mgm.Coll(u).UpdateOne(
			mgm.Ctx(),
			bson.M{
				"_id": bson.M{
					operator.Eq: u.ID,
				},
			},
			bson.M{
				"$set": bson.M{
					"emailverified": true,
				},
			},
		)

		if err != nil || res.MatchedCount != 1 {
			utils.Log().WithFields(
				log.Fields{
					"Provided OTT": ott,
					"Provided UID": userID,
					"User OTT":     u.OTT,
					"UserID":       u.ID.Hex(),
					"error":        err,
				},
			).Error("Error updating user")
		}

		c, ok := a.VerificationListeners[u.ID.Hex()]
		if ok {
			c <- true
		}

		delete(a.LastVerificationEmailSentAt, u.ID.Hex())
	}

	return res == 1
}

func (a *AuthenticationManager) Generate2FALink(u *models.User) (string, error) {
	otp, err := twofactor.NewTOTP(u.Email, "Fenix", crypto.SHA512, 8)

	if err != nil {
		return "", err
	}

	return otp.Url()
}

func (a *AuthenticationManager) DeleteToken(u *models.User, tokenID string) error {
	delete(u.Tokens, tokenID)
	return mgm.Coll(u).Update(u)
}

func (a *AuthenticationManager) CreateToken() (*models.Token, error) {
	token, err := GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	t := &models.Token{
		Token:  base64.URLEncoding.EncodeToString(token),
		Expires: time.Now().Add(time.Hour * 24 * 7),
		TokenID: tokenID.String(),
	}

	return t, nil
}

func (a *AuthenticationManager) AddToken(u *models.User, t *models.Token) error {
	u.Tokens[t.TokenID] = *t
	return mgm.Coll(u).Update(u)
}

func (a *AuthenticationManager) ChangeMFA(u *models.User, status bool) error {
	u.MFAEnabled = status
	return mgm.Coll(u).Update(u)
}

func (a *AuthenticationManager) ChangeUsername(u *models.User, newUsername string) error {
	u.Username = newUsername
	u.Discriminator = fmt.Sprint(mrand.Intn(9999))
	return mgm.Coll(u).Update(u)
}

func (a *AuthenticationManager) ChangePassword(u *models.User, newPassword string) (*models.Token, error) {
	salt, err := GenerateRandomBytes(16)
	if err != nil {
		return nil, err
	}

	t, err := a.CreateToken()
	if err != nil {
		return nil, err
	}

	u.Salt = salt
	u.Password = a.hashPassword([]byte(newPassword), salt)
	u.Tokens = map[string]models.Token{t.TokenID: *t}

	return t, mgm.Coll(u).Update(u)
}

type UserExistsError struct {
}

func (e UserExistsError) Error() string {
	return "UserExists"
}

type VerificationEmailCooldown struct {
}

func (e VerificationEmailCooldown) Error() string {
	return "VerificationEmailCooldown"
}
