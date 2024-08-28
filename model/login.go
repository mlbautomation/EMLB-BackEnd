package model

// Con esta información en el body solicitaremos el token
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
