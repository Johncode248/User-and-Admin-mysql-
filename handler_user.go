package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var err error

type UserHandler struct {
	userManager *UserManager
}

func (h *UserHandler) createAccountHandler(w http.ResponseWriter, r *http.Request) {

	var user_instance User

	_ = json.NewDecoder(r.Body).Decode(&user_instance)

	user, err := h.userManager.CreateUser(user_instance.Name, user_instance.Surname, user_instance.Email, user_instance.Password, user_instance.Date_birth)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)

}

func (h *UserHandler) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user_decode User

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user_decode); err != nil {
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	err = h.userManager.loginUserManager(&user_decode)
	if err != nil {
		http.Error(w, "Incorrect password or name"+err.Error(), http.StatusBadRequest)
		return
	}

	CreateTokenHandler(w, r, &user_decode)
}

func (h *UserHandler) getInfoUserHandler(w http.ResponseWriter, r *http.Request) {

	var user User

	if verfied, claims := VerifyTokenHandler(w, r); verfied == true {
		user, err = h.userManager.getInfoUser(claims)
		if err != nil {
			log.Fatal(err)
		}

		userInfo := map[string]interface{}{
			"name":    user.Name,
			"surname": user.Surname,
			"email":   user.Email,
		}
		// Returning JSON
		if err := json.NewEncoder(w).Encode(userInfo); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "invalid token", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (h *UserHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user_decodes User

	// Decode request body into the user_info struct
	err = json.NewDecoder(r.Body).Decode(&user_decodes)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if verfied, claims := VerifyTokenHandler(w, r); verfied == true {
		err = h.userManager.updateUser(user_decodes, claims)

		userInfo := map[string]interface{}{
			"name":       user_decodes.Name,
			"surname":    user_decodes.Surname,
			"date_birth": user_decodes.Date_birth,
			"email":      user_decodes.Email,
			"password":   user_decodes.Password,
		}
		// Returning JSON
		if err := json.NewEncoder(w).Encode(userInfo); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
