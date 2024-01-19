package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	time.Sleep(5 * time.Second)
	err := createTableIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// User functions
	r.HandleFunc("/create_account", createAccount).Methods("POST")
	r.HandleFunc("/start/login/user", loginUser).Methods("POST")
	r.HandleFunc("/login/user", getInfoUser).Methods("GET")
	r.HandleFunc("/update/user", updateUser).Methods("PUT")
	// Admin functions
	r.HandleFunc("/admin/login", login_Admin).Methods("POST")
	r.HandleFunc("/admin/get_users", getUsers_admin).Methods("GET")
	r.HandleFunc("/admin/delete/{id}", deleteUser_admin).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":7943", r))
}
