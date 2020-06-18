package api

var emptyMap = make(map[string]string)

type AuthMechanism interface {
	headers() map[string]string
	body() map[string]string
}

type UnsecureAuthentication struct {
	Email    string
	Username string
}

func (a UnsecureAuthentication) headers() map[string]string {
	return emptyMap
}

func (a UnsecureAuthentication) body() map[string]string {
	return map[string]string{
		"email":    a.Email,
		"username": a.Username,
	}
}
