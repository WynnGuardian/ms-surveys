-- name: CreateSurvey :exec
INSERT INTO WG_Surveys (Id, Status, ItemName, OpenedAt, Deadline, ChannelId, AnnouncementMessageID) VALUES(?,?,?,?,?,?,?);

-- name: UpdateSurvey :exec
UPDATE WG_Surveys SET Status = ?, ItemName = ?, OpenedAt = ?, Deadline = ?, ChannelID = ?, AnnouncementMessageID = ? WHERE Id = ?;

-- name: FindSurvey :many
SELECT * FROM WG_Surveys WHERE 
(sqlc.arg(status) = 0 OR Status = sqlc.arg(status)) 
AND (sqlc.arg(ItemName) = "" OR ItemName LIKE sqlc.arg(ItemNameOpt))
AND (sqlc.arg(Id) = "" OR Id LIKE sqlc.arg(IdOpt))
LIMIT ? OFFSET ?;

-- name: FindExpiredSurveys :many
SELECT * FROM WG_Surveys WHERE Deadline < NOW() && Status = ?;

-- name: DeleteSurveyEntries :exec
DELETE FROM WG_VoteEntries WHERE SurveyId = ?;

-- name: DeleteVotes :exec
DELETE FROM WG_Votes WHERE SurveyId = ?;

-- name: DeleteSurvey :exec
DELETE FROM WG_Surveys WHERE Id = ?;