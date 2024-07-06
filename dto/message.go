package dto

type GenericMessage struct {
	Message string `json:"message"`
}

func NewGenericeMessage(message string) *GenericMessage {
	m := GenericMessage{Message: message}
	return &m
}
