package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

// Company

// Company stuct represents a company model
type Company struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
	Code      string    `db:"code"`
	Country   string    `db:"country"`
	Website   string    `db:"website"`
	Phone     string    `db:"phone"`
}

// CompanyFilters stuct sql filter
type CompanyFilters struct {
	Name    string `db:"name"`
	Code    string `db:"code"`
	Country string `db:"country"`
	Website string `db:"website"`
	Phone   string `db:"phone"`
	Limit   int64  `db:"limit"`
	Offset  int64  `db:"offset"`
}

// SQL returns the SQL filter.
func (f CompanyFilters) SQL() string {
	var filters []string

	if f.Name != "" {
		filters = append(filters, "c.name = :name")
	}

	if f.Code != "" {
		filters = append(filters, "c.code = :code")
	}

	if f.Country != "" {
		filters = append(filters, "c.country = :country")
	}

	if f.Website != "" {
		filters = append(filters, "c.website = :website")
	}

	if f.Phone != "" {
		filters = append(filters, "c.phone = :phone")
	}

	if len(filters) == 0 {
		return ""
	}

	return "WHERE " + strings.Join(filters, " AND ")
}

// CreateCompany creates the given Company in db.
func CreateCompany(ctx context.Context, db sqlx.Execer, c *Company) error {
	now := time.Now()

	_, err := db.Exec(`
		INSERT INTO company (
			created_at,
			updated_at,
			name,
			code,
			country,
			website,
			phone
		) values ($1, $2, $3, $4, $5, $6, $7)
		`,
		now,
		now,
		c.Name,
		c.Code,
		c.Country,
		c.Website,
		c.Phone,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"name": c.Name,
	}).Info("company created")
	return nil
}

// GetCompanies returns a slice of companies according to the filter.
func GetCompanies(ctx context.Context,
	db sqlx.Queryer, filters CompanyFilters) ([]Company, error) {

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		SELECT 
			c.*
		FROM
			company c	
		`+filters.SQL()+`
		limit :limit
		offset :offset
	`, filters)
	if err != nil {
		return nil, fmt.Errorf("unable to make BindNamed %v", err)
	}

	var result []Company
	err = sqlx.Select(db, &result, query, args...)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return result, nil
}

// UpdateCompany updates the given company by its ID.
func UpdateCompany(ctx context.Context, db sqlx.Ext, c Company) error {
	res, err := db.Exec(`
		UPDATE company
		SET
			updated_at = $2,
			name = $3,
			code = $4,
			country = $5,
			website = $6,
			phone = $7	
		WHERE
			id = $1`,
		c.ID,
		time.Now(),
		c.Name,
		c.Code,
		c.Country,
		c.Website,
		c.Phone,
	)
	if err != nil {
		return handlePSQLError(Update, err, "can't update")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows affected %v", err)
	}
	if ra == 0 {
		return ErrDoesNotExist
	}
	return nil
}

// DeleteCompany deletes a company that matches the given ID.
func DeleteCompany(ctx context.Context, db sqlx.Ext, id int64) error {
	res, err := db.Exec("DELETE FROM company WHERE id = $1", id)
	if err != nil {
		return handlePSQLError(Delete, err, "can't delete")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows affected %v", err)
	}
	if ra == 0 {
		return ErrDoesNotExist
	}
	return nil
}

// GetCompany gets a company that matches the given ID.
func GetCompany(ctx context.Context, db sqlx.Ext, id int64) (Company, error) {
	var result Company

	err := sqlx.Get(db, &result, "SELECT * FROM company WHERE id = $1", id)
	if err != nil {
		return result, handlePSQLError(Select, err, "can't select Company")
	}
	return result, nil
}
