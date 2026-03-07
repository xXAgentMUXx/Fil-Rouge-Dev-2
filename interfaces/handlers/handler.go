package interfaces

import (
	"database/sql"
	"filrouge/interfaces/models"
	"filrouge/interfaces/services"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"io"
	"os"

	_ "github.com/lib/pq"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

type PropertyPageData struct {
	Properties []models.Property
	UserRole   string
}

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
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var role string
	err = DB.QueryRow(
		"SELECT role FROM users WHERE email = $1",
		cookie.Value,
	).Scan(&role)
	if err == sql.ErrNoRows {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
	}
	if err != nil {
		log.Println("Erreur récupération rôle:", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	data := PropertyPageData{
		Properties: properties,
		UserRole:   role,
	}
	tmpl := template.Must(template.ParseFiles("web/html/index.html"))
	tmpl.Execute(w, data)
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

func RequireRoles(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			var role string
			err = DB.QueryRow(
				"SELECT role FROM users WHERE email = $1",
				cookie.Value,
			).Scan(&role)
			if err != nil {
				http.Error(w, "Accès interdit", http.StatusForbidden)
				return
			}

			for _, allowed := range roles {
				if role == allowed {
					next(w, r)
					return
				}
			}

			http.Error(w, "Accès interdit", http.StatusForbidden)
		}
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
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var role string
	err = DB.QueryRow(
		"SELECT role FROM users WHERE email = $1",
		cookie.Value,
	).Scan(&role)
	if err == sql.ErrNoRows {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
	}
	if err != nil {
		log.Println("Erreur récupération rôle:", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	data := PropertyPageData{
		Properties: properties,
		UserRole:   role,
	}
	tmpl := template.Must(template.ParseFiles("web/html/index.html"))
	tmpl.Execute(w, data)
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
		err := DB.QueryRow(
			"SELECT EXISTS (SELECT 1 FROM agencies WHERE id = $1)",
			agencyID,
		).Scan(&agencyExists)

		if err != nil || !agencyExists {
			log.Println("L'agence avec l'ID", agencyID, "n'existe pas.")
			http.Error(w, "L'agence spécifiée n'existe pas.", http.StatusBadRequest)
			return
		}

		cookie, _ := r.Cookie("session")

		var agentID int
		DB.QueryRow(
			"SELECT id FROM users WHERE email=$1",
			cookie.Value,
		).Scan(&agentID)

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

		// -------------------------------
		// Upload image (optionnel)
		// -------------------------------

		var imagePath string

		file, handler, err := r.FormFile("image")

		if err == nil {
			defer file.Close()

			os.MkdirAll("web/uploads", os.ModePerm)

			dst, err := os.Create("web/uploads/" + handler.Filename)
			if err == nil {
				defer dst.Close()

				_, err = io.Copy(dst, file)
				if err == nil {
					imagePath = "/uploads/" + handler.Filename
				}
			}
		}

		property := models.Property{
			Title:       title,
			Description: description,
			City:        city,
			Price:       price,
			Surface:     surface,
			AgencyID:    agencyID,
			AgentID:     agentID,
			Image:       imagePath,
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

type SaleHandler struct {
	Service *services.SaleService
}

func (h *SaleHandler) Sell(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// --- Conversion des IDs ---
	propertyIDStr := r.FormValue("property_id")
	buyerIDStr := r.FormValue("buyer_id")
	priceStr := r.FormValue("sale_price")

	propertyID, err := strconv.Atoi(propertyIDStr)
	if err != nil {
		http.Error(w, "ID de propriété invalide", http.StatusBadRequest)
		return
	}

	buyerID, err := strconv.Atoi(buyerIDStr)
	if err != nil {
		http.Error(w, "ID acheteur invalide", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		http.Error(w, "Prix invalide", http.StatusBadRequest)
		return
	}
	err = h.Service.Sell(propertyID, buyerID, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/properties", http.StatusSeeOther)
}

type PaymentHandler struct {
	Service *services.PaymentService
}

func (h *PaymentHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.Atoi(r.FormValue("property_id"))
	if err != nil {
		http.Error(w, "ID propriété invalide", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
		return
	}
	var buyerID int
	err = DB.QueryRow(
		"SELECT id FROM users WHERE email=$1",
		cookie.Value,
	).Scan(&buyerID)

	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusInternalServerError)
		return
	}
	var price float64
	var isSold bool

	err = DB.QueryRow(
		"SELECT price, is_sold FROM properties WHERE id=$1",
		propertyID,
	).Scan(&price, &isSold)

	if err != nil {
		http.Error(w, "Propriété introuvable", http.StatusInternalServerError)
		return
	}
	if isSold {
		http.Error(w, "Propriété déjà vendue", http.StatusBadRequest)
		return
	}
	stripe.Key = "sk_test_51T7v2sPVJxAvAHPXWguzBmPqfEqWzQfMmWF3s6s9Bpm2StsMRvD6Xf3SKnCbJEB2LhH3rLtYhPvUny0dS9aw4vUS00TWfSIIRJ"

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(
			"http://localhost:8080/payment-success?session_id={CHECKOUT_SESSION_ID}",
		),
		CancelURL: stripe.String("http://localhost:8080/payment-cancel"),

		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{

					Currency: stripe.String("eur"),

					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(
							fmt.Sprintf("Achat propriété #%d", propertyID),
						),
					},
					UnitAmount: stripe.Int64(int64(price * 100)),
				},
			},
		},
	}
	s, err := session.New(params)
	if err != nil {
		log.Println("Erreur création session Stripe:", err)
		http.Error(w, "Erreur Stripe", http.StatusInternalServerError)
		return
	}
	payment := models.Payment{
		PropertyID:      propertyID,
		BuyerID:         buyerID,
		Amount:          price,
		StripeSessionID: s.ID,
		Status:          "pending",
	}
	err = h.Service.CreatePayment(payment)
	if err != nil {
		log.Println("Erreur insertion paiement :", err)
		http.Error(w, "Erreur enregistrement paiement", 500)
		return
	}
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
}

type SoldPageData struct {
	PropertyID string
	Price      string
}

func SoldPage(w http.ResponseWriter, r *http.Request) {

	propertyID := r.URL.Query().Get("property_id")
	price := r.URL.Query().Get("price")

	data := SoldPageData{
		PropertyID: propertyID,
		Price:      price,
	}
	tmpl := template.Must(template.ParseFiles("web/html/sold.html"))
	tmpl.Execute(w, data)
}
func PaymentSuccess(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	if sessionID == "" {
		http.Error(w, "Session Stripe manquante", http.StatusBadRequest)
		return
	}
	var propertyID int
	err := DB.QueryRow(
		"SELECT property_id FROM payments WHERE stripe_session_id=$1",
		sessionID,
	).Scan(&propertyID)

	if err != nil {
		log.Println("Paiement introuvable:", err)
		http.Error(w, "Paiement introuvable", 500)
		return
	}
	_, err = DB.Exec(
		"UPDATE payments SET status='paid' WHERE stripe_session_id=$1",
		sessionID,
	)

	if err != nil {
		log.Println("Erreur update paiement:", err)
		http.Error(w, "Erreur paiement", 500)
		return
	}
	_, err = DB.Exec(
		"UPDATE properties SET is_sold=true WHERE id=$1",
		propertyID,
	)

	if err != nil {
		log.Println("Erreur update propriété:", err)
		http.Error(w, "Erreur propriété", 500)
		return
	}
	http.Redirect(w, r, "/properties", http.StatusSeeOther)
}

func (h *PropertyHandler) EditProperty(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	property, err := h.Service.GetProperty(id)

	if err != nil {
		http.Error(w, "Propriété introuvable", 404)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/html/edit_property.html"))
	tmpl.Execute(w, property)
}

func (h *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Erreur formulaire", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	title := r.FormValue("title")
	description := r.FormValue("description")
	city := r.FormValue("city")

	file, handler, err := r.FormFile("image")

	var filename string

	if err == nil {
		defer file.Close()

		filename = handler.Filename
		path := "web/uploads/" + filename

		dst, err := os.Create(path)
		if err != nil {
			http.Error(w, "Erreur upload", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		io.Copy(dst, file)
	}

	property := models.Property{
	ID:          id,
	Title:       title,
	Description: description,
	City:        city,
	Image:       filename,
	}

	err = h.Service.UpdateProperty(property)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/properties", http.StatusSeeOther)
}

func (h *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.FormValue("id"))

	property, err := h.Service.GetProperty(id)

	if err != nil {
		http.Error(w, "Propriété introuvable", 404)
		return
	}

	cookie, _ := r.Cookie("session")

	var userID int
	DB.QueryRow(
		"SELECT id FROM users WHERE email=$1",
		cookie.Value,
	).Scan(&userID)

	if property.AgentID != userID {
		http.Error(w, "Non autorisé", 403)
		return
	}

	err = h.Service.DeleteProperty(id)

	if err != nil {
		http.Error(w, "Erreur suppression", 500)
		return
	}

	http.Redirect(w, r, "/properties", http.StatusSeeOther)
}

func (h *PropertyHandler) Dashboard(w http.ResponseWriter, r *http.Request) {

	totalSales, _ := h.Service.GetTotalSales()
	totalSold, _ := h.Service.GetTotalSoldProperties()
	topCities, _ := h.Service.GetTopCities()
	expensiveProperties, _ := h.Service.GetMostExpensive()

	data := map[string]interface{}{
		"TotalSales": totalSales,
		"TotalSold": totalSold,
		"TopCities": topCities,
		"ExpensiveProperties": expensiveProperties,
	}

	tmpl := template.Must(template.ParseFiles("web/html/dashboard.html"))
	tmpl.Execute(w, data)
}

func (h *PropertyHandler) PropertyDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	property, err := h.Service.GetProperty(id)
	if err != nil {
		http.Error(w, "Propriété introuvable", http.StatusNotFound)
		return
	}
	tmpl := template.Must(template.ParseFiles("web/html/property_detail.html"))
	tmpl.Execute(w, property)
}