package main

import (
	web "filrouge/interfaces/handlers"
	"filrouge/interfaces/repositories"
	"filrouge/interfaces/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/gorilla/sessions"
    "github.com/markbates/goth/gothic"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Pas de fichier .env trouvé")
	}
	web.InitDB()
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		log.Fatal("SESSION_SECRET manquant")
	}
	gothic.Store = sessions.NewCookieStore([]byte(secret))

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			os.Getenv("APP_URL")+"/auth/google/callback",
			"email", "profile",
		),
		github.New(
			os.Getenv("GITHUB_KEY"),
			os.Getenv("GITHUB_SECRET"),
			os.Getenv("APP_URL")+"/auth/github/callback",
			"user:email",
		),
	)
	propertyRepo := &repositories.PropertyRepository{DB: web.DB}
	saleRepo := &repositories.SaleRepository{DB: web.DB}
	paymentRepo := &repositories.PaymentRepository{DB: web.DB}
	favoriteRepo := &repositories.FavoriteRepository{DB: web.DB}
	propertyService := &services.PropertyService{Repo: propertyRepo}
	saleService := &services.SaleService{
		SaleRepo:     saleRepo,
		PropertyRepo: propertyRepo,
	}
	paymentService := &services.PaymentService{Repo: paymentRepo}
	favoriteService := &services.FavoriteService{Repo: favoriteRepo}
	propertyHandler := &web.PropertyHandler{Service: propertyService}
	saleHandler := &web.SaleHandler{Service: saleService}
	paymentHandler := &web.PaymentHandler{Service: paymentService}
	favoriteHandler := &web.FavoriteHandler{Service: favoriteService}
	userHandler := &web.UserHandler{}

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.RequireAuth(propertyHandler.Mainpage))
	mux.HandleFunc("/properties", web.RequireAuth(propertyHandler.ListProperties))
	mux.HandleFunc("/property", web.RequireAuth(propertyHandler.PropertyDetail))
	mux.HandleFunc("/login", web.GetLogin)
	mux.HandleFunc("/register", web.GetRegister)
	mux.HandleFunc("/post-login", web.PostLogin)
	mux.HandleFunc("/post-register", web.PostRegister)
	mux.HandleFunc("/logout", web.Logout)
	mux.HandleFunc("/add-property",web.RequireRoles("agent", "admin")(propertyHandler.AddProperty),)
	mux.HandleFunc("/sell-property",web.RequireRoles("agent", "admin")(saleHandler.Sell),)
	mux.HandleFunc("/pay-property",web.RequireAuth(paymentHandler.CreateCheckoutSession),)
	mux.HandleFunc("/sold", web.RequireAuth(web.SoldPage))
	mux.HandleFunc("/payment-success", web.RequireAuth(web.PaymentSuccess))
	mux.HandleFunc("/favorite/add", web.RequireAuth(favoriteHandler.AddFavorite))
	mux.HandleFunc("/favorite/remove", web.RequireAuth(favoriteHandler.RemoveFavorite))
	mux.HandleFunc("/favorites", web.RequireAuth(favoriteHandler.FavoritesPage))
	mux.HandleFunc("/profile", web.RequireAuth(userHandler.ProfilePage))
	mux.HandleFunc("/update-account", web.RequireAuth(userHandler.UpdateAccount))
	mux.HandleFunc("/update", web.RequireAuth(userHandler.UpdatePage))
	mux.HandleFunc("/edit-property",web.RequireRoles("agent", "admin")(propertyHandler.EditProperty),)
	mux.HandleFunc("/update-property",web.RequireRoles("agent", "admin")(propertyHandler.UpdateProperty),)
	mux.HandleFunc("/delete-property",web.RequireRoles("agent", "admin")(propertyHandler.DeleteProperty),)
	mux.HandleFunc("/dashboard",web.RequireRoles("admin")(propertyHandler.Dashboard),)

	mux.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {

	
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		http.Error(w, "Provider manquant", http.StatusBadRequest)
		return
	}

	provider := parts[2]

	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	if strings.Contains(r.URL.Path, "/callback") {
		web.CallbackAuth(w, r)
		return
	}

	web.BeginAuth(w, r)
})
	mux.HandleFunc("/auth/callback", web.CallbackAuth)
	mux.Handle("/web/",http.StripPrefix("/web/", http.FileServer(http.Dir("web"))),)
	mux.Handle("/uploads/",http.StripPrefix("/uploads/", http.FileServer(http.Dir("web/uploads"))),)
	log.Println("Serveur lancé sur http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
