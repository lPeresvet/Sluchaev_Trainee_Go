package core

type FailureMessage struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
