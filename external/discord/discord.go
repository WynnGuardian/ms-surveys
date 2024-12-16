package discord

import (
	"errors"
	"victo/wynnguardian/internal/domain/entity"
)

func NotifySurveyVote(vote *entity.SurveyVote) error {
	resp, err := discordServer.notifyVote(vote)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return errors.New(resp.Message)
	}
	return nil
}

func NotifySurveyCreate(survey *entity.Survey) error {
	resp, err := discordServer.notifySurveyCreate(survey)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return errors.New(resp.Message)
	}
	return nil
}

func SendResultForApproval(survey *entity.SurveyResult) error {
	resp, err := discordServer.sendResultForApproval(survey)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return errors.New(resp.Message)
	}
	return nil
}
