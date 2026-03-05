package main

import (
	web "filrouge/interfaces/handlers"
	"filrouge/interfaces/repositories"
	"filrouge/interfaces/services"
	"net/http"
)

func main() {
	web.InitDB()
	propertyHandler := &web.PropertyHandler{
		Service: &services.PropertyService{
			Repo: &repositories.PropertyRepository{DB: web.DB},
		},
	}

	http.HandleFunc("/", web.RequireAuth(propertyHandler.Mainpage))
	http.HandleFunc("/login", web.GetLogin) 
	http.HandleFunc("/register", web.GetRegister)
	http.HandleFunc("/post-login", web.PostLogin)
	http.HandleFunc("/post-register", web.PostRegister) 
	http.HandleFunc("/logout", web.Logout)
	http.HandleFunc("/properties", propertyHandler.ListProperties)      
	http.HandleFunc("/add-property", propertyHandler.AddProperty)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	print("serveur lancé sur http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
