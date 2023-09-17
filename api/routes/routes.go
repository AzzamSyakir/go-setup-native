package routes

import (
	"golang-api/api/controller"
	"golang-api/config"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Rute untuk handler CreateUser, UpdateUser, dan DeleteUser
	r.HandleFunc("/api/user/Register", controller.Register).Methods("POST")
	r.HandleFunc("/api/user/LoginUser", controller.LoginUser).Methods("POST")
	r.HandleFunc("/api/user/Update", controller.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/user/Delete", controller.DeleteUser).Methods("DELETE")

	// Rute untuk handler GetUser dengan variabel ID
	r.HandleFunc("/api/user/GetUser/{id}", controller.GetUser).Methods("GET")

	return r
}

func RunServer() {
	config.InitDB()
	router := SetupRoutes()

	// Mulai server HTTP dengan router yang telah dikonfigurasi
	http.Handle("/", router)
	http.ListenAndServe(":9000", nil)
}
