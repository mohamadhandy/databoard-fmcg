package models

type RequestableAuthenticationInterface interface {
	ForAuthentication() (string, string)
}

type RequestableAuthentication struct {
	Email    string
	Password string
}

func (r *RequestableAuthentication) ForAuthentication() (string, string) {
	return r.Email, r.Password
}
