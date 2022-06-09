package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fancar/tmp_xm/internal/storage"
	"github.com/fancar/tmp_xm/internal/test"
)

type CompanyAPITestSuite struct {
	suite.Suite
	api CompanyServiceServer
}

func TestCompanyServiceServer(t *testing.T) {
	suite.Run(t, new(CompanyAPITestSuite))
}

func (ts *CompanyAPITestSuite) SetupSuite() {

	assert := require.New(ts.T())
	conf := test.GetConfig()

	assert.NoError(storage.Setup(conf))
	assert.NoError(storage.MigrateDown(storage.DB().DB))
	assert.NoError(storage.MigrateUp(storage.DB().DB))

	ts.api = NewCompanyAPI()
}

func (ts *CompanyAPITestSuite) TestCompany() {
	c := Company{
		Name:    "test_name",
		Code:    "167000",
		Country: "Uruguay",
		Website: "http://uruguay.is.cool",
	}

	ts.T().Run("Create", func(t *testing.T) {

		assert := require.New(t)
		_, err := ts.api.Create(
			context.Background(),
			&CreateCompanyRequest{
				Company: &c,
			},
		)
		assert.Nil(err)

		c.Name = "second_company"
		_, err = ts.api.Create(
			context.Background(),
			&CreateCompanyRequest{
				Company: &c,
			},
		)
		assert.Nil(err)

		t.Run("List", func(t *testing.T) {
			assert := require.New(t)

			r := ListCompanyRequest{
				Name: "test_name",
			}

			getResp, err := ts.api.List(context.Background(), &r)
			assert.Nil(err)
			assert.Len(getResp.Result, 1)
			assert.Equal(getResp.Result[0].Name, r.Name)
			assert.Equal(getResp.Result[0].Code, c.Code)
			assert.Equal(getResp.Result[0].Country, c.Country)
			assert.Equal(getResp.Result[0].Website, c.Website)
			c.Id = getResp.Result[0].Id

			r = ListCompanyRequest{
				Code: c.Code,
			}

			getResp, err = ts.api.List(context.Background(), &r)
			assert.Nil(err)
			assert.Len(getResp.Result, 2)

			// assert.Equal(getResp, getResp.Company)
		})

		t.Run("Update", func(t *testing.T) {
			assert := require.New(t)

			c.Name = "test_name_changed"
			c.Code = "167000-changed"
			c.Country = "Paraguay"
			c.Website = "http://paraguay.isnt.cool"
			_, err = ts.api.Update(
				context.Background(),
				&UpdateCompanyRequest{
					Company: &c,
				},
			)
			assert.Nil(err)

			r := GetCompanyRequest{
				Id: c.Id,
			}

			getResp, err := ts.api.Get(context.Background(), &r)
			assert.Nil(err)
			assert.Equal(getResp.Company, &c)
		})

		t.Run("Delete", func(t *testing.T) {
			assert := require.New(t)

			r := DeleteCompanyRequest{
				Id: c.Id,
			}
			_, err := ts.api.Delete(context.Background(), &r)
			assert.Nil(err)

			gr := GetCompanyRequest{
				Id: c.Id,
			}
			_, err = ts.api.Get(context.Background(), &gr)
			assert.NotNil(err)
		})
	})
}
