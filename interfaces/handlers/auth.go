package interfaces

import (
	"net/http"

	"github.com/markbates/goth/gothic"
)

func BeginAuth(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func CallbackAuth(w http.ResponseWriter, r *http.Request) {

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := user.Email
	provider := user.Provider

	var exists bool

	err = DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)",
		email,
	).Scan(&exists)

	if err != nil {
		http.Error(w, "Erreur DB", 500)
		return
	}

	if exists {

		_, err = DB.Exec(
			"UPDATE users SET provider=$1 WHERE email=$2",
			provider,
			email,
		)

	} else {

		_, err = DB.Exec(
			"INSERT INTO users(email,password,role,provider) VALUES($1,'oauth','client',$2)",
			email,
			provider,
		)

	}

	if err != nil {
		http.Error(w, "Erreur utilisateur", 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    email,
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}