// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"
)

type WgAuthenticateditem struct {
	ID           string    `json:"id"`
	Lastranked   time.Time `json:"lastranked"`
	Itemname     string    `json:"itemname"`
	Ownermcuuid  string    `json:"ownermcuuid"`
	Owneruserid  string    `json:"owneruserid"`
	Position     int32     `json:"position"`
	Trackingcode string    `json:"trackingcode"`
	Ownerpublic  int32     `json:"ownerpublic"`
	Bytes        string    `json:"bytes"`
}

type WgAuthenticateditemstat struct {
	Itemid string `json:"itemid"`
	Statid string `json:"statid"`
	Value  int32  `json:"value"`
}

type WgCriterium struct {
	Itemname string  `json:"itemname"`
	Statid   string  `json:"statid"`
	Value    float64 `json:"value"`
}

type WgSurvey struct {
	ID                    string    `json:"id"`
	Channelid             string    `json:"channelid"`
	Announcementmessageid string    `json:"announcementmessageid"`
	Status                int8      `json:"status"`
	Itemname              string    `json:"itemname"`
	Openedat              time.Time `json:"openedat"`
	Deadline              time.Time `json:"deadline"`
}

type WgVote struct {
	Messageid string    `json:"messageid"`
	Userid    string    `json:"userid"`
	Surveyid  string    `json:"surveyid"`
	Token     string    `json:"token"`
	Status    int8      `json:"status"`
	Votedat   time.Time `json:"votedat"`
}

type WgVoteentry struct {
	Surveyid string  `json:"surveyid"`
	Userid   string  `json:"userid"`
	Statid   string  `json:"statid"`
	Value    float64 `json:"value"`
}

type WgWynnitem struct {
	Name            string `json:"name"`
	Sprite          string `json:"sprite"`
	Reqlevel        int32  `json:"reqlevel"`
	Reqstrenght     int32  `json:"reqstrenght"`
	Reqagility      int32  `json:"reqagility"`
	Reqdefence      int32  `json:"reqdefence"`
	Reqintelligence int32  `json:"reqintelligence"`
	Reqdexterity    int32  `json:"reqdexterity"`
}

type WgWynnitemstat struct {
	Itemname string `json:"itemname"`
	Statid   string `json:"statid"`
	Lower    int32  `json:"lower"`
	Upper    int32  `json:"upper"`
}