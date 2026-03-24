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

	user, err := GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	propertyID, _ := strconv.Atoi(r.FormValue("property_id"))

	err = h.Service.AddFavorite(user.ID, propertyID)
	if err != nil {
		http.Error(w, "Erreur favoris", 500)
		return
	}

	http.Redirect(w, r, "/properties", http.StatusSeeOther)
}

func (h *FavoriteHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {

	user, err := GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	propertyID, err := strconv.Atoi(r.FormValue("property_id"))
	if err != nil {
		http.Error(w, "ID propriété invalide", http.StatusBadRequest)
		return
	}

	err = h.Service.RemoveFavorite(user.ID, propertyID)
	if err != nil {
		http.Error(w, "Erreur suppression favori", 500)
		return
	}

	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}
func (h *FavoriteHandler) FavoritesPage(w http.ResponseWriter, r *http.Request) {

	user, err := GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	properties, err := h.Service.GetFavorites(user.ID)
	if err != nil {
		http.Error(w, "Erreur favoris", 500)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/html/favorites.html"))
	err = tmpl.Execute(w, properties)
	if err != nil {
		http.Error(w, "Erreur affichage", 500)
	}
}