package dto

type LoginSuccess struct {
	EmailAddress string `json:"emailAddress"`
	Name         string `json:"name"`
	SessionID    string `json:"sessionID"`
}
