package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driverName1     = "mysql"
	dataSourceName1 = "root:password@tcp(localhost:3306)/database_bigproject"
)

func main() {

	err := createTableIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	//db, err := sql.Open(driverName, dataSourceName)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()

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
