package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

// Company

// Company stuct represents a company model
type Company struct {
	ID           uuid.UUID `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	EmployeesCnt int32     `db:"employees_cnt"`
	Registered   bool      `db:"registered"`
	Type         uint32    `db:"type"`
}

// CreateCompany creates the given Company in db.
func CreateCompany(ctx context.Context, db sqlx.Execer, c *Company) error {
	now := time.Now()

	_, err := db.Exec(`
		INSERT INTO company (
			created_at,
			updated_at,
			id,
			name,
			description,
			employees_cnt,
			registered,
			type
		) values ($1, $2, $3, $4, $5, $6, $7, $8)
		`,
		now,
		now,
		c.ID,
		c.Name,
		c.Description,
		c.EmployeesCnt,
		c.Registered,
		c.Type,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"name": c.Name,
	}).Info("company created")
	return nil
}

// UpdateCompany updates the given company by its ID.
func UpdateCompany(ctx context.Context, db sqlx.Ext, c *Company) error {
	res, err := db.Exec(`
		UPDATE company
		SET
			updated_at = $2,
			name = $3,
			description = $4,
			employees_cnt = $5,
			registered = $6,
			type = $7	
		WHERE
			id = $1`,
		c.ID,
		time.Now(),
		c.Name,
		c.Description,
		c.EmployeesCnt,
		c.Registered,
		c.Type,
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
func DeleteCompany(ctx context.Context, db sqlx.Ext, id uuid.UUID) error {
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
func GetCompany(ctx context.Context, db sqlx.Ext, id uuid.UUID) (Company, error) {
	var result Company

	err := sqlx.Get(db, &result, "SELECT * FROM company WHERE id = $1", id)
	if err != nil {
		return result, handlePSQLError(Select, err, "can't select Company")
	}
	return result, nil
}
