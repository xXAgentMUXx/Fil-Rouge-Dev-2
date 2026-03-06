package interfaces

import (
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct{}

type ProfileData struct {
	Email     string
	Role      string
	Purchases []string
	Sales     []string
}

func (h *UserHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var email string
	var role string

	err = DB.QueryRow(
		"SELECT email, role FROM users WHERE email=$1",
		cookie.Value,
	).Scan(&email, &role)

	if err != nil {
		http.Error(w, "Utilisateur introuvable", 500)
		return
	}

	rowsBuy, err := DB.Query(`
		SELECT p.title
		FROM payments pay
		JOIN properties p ON p.id = pay.property_id
		WHERE pay.buyer_id = (
			SELECT id FROM users WHERE email=$1
		) AND pay.status='paid'
	`, email)

	var purchases []string

	if err == nil {
		defer rowsBuy.Close()
		for rowsBuy.Next() {
			var title string
			if err := rowsBuy.Scan(&title); err == nil {
				purchases = append(purchases, title)
			}
		}
	}

	rowsSell, err := DB.Query(`
		SELECT p.title
		FROM sales s
		JOIN properties p ON p.id = s.property_id
		WHERE s.agent_id = (
			SELECT id FROM users WHERE email=$1
		)
	`, email)

	var sales []string

	if err == nil {
		defer rowsSell.Close()
		for rowsSell.Next() {
			var title string
			if err := rowsSell.Scan(&title); err == nil {
				sales = append(sales, title)
			}
		}
	}

	data := ProfileData{
		Email:     email,
		Role:      role,
		Purchases: purchases,
		Sales:     sales,
	}

	tmpl := template.Must(template.ParseFiles("web/html/profile.html"))
	tmpl.Execute(w, data)
}

func (h *UserHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	newEmail := r.FormValue("email")
	newPassword := r.FormValue("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur hash password", 500)
		return
	}

	_, err = DB.Exec(
		"UPDATE users SET email=$1, password=$2 WHERE email=$3",
		newEmail,
		string(hash),
		cookie.Value,
	)

	if err != nil {
		http.Error(w, "Erreur mise à jour compte", 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: newEmail,
		Path:  "/",
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *UserHandler) UpdatePage(w http.ResponseWriter, r *http.Request) {

	_, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/html/update.html"))
	tmpl.Execute(w, nil)
}