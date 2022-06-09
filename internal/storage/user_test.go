package storage

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func (ts *StorageTestSuite) TestUser() {
	// Set a user secret so JWTs can be assigned
	jwtsecret = []byte("DoWahDiddy")

	ts.T().Run("Create with invalid password", func(t *testing.T) {
		assert := require.New(t)

		user := User{
			IsAdmin:    false,
			SessionTTL: 40,
			Username:   "foo@bar.com",
		}
		err := user.SetPasswordHash("bad")
		assert.Equal(ErrUserPasswordLength, errors.Cause(err))
	})

	ts.T().Run("Create", func(t *testing.T) {
		assert := require.New(t)

		user := User{
			IsAdmin:    false,
			SessionTTL: 20,
			Username:   "foo@bar.com",
		}
		password := "somepassword"
		assert.NoError(user.SetPasswordHash(password))

		err := CreateUser(context.Background(), DB(), &user)
		assert.NoError(err)

		t.Run("LoginUserByPassword", func(t *testing.T) {
			assert := require.New(t)

			jwt, err := LoginUserByPassword(context.Background(), DB(), user.Username, password)
			assert.NoError(err)
			assert.NotEqual("", jwt)
		})

		t.Run("GetUserToken", func(t *testing.T) {
			assert := require.New(t)

			token, err := GetUserToken(user)
			assert.NoError(err)
			assert.NotEqual("", token)
		})
	})
}
