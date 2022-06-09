package api

import (
	"context"

	log "github.com/sirupsen/logrus"

	// "github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/fancar/tmp_xm/internal/api/auth"
	"github.com/fancar/tmp_xm/internal/api/helpers"
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

	if err := helpers.IPaddrCheker(ctx); err != nil {
		return &empty.Empty{}, err
	}

	item := storage.Company{
		Name:    req.Company.Name,
		Code:    req.Company.Code,
		Country: req.Company.Country,
		Website: req.Company.Website,
		Phone:   req.Company.Phone,
	}

	err := storage.CreateCompany(ctx, storage.DB(), &item)
	if err != nil {
		return &empty.Empty{}, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Get returns an item
func (a *CompanyAPI) Get(ctx context.Context, req *GetCompanyRequest) (*GetCompanyResponse, error) {
	log.Debug("api/Get request:", req)

	result := &GetCompanyResponse{}
	d, err := storage.GetCompany(ctx, storage.DB(), req.Id)
	if err != nil {
		return result, helpers.ErrToRPCError(err)
	}

	return &GetCompanyResponse{
		Company: &Company{
			Id:      d.ID,
			Name:    d.Name,
			Code:    d.Code,
			Country: d.Country,
			Website: d.Website,
			Phone:   d.Phone,
		},
	}, nil
}

// List items
func (a *CompanyAPI) List(ctx context.Context, req *ListCompanyRequest) (*ListCompanyResponse, error) {
	log.Debug("api/List request:", req)

	filters := storage.CompanyFilters{
		Name:    req.Name,
		Code:    req.Code,
		Country: req.Country,
		Website: req.Website,
		Phone:   req.Phone,
		Offset:  req.Offset,
		Limit:   req.Limit,
	}

	if filters.Limit == 0 {
		filters.Limit = 10000
	}

	companies, err := storage.GetCompanies(ctx, storage.DB(), filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	result := &ListCompanyResponse{}

	for _, c := range companies {
		result.Result = append(result.Result, &Company{
			Id:      c.ID,
			Name:    c.Name,
			Code:    c.Code,
			Country: c.Country,
			Website: c.Website,
			Phone:   c.Phone,
		})
	}

	return result, nil
}

// Update the item
func (a *CompanyAPI) Update(ctx context.Context, req *UpdateCompanyRequest) (*empty.Empty, error) {
	log.Debug("api/Update request:", req)

	item := storage.Company{
		ID:      req.Company.Id,
		Name:    req.Company.Name,
		Code:    req.Company.Code,
		Country: req.Company.Country,
		Website: req.Company.Website,
		Phone:   req.Company.Phone,
	}

	err := storage.UpdateCompany(ctx, storage.DB(), item)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Delete the item
func (a *CompanyAPI) Delete(ctx context.Context, req *DeleteCompanyRequest) (*empty.Empty, error) {
	log.Debug("api/Delete request:", req)

	if err := a.validator.Validate(ctx, auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := helpers.IPaddrCheker(ctx); err != nil {
		return &empty.Empty{}, err
	}

	err := storage.DeleteCompany(ctx, storage.DB(), req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	return &empty.Empty{}, nil
}
