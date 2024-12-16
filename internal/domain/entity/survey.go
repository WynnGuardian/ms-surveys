package entity

import (
	"time"
	"victo/wynnguardian/internal/infra/enums"
)

type Survey struct {
	ID                    string             `json:"id"`
	ChannelID             string             `json:"channel_id"`
	AnnouncementMessageID string             `json:"announcement_message_id"`
	ItemName              string             `json:"item_name"`
	OpenedAt              time.Time          `json:"opened_at"`
	Deadline              time.Time          `json:"deadline"`
	Status                enums.SurveyStatus `json:"status"`
}

type SurveyResult struct {
	SurveyID   string             `json:"survey_id"`
	ItemName   string             `json:"item_name"`
	TotalVotes int                `json:"total_votes"`
	Results    map[string]float64 `json:"results"`
}
