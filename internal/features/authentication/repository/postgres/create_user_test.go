package authentication_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

type poolStub struct {
	row core_postgres_pool.Row
}

func (p poolStub) Query(
	context.Context,
	string,
	...any,
) (core_postgres_pool.Rows, error) {
	return nil, errors.New("unexpected Query call")
}

func (p poolStub) QueryRow(context.Context, string, ...any) core_postgres_pool.Row {
	return p.row
}

func (p poolStub) Exec(
	context.Context,
	string,
	...any,
) (core_postgres_pool.CommandTag, error) {
	return nil, errors.New("unexpected Exec call")
}

func (poolStub) Close() {}

func (poolStub) OpTimeout() time.Duration {
	return time.Second
}

type rowStub struct {
	err error
}

func (r rowStub) Scan(...any) error {
	return r.err
}

func TestCreateUserMapsUniqueViolationToConflict(t *testing.T) {
	repository := NewAuthenticationRepository(poolStub{
		row: rowStub{err: core_postgres_pool.ErrViolatesUnique},
	})

	_, err := repository.CreateUser(
		context.Background(),
		"ivan@example.com",
		"password-hash",
		"Ivan Ivanov",
	)
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("CreateUser() error = %v, want ErrConflict", err)
	}
}
