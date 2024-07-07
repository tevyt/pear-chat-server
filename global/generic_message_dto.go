package global

type GenericMessageDTO struct {
	Message string `json:"message"`
}

func NewGenericeMessageDTO(message string) *GenericMessageDTO {
	m := GenericMessageDTO{Message: message}
	return &m
}
