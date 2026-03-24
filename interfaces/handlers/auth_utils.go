package interfaces

import (
	"net/http"
)

type User struct {
	ID       int
	Email    string
	Role     string
	Provider string
}

func GetCurrentUser(r *http.Request) (*User, error) {

	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		return nil, err
	}

	var user User

	err = DB.QueryRow(`
		SELECT id, email, role, provider
		FROM users
		WHERE email = $1
	`, cookie.Value).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.Provider,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}