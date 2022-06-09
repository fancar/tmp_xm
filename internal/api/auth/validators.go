package auth

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// API key subjects.
const (
	SubjectUser   = "user"
	SubjectAPIKey = "api_key"
)

// Flag defines the authorization flag.
type Flag int

// Authorization flags.
const (
	Create Flag = iota
	Read
	Update
	Delete
	List
	UpdateProfile
	ADRAlgorithms
)

func (f Flag) String() string {
	return [...]string{"Create", "Read", "Update", "Delete", "List", "UpdateProfile", "ADRAlgorithms"}[f]
}

// ValidateActiveUser validates if the user in the JWT claim is active.
func ValidateActiveUser() ValidatorFunc {
	query := `
		select
			1
		from
			"user" u
	`

	where := [][]string{
		{"(u.username = $1 or u.id = $2)"},
	}

	return func(db sqlx.Queryer, claims *Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return executeQuery(db, query, where, claims.Username, claims.UserID)
		case SubjectAPIKey:
			return false, nil
		default:
			return false, nil
		}
	}
}

func executeQuery(db sqlx.Queryer, query string, where [][]string, args ...interface{}) (bool, error) {
	var ors []string
	for _, ands := range where {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	query = "select count(*) from (" + query + " where " + whereStr + " limit 1) count_only"

	var count int64

	if err := sqlx.Get(db, &count, query, args...); err != nil {
		return false, fmt.Errorf("validator select error %v", err)
	}
	return count > 0, nil
}
