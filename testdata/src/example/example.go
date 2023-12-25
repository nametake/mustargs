package example

import (
	"context"
	"database/sql"
)

type TenantID string
type UserID string

type Usecase struct{}

func (u *Usecase) GetUser(ctx context.Context, tenantID TenantID, userID string) {
}

type DB struct{}

func (db *DB) GetUser(ctx context.Context, tx *sql.Tx, tenantID TenantID, userID UserID) {}

func (db *DB) GetUsers(ctx context.Context, tx *sql.Tx, tenantID TenantID, limit, offset int) {
}
