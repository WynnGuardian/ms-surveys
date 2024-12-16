package response

type APIResponse[E any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Body    E      `json:"body"`
}

func Of[E any](body E, status int, message string) APIResponse[E] {
	return APIResponse[E]{
		Message: message,
		Body:    body,
		Status:  status,
	}
}
