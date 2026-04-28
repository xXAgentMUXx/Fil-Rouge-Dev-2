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
	Provider  string
	Purchases []string
	Sales     []string
}

func (h *UserHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {

	user, err := GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	rowsBuy, _ := DB.Query(`
		SELECT p.title
		FROM payments pay
		JOIN properties p ON p.id = pay.property_id
		WHERE pay.buyer_id = $1 AND pay.status='paid'
	`, user.ID)

	var purchases []string
	if rowsBuy != nil {
		defer rowsBuy.Close()
		for rowsBuy.Next() {
			var title string
			rowsBuy.Scan(&title)
			purchases = append(purchases, title)
		}
	}
	rowsSell, _ := DB.Query(`
		SELECT p.title
		FROM sales s
		JOIN properties p ON p.id = s.property_id
		WHERE p.agent_id = $1
	`, user.ID)

	var sales []string
	if rowsSell != nil {
		defer rowsSell.Close()
		for rowsSell.Next() {
			var title string
			rowsSell.Scan(&title)
			sales = append(sales, title)
		}
	}
	data := ProfileData{
		Email:     user.Email,
		Role:      user.Role,
		Provider:  user.Provider,
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
	user, err := GetCurrentUser(r)
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
	if user.Provider != "local" {
		_, err = DB.Exec(
			"UPDATE users SET email=$1, password=$2, provider='local' WHERE id=$3",
			newEmail,
			string(hash),
			user.ID,
		)
	} else {
		_, err = DB.Exec(
			"UPDATE users SET email=$1, password=$2 WHERE id=$3",
			newEmail,
			string(hash),
			user.ID,
		)
	}
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
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var provider string
	err = DB.QueryRow(
		"SELECT provider FROM users WHERE email=$1",
		cookie.Value,
	).Scan(&provider)

	if err != nil {
		http.Error(w, "Utilisateur introuvable", 500)
		return
	}
	data := struct {
		Provider string
	}{
		Provider: provider,
	}
	tmpl := template.Must(template.ParseFiles("web/html/update.html"))
	tmpl.Execute(w, data)
}