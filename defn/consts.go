package defn

import "time"

type ContextKey string

const ContextUserKey ContextKey = "user"

const (
	ReadTimeout       = 15 * time.Second
	WriteTimeout      = 15 * time.Second
	ReadHeaderTimeout = 5 * time.Second
)

const (
	ContentTypeJSON       = "application/json"
	ContentTypePlainText  = "text/plain; charset=UTF-8"
	HTTPHeaderContentType = "Content-Type"
)

const (
	MaximumPasswordLength = 24

	JWTExpirationTime = 24 * time.Hour

	AuthorizationHeader = "Authorization"
	AuthorizationBearer = "Bearer "

	RoleUser = "user"
)
