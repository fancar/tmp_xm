package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func (ts *StorageTestSuite) TestCompany() {
	// assert := require.New(ts.T())
	ctx := context.Background()

	ts.T().Run("Create", func(t *testing.T) {
		assert := require.New(t)

		c := Company{
			Name:    "test_name",
			Code:    "167000",
			Country: "Uruguay",
			Website: "http://uruguay.is.cool",
		}

		assert.Nil(CreateCompany(context.Background(), ts.Tx(), &c))
		c.Name = "one more"
		assert.Nil(CreateCompany(context.Background(), ts.Tx(), &c))

		c.CreatedAt = time.Now().Round(time.Second).UTC()
		c.UpdatedAt = time.Now().Round(time.Second).UTC()

		t.Run("Get", func(t *testing.T) {
			f := CompanyFilters{
				Limit: 100,
			}
			resp, err := GetCompanies(ctx, ts.Tx(), f)
			assert.Nil(err)
			assert.Len(resp, 2)
			c.ID = resp[1].ID

			f.Name = c.Name
			resp, err = GetCompanies(ctx, ts.Tx(), f)
			assert.Nil(err)
			assert.Len(resp, 1)
			dGet := resp[0]

			dGet.CreatedAt = dGet.CreatedAt.Round(time.Second).UTC()
			dGet.UpdatedAt = dGet.UpdatedAt.Round(time.Second).UTC()

			assert.Equal(c, dGet)
		})

		t.Run("Update", func(t *testing.T) {
			assert := require.New(t)

			upd := Company{
				ID:      c.ID,
				Name:    "new_name",
				Code:    "100100",
				Country: "Paraguay",
				Website: "http://Paraguay.is.not.dat.cool",
			}

			assert.Nil(UpdateCompany(ctx, ts.Tx(), upd))
			upd.UpdatedAt = time.Now().Round(time.Second).UTC()
			upd.CreatedAt = c.CreatedAt

			dGet, err := GetCompany(ctx, ts.Tx(), upd.ID)
			assert.Nil(err)

			dGet.CreatedAt = dGet.CreatedAt.Round(time.Second).UTC()
			dGet.UpdatedAt = dGet.UpdatedAt.Round(time.Second).UTC()

			assert.Equal(upd, dGet)
		})

		t.Run("Delete", func(t *testing.T) {
			assert := require.New(t)

			assert.Nil(DeleteCompany(context.Background(), ts.Tx(), c.ID))
			assert.Equal(ErrDoesNotExist, DeleteCompany(context.Background(), ts.Tx(), c.ID))
			_, err := GetCompany(ctx, ts.Tx(), c.ID)
			assert.Equal(ErrDoesNotExist, err)
		})
	})
}
