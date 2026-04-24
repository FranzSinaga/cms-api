package shared

import "net/http"

type UserClaim struct {
	UserID string
	Email  string
	Role   string
	Name   string
}

func GetUserFromContext(r *http.Request) *UserClaim {
	user, ok := r.Context().Value(UserContextKey).(*UserClaim)
	if !ok {
		return nil
	}
	return user
}
