package port

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) bool
}
