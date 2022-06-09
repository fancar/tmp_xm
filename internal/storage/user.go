package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	// "github.com/pkg/errors"
	// log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
)

// saltSize defines the salt size
const saltSize = 16

// defaultSessionTTL defines the default session TTL
const defaultSessionTTL = time.Hour * 24

// Any printable characters, at least 6 characters.
var passwordValidator = regexp.MustCompile(`^.{6,}$`)

// User defines the user structure.
type User struct {
	ID           int64     `db:"id"`
	IsAdmin      bool      `db:"is_admin"`
	SessionTTL   int32     `db:"session_ttl"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	PasswordHash string    `db:"password_hash"`
	Username     string    `db:"username"`
}

// LoginUserByPassword returns a JWT token for the user matching the given email
// and password combination.
func LoginUserByPassword(ctx context.Context, db sqlx.Ext, name string, password string) (string, error) {
	// get the user by email
	var user User
	err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			username = $1
	`, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrInvalidUsernameOrPassword
		}
		return "", fmt.Errorf("select error %v", err)
	}

	// Compare the passed in password with the hash in the database.
	if !hashCompare(password, user.PasswordHash) {
		return "", ErrInvalidUsernameOrPassword
	}

	return GetUserToken(user)
}

// hashCompare verifies that passed password hashes to the same value as the
// passed passwordHash.
func hashCompare(password string, passwordHash string) bool {
	// SPlit the hash string into its parts.
	hashSplit := strings.Split(passwordHash, "$")

	// Get the iterations and the salt and use them to encode the password
	// being compared.cre
	iterations, _ := strconv.Atoi(hashSplit[2])
	salt, _ := base64.StdEncoding.DecodeString(hashSplit[3])
	newHash := hashWithSalt(password, salt, iterations)
	return newHash == passwordHash
}

// GetUserToken returns a JWT token for the given user.
func GetUserToken(u User) (string, error) {
	// Generate the token.
	now := time.Now()
	nowSecondsSinceEpoch := now.Unix()
	var expSecondsSinceEpoch int64
	if u.SessionTTL > 0 {
		expSecondsSinceEpoch = nowSecondsSinceEpoch + (3600 * int64(u.SessionTTL))
	} else {
		expSecondsSinceEpoch = nowSecondsSinceEpoch + int64(defaultSessionTTL/time.Second)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "as",
		"aud":      "as",
		"nbf":      nowSecondsSinceEpoch,
		"exp":      expSecondsSinceEpoch,
		"sub":      "user",
		"id":       u.ID,
		"username": u.Username,
	})

	jwt, err := token.SignedString(jwtsecret)
	if err != nil {
		return jwt, fmt.Errorf("get jwt signed string error %v", err)
	}
	return jwt, err
}

// Generate the hash of a password for storage in the database.
// NOTE: We store the details of the hashing algorithm with the hash itself,
// making it easy to recreate the hash for password checking, even if we change
// the default criteria here.
func hash(password string, saltSize int, iterations int) (string, error) {
	// Generate a random salt value, 128 bits.
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("read random bytes error %v", err)
	}

	return hashWithSalt(password, salt, iterations), nil
}

func hashWithSalt(password string, salt []byte, iterations int) string {
	// Generate the hash.  This should be a little painful, adjust ITERATIONS
	// if it needs performance tweeking.  Greatly depends on the hardware.
	// NOTE: We store these details with the returned hash, so changes will not
	// affect our ability to do password compares.
	hash := pbkdf2.Key([]byte(password), salt, iterations, sha512.Size, sha512.New)

	// Build up the parameters and hash into a single string so we can compare
	// other string to the same hash.  Note that the hash algorithm is hard-
	// coded here, as it is above.  Introducing alternate encodings must support
	// old encodings as well, and build this string appropriately.
	var buffer bytes.Buffer

	buffer.WriteString("PBKDF2$")
	buffer.WriteString("sha512$")
	buffer.WriteString(strconv.Itoa(iterations))
	buffer.WriteString("$")
	buffer.WriteString(base64.StdEncoding.EncodeToString(salt))
	buffer.WriteString("$")
	buffer.WriteString(base64.StdEncoding.EncodeToString(hash))

	return buffer.String()
}

// SetPasswordHash hashes the given password and sets it.
func (u *User) SetPasswordHash(pw string) error {
	if !passwordValidator.MatchString(pw) {
		return ErrUserPasswordLength
	}

	pwHash, err := hash(pw, saltSize, HashIterations)
	if err != nil {
		return err
	}

	u.PasswordHash = pwHash

	return nil
}

// CreateUser creates the given user.
func CreateUser(ctx context.Context, db sqlx.Queryer, user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := sqlx.Get(db, &user.ID, `
		insert into "user" (
			created_at,
			updated_at,
			username,
			password_hash,
			session_ttl,
			is_admin
		)
		values (
			$1, $2, $3, $4, $5, $6)
		returning
			id`,
		user.CreatedAt,
		user.UpdatedAt,
		user.Username,
		user.PasswordHash,
		user.SessionTTL,
		user.IsAdmin,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	return nil
}
