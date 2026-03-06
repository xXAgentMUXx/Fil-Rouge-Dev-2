package interfaces

import (
	"filrouge/interfaces/services"
	"html/template"
	"net/http"
	"strconv"
	_ "github.com/lib/pq"
)

type FavoriteHandler struct {
	Service *services.FavoriteService
}

func (h *FavoriteHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {

	propertyID, _ := strconv.Atoi(r.FormValue("property_id"))

	cookie, _ := r.Cookie("session")

	var userID int

	DB.QueryRow(
		"SELECT id FROM users WHERE email=$1",
		cookie.Value,
	).Scan(&userID)

	err := h.Service.AddFavorite(userID, propertyID)

	if err != nil {
		http.Error(w, "Erreur favoris", 500)
		return
	}

	http.Redirect(w, r, "/properties", http.StatusSeeOther)
}

func (h *FavoriteHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {

	propertyID, _ := strconv.Atoi(r.FormValue("property_id"))

	cookie, _ := r.Cookie("session")

	var userID int

	DB.QueryRow(
		"SELECT id FROM users WHERE email=$1",
		cookie.Value,
	).Scan(&userID)

	h.Service.RemoveFavorite(userID, propertyID)

	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}

func (h *FavoriteHandler) FavoritesPage(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")

	var userID int

	DB.QueryRow(
		"SELECT id FROM users WHERE email=$1",
		cookie.Value,
	).Scan(&userID)

	properties, err := h.Service.GetFavorites(userID)

	if err != nil {
		http.Error(w, "Erreur favoris", 500)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/html/favorites.html"))

	tmpl.Execute(w, properties)
}