package discord

import "github.com/wynnguardian/common/entity"

var discordServer discordInternalServerInterface

type discordInternalServerInterface interface {
	notifySurveyEnd(*entity.Survey) (*apiResponse, error)
	notifySurveyCreate(*entity.Survey) (*apiResponse, error)
	notifyVote(*entity.SurveyVote) (*apiResponse, error)
	sendResultForApproval(*entity.SurveyResult) (*apiResponse, error)
}

type discordInternalServer struct {
	discordInternalServerInterface
}

func SetupDiscordServer() {
	discordServer = &discordInternalServer{}
}

func (s *discordInternalServer) notifySurveyEnd(surv *entity.Survey) (*apiResponse, error) {
	return post("surveyEnd", surv)
}

func (s *discordInternalServer) sendResultForApproval(surv *entity.SurveyResult) (*apiResponse, error) {
	return post("surveyApproval", surv)
}

func (s *discordInternalServer) notifySurveyCreate(surv *entity.Survey) (*apiResponse, error) {
	return post("surveyCreate", surv)
}

func (s *discordInternalServer) notifyVote(v *entity.SurveyVote) (*apiResponse, error) {
	return post("vote", v)
}
