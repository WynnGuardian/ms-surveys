// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: bans.sql

package db

import (
	"context"
)

const deleteBan = `-- name: DeleteBan :exec
DELETE FROM WG_SurveyBan WHERE UserID = ?
`

func (q *Queries) DeleteBan(ctx context.Context, userid string) error {
	_, err := q.db.ExecContext(ctx, deleteBan, userid)
	return err
}

const findBan = `-- name: FindBan :one
SELECT Reason FROM WG_SurveyBan WHERE UserID = ?
`

func (q *Queries) FindBan(ctx context.Context, userid string) (string, error) {
	row := q.db.QueryRowContext(ctx, findBan, userid)
	var reason string
	err := row.Scan(&reason)
	return reason, err
}

const insertBan = `-- name: InsertBan :exec
INSERT INTO WG_SurveyBan (UserID, Reason) VALUES (?,?)
`

type InsertBanParams struct {
	Userid string `json:"userid"`
	Reason string `json:"reason"`
}

func (q *Queries) InsertBan(ctx context.Context, arg InsertBanParams) error {
	_, err := q.db.ExecContext(ctx, insertBan, arg.Userid, arg.Reason)
	return err
}
