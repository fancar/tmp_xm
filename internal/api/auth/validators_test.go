package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"

	"github.com/fancar/tmp_xm/internal/storage"
	"github.com/fancar/tmp_xm/internal/test"
)

type validatorTest struct {
	Name       string
	Claims     Claims
	Validators []ValidatorFunc
	ExpectedOK bool
}

type ValidatorTestSuite struct {
	suite.Suite
}

func TestValidators(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}

func (ts *ValidatorTestSuite) SetupSuite() {
	assert := require.New(ts.T())

	conf := test.GetConfig()
	assert.NoError(storage.Setup(conf))
}

func (ts *ValidatorTestSuite) SetupTest() {
	assert := require.New(ts.T())

	assert.NoError(storage.MigrateDown(storage.DB().DB))
	assert.NoError(storage.MigrateUp(storage.DB().DB))
}

func (ts *ValidatorTestSuite) CreateUser(username string, isActive, isAdmin bool) (int64, error) {
	u := storage.User{
		IsAdmin:  isAdmin,
		Username: username,
	}

	err := storage.CreateUser(context.Background(), storage.DB(), &u)
	return u.ID, err
}

func (ts *ValidatorTestSuite) RunTests(t *testing.T, tests []validatorTest) {
	for _, tst := range tests {
		t.Run(tst.Name, func(t *testing.T) {
			assert := require.New(t)

			if tst.Claims.Username != "" || tst.Claims.UserID != 0 {
				tst.Claims.Subject = "user"
			} else {
				tst.Claims.Subject = "api_key"
			}

			for _, v := range tst.Validators {
				ok, err := v(storage.DB(), &tst.Claims)
				assert.NoError(err)
				assert.Equal(tst.ExpectedOK, ok)
			}
		})
	}
}

func (ts *ValidatorTestSuite) TestUser() {
	assert := require.New(ts.T())

	users := []struct {
		id       int64
		username string
		isActive bool
		isAdmin  bool
	}{
		// {username: "admin", isAdmin: true},
		{username: "user", isAdmin: false},
	}
	for i, user := range users {
		id, err := ts.CreateUser(user.username, user.isActive, user.isAdmin)
		assert.NoError(err)
		users[i].id = id
	}

	ts.T().Run("User", func(t *testing.T) {
		tests := []validatorTest{
			{
				Name:       "user",
				Validators: []ValidatorFunc{ValidateActiveUser()},
				Claims:     Claims{UserID: users[0].id},
				ExpectedOK: true,
			},
			{
				Name:       "invalid user",
				Validators: []ValidatorFunc{ValidateActiveUser()},
				Claims:     Claims{UserID: 9999},
				ExpectedOK: false,
			},
		}

		ts.RunTests(t, tests)
	})
}
