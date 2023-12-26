package example

import (
	"context"
	"database/sql"
)

type TenantID string

type Usecase struct{}

func (u *Usecase) GetUser(ctx context.Context, tenantID TenantID, userID string) {
}

func (u *Usecase) GetPost(ctx context.Context, userID string) { // want "no TenantID type arg at index 1 found for func GetPost"
}

type DB struct{}

func (db *DB) GetUser(ctx context.Context, tx *sql.Tx, tenantID TenantID, userID string) {
}

func (db *DB) GetPost(ctx context.Context, tenantID TenantID, postID string) { // want "no \\*sql.Tx type arg at index 1 found for func GetPost"
}

func (db *DB) GetMultipleUsers(ctx context.Context, tx *sql.Tx, tenantID TenantID, limit, offset int) {
}

func (db *DB) GetMultiplePosts(ctx context.Context, tx *sql.Tx, tenantID TenantID) { // want "no int type arg at index -1, no int type arg at index -2 found for func GetMultiplePosts"
}
