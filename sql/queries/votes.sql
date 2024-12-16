-- name: FindVote :many
SELECT
WG_Surveys.Id AS SurveyId,
WG_Surveys.ChannelID AS SurveyChannel,
WG_Surveys.AnnouncementMessageID,
WG_Surveys.ItemName,
WG_Surveys.OpenedAt,
WG_Surveys.Deadline,
WG_Surveys.Status AS SurveyStatus,
WG_Votes.Status as VoteStatus,
WG_Votes.Token as Token,
WG_Votes.MessageId AS VoteMessage,
WG_Votes.UserId,
WG_Votes.VotedAt
FROM WG_Votes
INNER JOIN WG_Surveys ON WG_Surveys.Id = WG_Votes.SurveyId
WHERE (sqlc.arg(status) = 0 OR WG_Votes.Status = sqlc.arg(status)) 
AND (sqlc.arg(SurveyId) = "" OR WG_Votes.SurveyId LIKE sqlc.arg(SurveyOpt))
AND (sqlc.arg(UserId) = "" OR WG_Votes.UserId LIKE sqlc.arg(IdOpt))
AND (sqlc.arg(Token) = "" OR WG_Votes.Token LIKE sqlc.arg(TokenOpt))
LIMIT ? OFFSET ?;

-- name: CreateVote :exec
INSERT INTO WG_Votes (SurveyId, UserId, Token, MessageId, Status, VotedAt) VALUES (?,?,?,?,?,?);

-- name: CreateVoteEntry :exec
INSERT INTO WG_VoteEntries (SurveyId, UserId, StatId, Value) VALUES (?,?,?,?);

-- name: UpdateVote :exec
UPDATE WG_Votes SET MessageId = ?, Status = ? WHERE SurveyId = ? AND UserId = ?;

-- name: FindVoteEntries :many
SELECT
WG_Surveys.Id AS SurveyId,
WG_Surveys.ChannelID AS SurveyChannel,
WG_Surveys.AnnouncementMessageID,
WG_Surveys.ItemName,
WG_Surveys.OpenedAt,
WG_Surveys.Deadline,
WG_Surveys.Status AS SurveyStatus,
WG_VoteEntries.UserId,
WG_VoteEntries.StatId,
WG_VoteEntries.Value
FROM WG_VoteEntries
INNER JOIN WG_Votes ON WG_Votes.UserId = WG_VoteEntries.UserId AND WG_Votes.SurveyId = WG_VoteEntries.SurveyId
INNER JOIN WG_Surveys ON WG_Surveys.Id = WG_Votes.SurveyId
WHERE (sqlc.arg(status) = 0 OR WG_Votes.Status = sqlc.arg(status)) 
AND (sqlc.arg(SurveyId) = "" OR WG_Votes.SurveyId LIKE sqlc.arg(SurveyOpt))
AND (sqlc.arg(UserId) = "" OR WG_Votes.UserId LIKE sqlc.arg(IdOpt))
AND (sqlc.arg(Token) = "" OR WG_Votes.Token LIKE sqlc.arg(TokenOpt))
LIMIT ? OFFSET ?;

-- name: SumTotalVotes :one
SELECT COUNT(UserId) FROM WG_Votes WHERE SurveyId = ? AND Status != 0;

-- name: SumStatEntries :one
SELECT CAST(SUM(value) AS DECIMAL(10,3)) FROM WG_VoteEntries
INNER JOIN WG_Votes ON WG_Votes.SurveyId = WG_VoteEntries.SurveyId AND WG_Votes.UserId = WG_VoteEntries.UserId
WHERE WG_Votes.SurveyId = ? AND WG_VoteEntries.StatId = ? AND WG_Votes.Status = ?;

-- name: FindResult :many
SELECT StatId, CAST(AVG(Value) AS DECIMAL(10,4)) AS AverageValue, COUNT(*) AS TotalVotes
FROM WG_VoteEntries
INNER JOIN WG_Votes ON WG_Votes.UserId = WG_VoteEntries.UserId AND WG_VoteEntries.SurveyId = WG_VoteEntries.SurveyId
WHERE WG_VoteEntries.SurveyId = ? AND WG_Votes.SurveyId = ? AND WG_Votes.Status = ?
GROUP BY StatId;

-- name: DeleteVote :exec
DELETE FROM WG_Votes WHERE SurveyId = ? AND UserId = ?;

-- name: AlreadyVoted :one
SELECT SurveyId FROM WG_VoteEntries WHERE SurveyId = ? AND UserId = ?;

-- name: HasOpenVote :one
SELECT SurveyId, Token FROM WG_Votes WHERE SurveyId = ? AND UserId = ?;

-- name: IsContabilized :one
SELECT * FROM WG_Votes WHERE UserId = ? AND SurveyId = ? AND Status = 1;