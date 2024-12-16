package entity

import (
	"time"
	"victo/wynnguardian/internal/infra/enums"
)

type SurveyVote struct {
	Survey        *Survey            `json:"survey"`
	DiscordUserID string             `json:"user_id"`
	MessageID     string             `json:"message_id"`
	ChannelID     string             `json:"channel_id"`
	Token         string             `json:"token"`
	Votes         map[string]float64 `json:"votes"`
	VotedAt       time.Time          `json:"voted_at"`
	Status        enums.VoteStatus   `json:"status"`
}

type SurveyVoteEntry struct {
	Survey *Survey `json:"survey_id"`
	UserID string  `json:"user_dc_id"`
	Stat   string  `json:"stat"`
	Value  float64 `json:"value"`
}
