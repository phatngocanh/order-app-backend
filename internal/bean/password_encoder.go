package bean

type PasswordEncoder interface {
	Encrypt(password string) (string, error)
	Compare(password, hash string) bool
}
