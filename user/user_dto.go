package user

type UserDTO struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
	PublicKey    string `json:"publicKey"`
}
