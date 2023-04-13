package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fancar/tmp_xm/internal/kafka"
	"github.com/fancar/tmp_xm/internal/storage"
	"github.com/fancar/tmp_xm/internal/test"
	log "github.com/sirupsen/logrus"
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
	var wg sync.WaitGroup

	assert.NoError(storage.Setup(conf))
	assert.NoError(storage.MigrateDown(storage.DB().DB))
	assert.NoError(storage.MigrateUp(storage.DB().DB))

	assert.NoError(kafka.Setup(context.Background(), &wg, conf))
	assert.NoError(test.KafkaConsumer(conf))

	validator := &TestValidator{returnSubject: "user"}
	ts.api = NewCompanyAPI(validator)
}

func (ts *CompanyAPITestSuite) TestCompany() {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	fits := []*Company{
		{
			Id:           "6235fee7-d12b-4a1b-b35d-c87838b8a4a4",
			Name:         "test_name_1",
			Description:  "some description",
			Employeescnt: 1,
			Registered:   false,
			Type:         CompanyType(1),
		},
		{
			Id:           "98a75f8f-beb7-4abe-8971-bac6ab34affa",
			Name:         "test_name_2",
			Description:  "another description",
			Employeescnt: 99999,
			Registered:   false,
			Type:         CompanyType(3),
		},
	}

	doesntFit := []*Company{
		// no uuid
		{
			Name:         "the name is too long",
			Description:  "some description",
			Employeescnt: 12,
			Registered:   true,
			Type:         CompanyType(3),
		},
		// long name
		{
			Id:           "98a75f8f-beb7-4abe-8971-bac6ab34affa",
			Name:         "the name is too long",
			Description:  "some description",
			Employeescnt: 12,
			Registered:   true,
			Type:         CompanyType(3),
		},
		// no employees
		{
			Id:           "d793c371-7320-44ec-949d-d7334c27fe41",
			Name:         "name_22",
			Description:  "some description",
			Employeescnt: 0,
			Registered:   true,
			Type:         CompanyType(3),
		},
		// no type
		{
			Id:           "5658bd48-8c1b-46c9-bf96-3a20eee0bbf0",
			Name:         "name_33",
			Description:  "some description",
			Employeescnt: 66,
			Registered:   false,
		},
	}

	ts.T().Run("Create", func(t *testing.T) {

		for _, c := range fits {
			assert := require.New(t)
			_, err := ts.api.Create(
				context.Background(),
				&CreateCompanyRequest{
					Company: c,
				},
			)
			if err != nil {
				fmt.Println("create err: ", err)
			}
			assert.Nil(err)
			msg, err := test.GetMessage(fmt.Sprintf("company.%s.event.created", c.Id))
			assert.Nil(err)

			var recieved Company
			assert.NoError(json.Unmarshal(msg.Value, &recieved))
			assert.Equal(c, &recieved)

		}

		for _, c := range doesntFit {
			assert := require.New(t)
			_, err := ts.api.Create(
				context.Background(),
				&CreateCompanyRequest{
					Company: c,
				},
			)
			if err == nil {
				fmt.Println("expected err, nil returned ", c)
			}
			assert.NotNil(err)
		}

		c := fits[0]

		t.Run("Update", func(t *testing.T) {
			assert := require.New(t)

			c.Name = "name_changed"
			c.Description = "description has changed also"
			c.Employeescnt = 1
			c.Type = 2
			_, err := ts.api.Update(
				context.Background(),
				&UpdateCompanyRequest{
					Company: c,
				},
			)
			assert.Nil(err)
			msg, err := test.GetMessage(fmt.Sprintf("company.%s.event.updated", c.Id))
			assert.Nil(err)

			var recieved Company
			assert.NoError(json.Unmarshal(msg.Value, &recieved))
			assert.Equal(c, &recieved)

			r := GetCompanyRequest{
				Id: c.Id,
			}

			getResp, err := ts.api.Get(context.Background(), &r)
			assert.Nil(err)
			assert.Equal(getResp.Company, c)
			// fmt.Println("got:", getResp.Company)
			// fmt.Println("changed:", &c)
		})

		t.Run("Delete", func(t *testing.T) {
			assert := require.New(t)

			for _, c := range fits {
				r := DeleteCompanyRequest{
					Id: c.Id,
				}
				_, err := ts.api.Delete(context.Background(), &r)
				assert.Nil(err)
				msg, err := test.GetMessage(fmt.Sprintf("company.%s.event.deleted", c.Id))
				assert.Nil(err)
				assert.Nil(msg.Value)

				gr := GetCompanyRequest{
					Id: c.Id,
				}
				_, err = ts.api.Get(context.Background(), &gr)
				assert.NotNil(err)
			}
		})
	})
}
