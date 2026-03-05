package interfaces

import (
	"database/sql"
	"filrouge/interfaces/models"
	"filrouge/interfaces/services"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	
)

var DB *sql.DB


func InitDB() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=filrouge sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")
}

func (h *PropertyHandler) Mainpage(w http.ResponseWriter, r *http.Request) {
	properties, err := h.Service.ListProperties() 
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des propriétés", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/html/index.html"))
	tmpl.Execute(w, properties) 
}

func GetLogin(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("web/html/login.html"))
	tmpl.Execute(w, nil)
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	var storedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE email=$1", email).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		return
	}
	cookie := http.Cookie{
		Name:     "session",
		Value:    email,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetRegister(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("web/html/register.html"))
	tmpl.Execute(w, nil)
}
func PostRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Champs obligatoires", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	_, err = DB.Exec("INSERT INTO users(email, password) VALUES($1, $2)", email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Email déjà utilisé", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var userRole string
		err = DB.QueryRow("SELECT role FROM users WHERE email=$1", cookie.Value).Scan(&userRole)
		if err != nil || userRole != role {
			http.Error(w, "Accès interdit", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

type PropertyHandler struct {
	Service *services.PropertyService
}

func (h *PropertyHandler) ListProperties(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filter := models.PropertyFilter{
		City: query.Get("city"),
	}

	if v := query.Get("min_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MinPrice = &f
		}
	}

	if v := query.Get("max_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MaxPrice = &f
		}
	}

	if v := query.Get("min_surface"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			filter.MinSurface = &i
		}
	}

	properties, err := h.Service.Search(filter)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/html/index.html"))
	tmpl.Execute(w, properties)
}

func (h *PropertyHandler) AddProperty(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        title := r.FormValue("title")
        description := r.FormValue("description")
        city := r.FormValue("city")
        priceStr := r.FormValue("price")       
        surfaceStr := r.FormValue("surface")   
        agencyID := r.FormValue("agency_id")  

        var agencyExists bool
        err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM agencies WHERE id = $1)", agencyID).Scan(&agencyExists)
        if err != nil || !agencyExists {
            log.Println("L'agence avec l'ID", agencyID, "n'existe pas.")
            http.Error(w, "L'agence spécifiée n'existe pas.", http.StatusBadRequest)
            return
        }

        price, err := strconv.ParseFloat(priceStr, 64)
        if err != nil {
            log.Println("Erreur conversion prix:", err)
            http.Error(w, "Prix invalide", http.StatusBadRequest)
            return
        }

        surface, err := strconv.Atoi(surfaceStr)
        if err != nil {
            log.Println("Erreur conversion surface:", err)
            http.Error(w, "Surface invalide", http.StatusBadRequest)
            return
        }
        property := models.Property{
            Title:       title,
            Description: description,
            City:        city,
            Price:       price,
            Surface:     surface,
            AgencyID:    agencyID, 
        }

        err = h.Service.AddProperty(property)
        if err != nil {
            log.Println("Erreur ajout propriété:", err)
            http.Error(w, "Erreur serveur", http.StatusInternalServerError)
            return
        }

        
        http.Redirect(w, r, "/properties", http.StatusSeeOther)
        return
    }

    rows, err := DB.Query("SELECT id, name FROM agencies")
    if err != nil {
        log.Println("Erreur lors de la récupération des agences:", err)
        http.Error(w, "Erreur serveur", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var agencies []models.Agency
    for rows.Next() {
        var agency models.Agency
        if err := rows.Scan(&agency.ID, &agency.Name); err != nil {
            log.Println("Erreur de scan des agences:", err)
            http.Error(w, "Erreur serveur", http.StatusInternalServerError)
            return
        }
        agencies = append(agencies, agency)
    }
    tmpl := template.Must(template.ParseFiles("web/html/add_property.html"))
    tmpl.Execute(w, agencies)
}

