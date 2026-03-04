package main

import (
	"net/http"
	web "filrouge/interfaces"
)

func main() {
	web.InitDB()

	http.HandleFunc("/", web.RequireAuth(web.Mainpage)) 
	http.HandleFunc("/login", web.GetLogin) 
	http.HandleFunc("/register", web.GetRegister)
	http.HandleFunc("/post-login", web.PostLogin)
	http.HandleFunc("/post-register", web.PostRegister) 
	http.HandleFunc("/logout", web.Logout)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	print("serveur lancé sur http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
