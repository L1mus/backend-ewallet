package auth

type HasherPassword interface {
	Generate(password string) string
	Compare(hash, password string) error
}
