package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wynnguardian/ms-surveys/internal/config"
)

type apiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func newCall(endpoint string) string {
	return fmt.Sprintf("%s/%s", config.MainConfig.Private.Tokens.Self, endpoint)
}

func defaultHeaders(req *http.Request) {
	req.Header.Add("Authorization", config.MainConfig.Private.Tokens.Self)
}

func post[T any](path string, body *T) (*apiResponse, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest("POST", newCall(path), bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}

	defaultHeaders(r)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resp := apiResponse{}
	derr := json.NewDecoder(res.Body).Decode(&resp)
	if derr != nil {
		return nil, derr
	}

	return &resp, nil
}
