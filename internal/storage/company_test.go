package storage

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func (ts *StorageTestSuite) TestCompany() {
	ctx := context.Background()

	ts.T().Run("Create", func(t *testing.T) {
		assert := require.New(t)

		uuid1, err := uuid.NewV4()
		assert.Nil(err)
		uuid2, err := uuid.NewV4()
		assert.Nil(err)

		c1 := Company{
			ID:           uuid1,
			Name:         "test_company_1",
			Description:  "description for company 1",
			EmployeesCnt: 11,
			Registered:   true,
			Type:         1,
		}

		c2 := Company{
			ID:           uuid2,
			Name:         "test_company_2",
			Description:  "description for company 2",
			EmployeesCnt: 111,
			Registered:   false,
			Type:         2,
		}

		assert.Nil(CreateCompany(context.Background(), ts.Tx(), &c1))
		assert.Nil(CreateCompany(context.Background(), ts.Tx(), &c2))

		c1.CreatedAt = time.Now().Round(time.Second).UTC()
		c1.UpdatedAt = time.Now().Round(time.Second).UTC()
		c2.CreatedAt = time.Now().Round(time.Second).UTC()
		c2.UpdatedAt = time.Now().Round(time.Second).UTC()

		t.Run("Get", func(t *testing.T) {
			resp1, err := GetCompany(ctx, ts.Tx(), c1.ID)
			assert.Nil(err)

			resp2, err := GetCompany(ctx, ts.Tx(), c2.ID)
			assert.Nil(err)

			resp1.CreatedAt = resp1.CreatedAt.Round(time.Second).UTC()
			resp1.UpdatedAt = resp1.UpdatedAt.Round(time.Second).UTC()
			resp2.CreatedAt = resp2.CreatedAt.Round(time.Second).UTC()
			resp2.UpdatedAt = resp2.UpdatedAt.Round(time.Second).UTC()

			assert.Equal(c1, resp1)
			assert.Equal(c2, resp2)
		})

		t.Run("Update", func(t *testing.T) {
			assert := require.New(t)

			upd := &Company{
				ID:           c1.ID,
				Name:         "new_name",
				Description:  "another description",
				EmployeesCnt: 12,
				Registered:   false,
				Type:         3,
			}

			assert.Nil(UpdateCompany(ctx, ts.Tx(), upd))
			upd.UpdatedAt = time.Now().Round(time.Second).UTC()
			upd.CreatedAt = c1.CreatedAt

			dGet, err := GetCompany(ctx, ts.Tx(), upd.ID)
			assert.Nil(err)

			dGet.CreatedAt = dGet.CreatedAt.Round(time.Second).UTC()
			dGet.UpdatedAt = dGet.UpdatedAt.Round(time.Second).UTC()

			assert.Equal(*upd, dGet)
		})

		t.Run("Delete", func(t *testing.T) {
			assert := require.New(t)

			assert.Nil(DeleteCompany(context.Background(), ts.Tx(), c1.ID))
			assert.Equal(ErrDoesNotExist, DeleteCompany(context.Background(), ts.Tx(), c1.ID))
			_, err := GetCompany(ctx, ts.Tx(), c1.ID)
			assert.Equal(ErrDoesNotExist, err)

			assert.Nil(DeleteCompany(context.Background(), ts.Tx(), c2.ID))
			assert.Equal(ErrDoesNotExist, DeleteCompany(context.Background(), ts.Tx(), c2.ID))
			_, err = GetCompany(ctx, ts.Tx(), c2.ID)
			assert.Equal(ErrDoesNotExist, err)
		})
	})
}
