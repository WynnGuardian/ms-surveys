-- name: InsertBan :exec
INSERT INTO WG_SurveyBan (UserID, Reason) VALUES (?,?);

-- name: DeleteBan :exec
DELETE FROM WG_SurveyBan WHERE UserID = ?;

-- name: FindBan :one
SELECT Reason FROM WG_SurveyBan WHERE UserID = ?;