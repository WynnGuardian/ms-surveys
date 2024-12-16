package response

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

var (
	Ok = WGResponse{
		Status:  http.StatusOK,
		Message: "",
	}
	ErrInvalidItem = WGResponse{
		Status:  http.StatusNotFound,
		Message: "invalid item.",
	}
	ErrSurveyNotOpen = WGResponse{
		Status:  http.StatusForbidden,
		Message: "survey not open.",
	}
	ErrVoteInterval = WGResponse{
		Status:  http.StatusForbidden,
		Message: "vote sum must be between 99 and 101.",
	}
	ErrCriteriaMissing = WGResponse{
		Status:  http.StatusBadRequest,
		Message: "criteria is missing.",
	}
	ErrAlreadyVoted = WGResponse{
		Status:  http.StatusForbidden,
		Message: "you already voted in this survey.",
	}
	ErrNotWaitingApproval = WGResponse{
		Status:  http.StatusForbidden,
		Message: "survey is not waiting for approval.",
	}
	ErrWynnItemNotFound = WGResponse{
		Status:  http.StatusNotFound,
		Message: "WynnItem not found.",
	}
	ErrVoteNotFound = WGResponse{
		Status:  http.StatusNotFound,
		Message: "vote not found.",
	}
	ErrSurveyNotFound = WGResponse{
		Status:  http.StatusNotFound,
		Message: "Survey not found.",
	}
	ErrCriteriaNotFound = WGResponse{
		Status:  http.StatusNotFound,
		Message: "Criterias not found.",
	}
	ErrBadRequest = WGResponse{
		Status:  http.StatusBadRequest,
		Message: "An error has occured. Please try again soon.",
	}
	ErrUnauthorized = WGResponse{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized action. Please try again soon.",
	}
	ErrTryAgain = WGResponse{
		Status:  http.StatusInternalServerError,
		Message: "An error has occured. Please try again soon.",
	}
	ErrVoteAlreadyConfirmed = WGResponse{
		Status:  http.StatusForbidden,
		Message: "This vote was already confirmed.",
	}
)

func New[T any](status int, message string, body T) WGResponse {
	js, _ := json.Marshal(body)
	return WGResponse{
		Status:  status,
		Message: message,
		Body:    string(js),
	}
}

func ErrInternalServerErr(err error) WGResponse {
	log.Println(err.Error())
	debug.PrintStack()
	return ErrTryAgain
}
