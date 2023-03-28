package login

type LoginDTO struct {
	Email    string `json:"email" valid:"required~O campo E-mail é obrigatório"`
	Password string `json:"password" valid:"required~O campo Password é obrigatório"`
}
