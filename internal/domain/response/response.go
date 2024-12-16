package response

import "net/http"

type WGResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Body    string `json:"body"`
}

func (w *WGResponse) Ok() bool {
	return w.Status == http.StatusOK
}
