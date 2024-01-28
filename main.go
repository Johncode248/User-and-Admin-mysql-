package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	userManager := &UserManager{
		userRepository: &DatabaseUserRepository{},
	}

	userHandler := &UserHandler{
		userManager: userManager,
	}

	adminHandler := &AdminHandler{
		userManager: userManager,
	}
	fmt.Print(adminHandler)

	DatabaseConnect()

	r := mux.NewRouter()

	// User functions
	r.HandleFunc("/create_account", userHandler.createAccountHandler).Methods("POST")
	r.HandleFunc("/start/login/user", userHandler.loginUserHandler).Methods("POST")
	r.HandleFunc("/login/user", userHandler.getInfoUserHandler).Methods("GET")
	r.HandleFunc("/update/user", userHandler.updateUserHandler).Methods("PUT")
	// Admin functions
	r.HandleFunc("/admin/login", adminHandler.loginAdminHandler).Methods("POST")
	r.HandleFunc("/admin/get_users", adminHandler.getUsersAdminHandler).Methods("GET")
	r.HandleFunc("/admin/delete/{id}", adminHandler.deleteUserAdminHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":7943", r))
}
