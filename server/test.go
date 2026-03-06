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
	saleHandler := &web.SaleHandler{Service: saleService}
	paymentRepo := &repositories.PaymentRepository{DB: web.DB}
	paymentService := &services.PaymentService{Repo: paymentRepo}
	paymentHandler := &web.PaymentHandler{Service: paymentService}

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
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	println("Serveur lancé sur http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}