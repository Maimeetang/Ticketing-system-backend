package port

type AuthService interface {
	Login(username, password string) (string, error)
}