package registeration_dto

// RegistrationRequestDTO request dto client sends tos tart signup process
type RegistrationRequestDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
