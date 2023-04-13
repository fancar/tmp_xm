package api

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	// "github.com/golang/protobuf/ptypes"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/fancar/tmp_xm/internal/api/auth"
	"github.com/fancar/tmp_xm/internal/api/helpers"
	"github.com/fancar/tmp_xm/internal/kafka"
	"github.com/fancar/tmp_xm/internal/storage"
)

// CompanyAPI exports the internal User related functions.
type CompanyAPI struct {
	validator auth.Validator
}

// NewCompanyAPI creates a new api for companies funcs.
func NewCompanyAPI(validator auth.Validator) *CompanyAPI { // validator auth.Validator
	return &CompanyAPI{
		validator: validator,
	}
}

// Login validates the login request and returns a JWT token.
func (a *CompanyAPI) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	jwt, err := storage.LoginUserByPassword(ctx, storage.DB(), req.User, req.Password)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	return &LoginResponse{Jwt: jwt}, nil
}

// Create adds new item in storage
func (a *CompanyAPI) Create(ctx context.Context,
	req *CreateCompanyRequest) (*empty.Empty, error) {
	log.Debug("api/Create request:", req)

	if err := a.validator.Validate(ctx, auth.ValidateActiveUser()); err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	item, err := a.convertCompany(ctx, req.Company)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "check your body: %s", err)
	}

	err = storage.CreateCompany(ctx, storage.DB(), item)
	if err != nil {
		return &empty.Empty{}, helpers.ErrToRPCError(err)
	}

	go sendEvent(ctx, item, req.Company.Id, "created")

	return &empty.Empty{}, nil
}

// Get returns an item
func (a *CompanyAPI) Get(ctx context.Context, req *GetCompanyRequest) (*GetCompanyResponse, error) {
	log.Debug("api/Get request:", req)

	ID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad value: %s", err)
	}

	result := &GetCompanyResponse{}
	d, err := storage.GetCompany(ctx, storage.DB(), ID)
	if err != nil {
		return result, helpers.ErrToRPCError(err)
	}

	return &GetCompanyResponse{
		Company: &Company{
			Id:           d.ID.String(),
			Name:         d.Name,
			Description:  d.Description,
			Employeescnt: d.EmployeesCnt,
			Registered:   d.Registered,
			Type:         CompanyType(d.Type),
		},
	}, nil
}

// Update the item
func (a *CompanyAPI) Update(ctx context.Context, req *UpdateCompanyRequest) (*empty.Empty, error) {
	log.Debug("api/Update request:", req)

	if err := a.validator.Validate(ctx, auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	item, err := a.convertCompany(ctx, req.Company)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad value: %s", err)
	}

	err = storage.UpdateCompany(ctx, storage.DB(), item)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	go sendEvent(ctx, item, req.Company.Id, "updated")

	return &empty.Empty{}, nil
}

// Delete the item
func (a *CompanyAPI) Delete(ctx context.Context, req *DeleteCompanyRequest) (*empty.Empty, error) {
	log.Debug("api/Delete request:", req)

	if err := a.validator.Validate(ctx, auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	ID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad value: %s", err)
	}

	err = storage.DeleteCompany(ctx, storage.DB(), ID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	go sendEvent(ctx, nil, req.Id, "deleted")

	return &empty.Empty{}, nil
}

// convertCompany validates all the fields and converts it to local struct
func (a *CompanyAPI) convertCompany(ctx context.Context, in *Company) (*storage.Company, error) {

	ID, err := uuid.FromString(in.Id)
	if err != nil {
		return nil, err
	}

	if in.Name == "" {
		return nil, fmt.Errorf("you must specify the 'name' field")
	}

	if in.Employeescnt == 0 {
		return nil, fmt.Errorf("you must specify 'EmployeesCnt'")
	}

	if in.Type.Number() == 0 {
		return nil, fmt.Errorf("you must specify 'Type'")
	}

	result := &storage.Company{
		ID:           ID,
		Name:         in.Name,
		Description:  in.Description,
		EmployeesCnt: in.Employeescnt,
		Registered:   in.Registered,
		Type:         uint32(in.Type.Number()),
	}

	return result, nil
}

// sendEvent prepeares data and sends the event via kafka producer
func sendEvent(ctx context.Context, item *storage.Company, id, event string) {
	b := []byte{}
	var err error
	if item != nil {
		b, err = json.Marshal(item)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"company_id": id,
				"event":      event,
			}).Error("unable to marshal data")
			return
		}
	}
	kafka.PublishMessage(ctx, id, event, b)
}
