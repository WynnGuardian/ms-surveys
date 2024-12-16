package enums

type SurveyStatus int8

const (
	SURVEY_WAITING_APPROVAL SurveyStatus = iota + 1
	SURVEY_DENIED
	SURVEY_APPROVED
	SURVEY_OPEN
)

type VoteStatus int8

const (
	VOTE_NOT_CONFIRMED VoteStatus = iota + 1
	VOTE_CONTABILIZED
	VOTE_DENIED
)
