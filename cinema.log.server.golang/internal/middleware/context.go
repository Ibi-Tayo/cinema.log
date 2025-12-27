package middleware

type ContextKey int

const (
	KeyPrincipalID ContextKey = iota
	KeyUser
)
