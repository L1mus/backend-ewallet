package auth

type HasherPassword interface {
	Generate(password string) string
	Compare(hashPassword string) error
}
