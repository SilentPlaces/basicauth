package registeration_dto

// SignUpRequestDTO request dto client sends tos tart signup process
type SignUpRequestDTO struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
