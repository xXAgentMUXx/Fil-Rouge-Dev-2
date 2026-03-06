package main

import (
	web "filrouge/interfaces/handlers"
	"filrouge/interfaces/repositories"
	"filrouge/interfaces/services"
	"net/http"
)

func main() {
	web.InitDB()
	propertyRepo := &repositories.PropertyRepository{DB: web.DB}
	propertyService := &services.PropertyService{Repo: propertyRepo}
	propertyHandler := &web.PropertyHandler{Service: propertyService}
	saleRepo := &repositories.SaleRepository{DB: web.DB}
	saleService := &services.SaleService{
		SaleRepo:     saleRepo,
		PropertyRepo: propertyRepo,
	}
	userHandler := &web.UserHandler{}
	saleHandler := &web.SaleHandler{Service: saleService}
	paymentRepo := &repositories.PaymentRepository{DB: web.DB}
	paymentService := &services.PaymentService{Repo: paymentRepo}
	paymentHandler := &web.PaymentHandler{Service: paymentService}

	favoriteRepo := &repositories.FavoriteRepository{DB: web.DB}
	favoriteService := &services.FavoriteService{Repo: favoriteRepo}
	favoriteHandler := &web.FavoriteHandler{Service: favoriteService}

	http.HandleFunc("/", web.RequireAuth(propertyHandler.Mainpage))
	http.HandleFunc("/properties", web.RequireAuth(propertyHandler.ListProperties))
	http.HandleFunc("/login", web.GetLogin)
	http.HandleFunc("/register", web.GetRegister)
	http.HandleFunc("/post-login", web.PostLogin)
	http.HandleFunc("/post-register", web.PostRegister)
	http.HandleFunc("/logout", web.Logout)
	http.HandleFunc("/add-property",web.RequireRoles("agent", "admin")(propertyHandler.AddProperty),)
	http.HandleFunc("/sell-property",web.RequireRoles("agent", "admin")(saleHandler.Sell),)
	http.HandleFunc("/pay-property", web.RequireAuth(paymentHandler.CreateCheckoutSession))
	http.HandleFunc("/sold", web.RequireAuth(web.SoldPage))
	http.HandleFunc("/payment-success", web.RequireAuth(web.PaymentSuccess))
	http.HandleFunc("/favorite/add", web.RequireAuth(favoriteHandler.AddFavorite))
	http.HandleFunc("/favorite/remove", web.RequireAuth(favoriteHandler.RemoveFavorite))
	http.HandleFunc("/favorites", web.RequireAuth(favoriteHandler.FavoritesPage))
	http.HandleFunc("/profile", web.RequireAuth(userHandler.ProfilePage))
	http.HandleFunc("/update-account", userHandler.UpdateAccount)
	http.HandleFunc("/update", userHandler.UpdatePage)
	http.HandleFunc("/edit-property",web.RequireRoles("agent","admin")(propertyHandler.EditProperty))
	http.HandleFunc("/update-property",web.RequireRoles("agent","admin")(propertyHandler.UpdateProperty))
	http.HandleFunc("/delete-property",web.RequireRoles("agent","admin")(propertyHandler.DeleteProperty))
	http.HandleFunc("/dashboard", web.RequireRoles("admin",)(propertyHandler.Dashboard))
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	println("Serveur lancé sur http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}