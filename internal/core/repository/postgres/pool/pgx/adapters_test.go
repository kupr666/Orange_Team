package core_pgx_pool

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func TestMapErrorsMapsUniqueViolation(t *testing.T) {
	err := mapErrors(&pgconn.PgError{Code: "23505"})

	if !errors.Is(err, core_postgres_pool.ErrViolatesUnique) {
		t.Fatalf("mapErrors() error = %v, want ErrViolatesUnique", err)
	}
}

func TestMapErrorsMapsForeignKeyViolation(t *testing.T) {
	err := mapErrors(&pgconn.PgError{Code: "23503"})

	if !errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
		t.Fatalf("mapErrors() error = %v, want ErrViolatesForeignKey", err)
	}
}
