package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wynnguardian/common/handlerfunc"
	middleware "github.com/wynnguardian/common/middlewares"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/ms-surveys/internal/config"
	"github.com/wynnguardian/ms-surveys/internal/infra/http/handlers"
)

type RouterEntry struct {
	MustBeMod bool
	Handler   handlerfunc.HandlerFunc
	Path      string
	Method    string
}

var (
	entries = []RouterEntry{
		{Path: "/surveyCreate", MustBeMod: true, Method: "POST", Handler: handlers.OpenSurvey},
		{Path: "/findOpenSurvey", MustBeMod: false, Method: "POST", Handler: handlers.FindSurvey},
		{Path: "/sendVote", MustBeMod: false, Method: "POST", Handler: handlers.SendVote},
		{Path: "/createVote", MustBeMod: false, Method: "POST", Handler: handlers.CreateVote},
		{Path: "/defineSurveyInfo", MustBeMod: false, Method: "POST", Handler: handlers.DefineSurveyInfo},
		{Path: "/defineVoteMessage", MustBeMod: true, Method: "POST", Handler: handlers.DefineVoteChannelMessage},
		{Path: "/confirmVote", MustBeMod: true, Method: "POST", Handler: handlers.ConfirmVote},
		{Path: "/closeSurvey", MustBeMod: true, Method: "POST", Handler: handlers.CloseSurvey},
		{Path: "/cancelSurvey", MustBeMod: true, Method: "POST", Handler: handlers.CancelSurvey},
		{Path: "/approveSurvey", MustBeMod: true, Method: "POST", Handler: handlers.ApproveSurvey},
		{Path: "/discardSurvey", MustBeMod: true, Method: "POST", Handler: handlers.DiscardSurvey},
	}
)

func post(engine *gin.Engine, path string, handler handlerfunc.HandlerFunc) {
	engine.POST(path, middleware.Parse(handler))
}

func Build() *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "http://proxy:8088,http://bot:8087")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	for _, entry := range entries {
		switch entry.Method {
		case "POST":
			post(engine, entry.Path, Authorization(entry.Handler))
		}
	}
	return engine
}

func Authorization(hf handlerfunc.HandlerFunc) handlerfunc.HandlerFunc {
	return func(ctx *gin.Context) response.WGResponse {
		whitelist := config.MainConfig.Private.Tokens.Whitelist
		for _, w := range whitelist {
			if w == ctx.GetHeader("Authorization") {
				fmt.Println("passoy")
				return hf(ctx)
			}
		}
		return response.ErrUnauthorized
	}
}
