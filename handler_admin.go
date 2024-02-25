package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

type AdminHandler struct {
	userManager *UserManager
}

var (
	admin = Admin{Name: "Admin", Password: "12345"}
)

func (a *AdminHandler) loginAdminHandler(w http.ResponseWriter, r *http.Request) {
	var admin_decode Admin

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&admin_decode); err != nil {
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	if admin_decode.Name == admin.Name && admin_decode.Password == admin.Password {
		//json.NewEncoder(w).Encode(r.Body)  Empty JSON
		CreateTokenHandler(w, r, &admin_decode)
	} else {
		http.Error(w, "wrong name or password", http.StatusBadRequest)
		return
	}

}

func (a *AdminHandler) deleteUserAdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	productId := params["id"]
	//email := params["email"] // id = user's mail

	if verfied, claims := VerifyTokenHandler(w, r); verfied == true && claims["username"].(string) == "Admin" {
		err = a.userManager.userRepository.Delate(productId)
		//
		if err != nil {
			panic(err)
		}

		//rowsAffected := del.RowsAffected
		//fmt.Println("Rows Affected:", rowsAffected)

	} else {
		http.Error(w, "Invalid token", http.StatusBadRequest)
	}

}

func (a *AdminHandler) getUsersAdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if verfied, claims := VerifyTokenHandler(w, r); verfied == true && claims["username"].(string) == "Admin" {
		fmt.Println("Token accepted")
	} else {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	page := 1
	pageSize := 10

	var users []User
	if r.URL.Query().Get("page") != "" {
		if val, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
			page = val
		}
	}

	page, users = a.userManager.getUsersAdmin(page)

	json.NewEncoder(w).Encode(Page{Page: page, PageSize: pageSize, Users: users})
}
