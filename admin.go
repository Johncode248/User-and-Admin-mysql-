package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func login_Admin(w http.ResponseWriter, r *http.Request) {
	var admin_decode Admin
	admin := Admin{Name: "Admin", Password: "12345"}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&admin_decode); err != nil {
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	if admin_decode.Name == admin.Name && admin_decode.Password == admin.Password {
		//json.NewEncoder(w).Encode(r.Body)  Empty JSON
		CreateTokenHandler(w, r, admin_decode)
	} else {
		http.Error(w, "wrong name or password", http.StatusBadRequest)
		return
	}

}

type Page struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Users    []User `json:"users"`
}

func deleteUser_admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get connect
	db, err := OpenDatabaseConnection()
	if err != nil {
		fmt.Println("Błąd podczas otwierania połączenia do bazy danych:", err)
		return
	}
	defer db.Close()

	//var user User
	fmt.Println(w, "sdad")

	params := mux.Vars(r)
	productId := params["id"]
	//email := params["email"] // id = user's mail

	if verfied, claims := VerifyTokenHandler(w, r); verfied == true && claims["username"].(string) == "Admin" {
		del, err := db.Exec("DELETE FROM `bigproject`.`project_table` WHERE (`email` = ?)", productId)
		if err != nil {
			panic(err)
		}

		rowsAffected := del.RowsAffected
		fmt.Println("Rows Affected:", rowsAffected)

	} else {
		http.Error(w, "Invalid token", http.StatusBadRequest)
	}

}

func getUsers_admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if verfied, claims := VerifyTokenHandler(w, r); verfied == true && claims["username"].(string) == "Admin" {
		fmt.Println("Token accepted")
	} else {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	var users []User
	db, err := OpenDatabaseConnection()
	if err != nil {
		fmt.Println("Błąd podczas otwierania połączenia do bazy danych:", err)
		return
	}
	defer db.Close()

	page := 1
	pageSize := 10

	if r.URL.Query().Get("page") != "" {
		if val, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
			page = val
		}
	}
	// Construct the SQL query with LIMIT and OFFSET
	query := fmt.Sprintf("SELECT name, surname, date_birth, email, password, updated_at FROM bigproject.project_table LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var u string
	var u2 string // to trzeba sprawdic !!!
	var user User
	for rows.Next() {

		//Below is a bug       CHECK THESE &u !!!!!!!!!!!!!
		if err := rows.Scan(&user.Name, &user.Surname, &u2, &user.Email, &user.Password, &u); err == nil {

			//timeString := string(u)
			//timeString2 := string(u2)

			// Parse string to time.Time

			user.Date_birth, err = time.Parse("2006-01-02", u2)
			if err != nil {
				fmt.Println("Error parsing Date_birth:", err)
			} else {
				fmt.Println("good")
			}

			user.Updated_at, err = time.Parse("2006-01-02 15:04:05.9999999", u)
			if err != nil {
				fmt.Println("Error parsing Updated_at:", err)
			} else {
				fmt.Println("good")
			}

		} else {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError) // here shows a  bug !!!   if you delete this line it will return date_birth and updated_at as 0001-00-00 00:00:00
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Page{Page: page, PageSize: pageSize, Users: users})
}
