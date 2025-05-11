package contextkey

// ContextKey is a custom type for context keys to avoid collisions.
type ContextKey string

const (
	// UserIDKey is the key for the user ID in the context.
	UserIDKey ContextKey = "userID"
	// HomeIDKey is the key for the home ID in the context.
	HomeIDKey ContextKey = "homeID"
	// UserRoleKey is the key for the user role in the context.
	UserRoleKey ContextKey = "userRole"
)
