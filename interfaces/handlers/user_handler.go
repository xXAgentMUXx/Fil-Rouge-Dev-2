package interfaces

import (
	"html/template"
	"net/http"
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

	// -------- PROPRIETES ACHETEES --------

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
			err := rowsBuy.Scan(&title)
			if err == nil {
				purchases = append(purchases, title)
			}
		}
	}

	// -------- PROPRIETES VENDUES --------

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
			err := rowsSell.Scan(&title)
			if err == nil {
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