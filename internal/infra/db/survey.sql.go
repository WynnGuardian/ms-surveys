// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: survey.sql

package db

import (
	"context"
	"time"
)

const createSurvey = `-- name: CreateSurvey :exec
INSERT INTO WG_Surveys (Id, Status, ItemName, OpenedAt, Deadline, ChannelId, AnnouncementMessageID) VALUES(?,?,?,?,?,?,?)
`

type CreateSurveyParams struct {
	ID                    string    `json:"id"`
	Status                int8      `json:"status"`
	Itemname              string    `json:"itemname"`
	Openedat              time.Time `json:"openedat"`
	Deadline              time.Time `json:"deadline"`
	Channelid             string    `json:"channelid"`
	Announcementmessageid string    `json:"announcementmessageid"`
}

func (q *Queries) CreateSurvey(ctx context.Context, arg CreateSurveyParams) error {
	_, err := q.db.ExecContext(ctx, createSurvey,
		arg.ID,
		arg.Status,
		arg.Itemname,
		arg.Openedat,
		arg.Deadline,
		arg.Channelid,
		arg.Announcementmessageid,
	)
	return err
}

const deleteSurvey = `-- name: DeleteSurvey :exec
DELETE FROM WG_Surveys WHERE Id = ?
`

func (q *Queries) DeleteSurvey(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteSurvey, id)
	return err
}

const deleteSurveyEntries = `-- name: DeleteSurveyEntries :exec
DELETE FROM WG_VoteEntries WHERE SurveyId = ?
`

func (q *Queries) DeleteSurveyEntries(ctx context.Context, surveyid string) error {
	_, err := q.db.ExecContext(ctx, deleteSurveyEntries, surveyid)
	return err
}

const deleteVotes = `-- name: DeleteVotes :exec
DELETE FROM WG_Votes WHERE SurveyId = ?
`

func (q *Queries) DeleteVotes(ctx context.Context, surveyid string) error {
	_, err := q.db.ExecContext(ctx, deleteVotes, surveyid)
	return err
}

const findExpiredSurveys = `-- name: FindExpiredSurveys :many
SELECT id, channelid, announcementmessageid, status, itemname, openedat, deadline FROM WG_Surveys WHERE Deadline < NOW() && Status = ?
`

func (q *Queries) FindExpiredSurveys(ctx context.Context, status int8) ([]WgSurvey, error) {
	rows, err := q.db.QueryContext(ctx, findExpiredSurveys, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WgSurvey
	for rows.Next() {
		var i WgSurvey
		if err := rows.Scan(
			&i.ID,
			&i.Channelid,
			&i.Announcementmessageid,
			&i.Status,
			&i.Itemname,
			&i.Openedat,
			&i.Deadline,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findSurvey = `-- name: FindSurvey :many
SELECT id, channelid, announcementmessageid, status, itemname, openedat, deadline FROM WG_Surveys WHERE 
(? = 0 OR Status = ?) 
AND (? = "" OR ItemName LIKE ?)
AND (? = "" OR Id LIKE ?)
LIMIT ? OFFSET ?
`

type FindSurveyParams struct {
	Status      int8        `json:"status"`
	Itemname    interface{} `json:"itemname"`
	Itemnameopt string      `json:"itemnameopt"`
	ID          interface{} `json:"id"`
	Idopt       string      `json:"idopt"`
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
}

func (q *Queries) FindSurvey(ctx context.Context, arg FindSurveyParams) ([]WgSurvey, error) {
	rows, err := q.db.QueryContext(ctx, findSurvey,
		arg.Status,
		arg.Status,
		arg.Itemname,
		arg.Itemnameopt,
		arg.ID,
		arg.Idopt,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WgSurvey
	for rows.Next() {
		var i WgSurvey
		if err := rows.Scan(
			&i.ID,
			&i.Channelid,
			&i.Announcementmessageid,
			&i.Status,
			&i.Itemname,
			&i.Openedat,
			&i.Deadline,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSurvey = `-- name: UpdateSurvey :exec
UPDATE WG_Surveys SET Status = ?, ItemName = ?, OpenedAt = ?, Deadline = ?, ChannelID = ?, AnnouncementMessageID = ? WHERE Id = ?
`

type UpdateSurveyParams struct {
	Status                int8      `json:"status"`
	Itemname              string    `json:"itemname"`
	Openedat              time.Time `json:"openedat"`
	Deadline              time.Time `json:"deadline"`
	Channelid             string    `json:"channelid"`
	Announcementmessageid string    `json:"announcementmessageid"`
	ID                    string    `json:"id"`
}

func (q *Queries) UpdateSurvey(ctx context.Context, arg UpdateSurveyParams) error {
	_, err := q.db.ExecContext(ctx, updateSurvey,
		arg.Status,
		arg.Itemname,
		arg.Openedat,
		arg.Deadline,
		arg.Channelid,
		arg.Announcementmessageid,
		arg.ID,
	)
	return err
}
