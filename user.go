package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"database/sql"
)

var err error

func createAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user_instance := &User{
		Updated_at: time.Now().Add(-24 * time.Hour),
	}
	// Decoding request
	_ = json.NewDecoder(r.Body).Decode(&user_instance)

	// converting password for database
	hashedPassword, err := hashPassword(user_instance.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user_instance.Password = hashedPassword

	//Inserting into database
	//
	var ins *sql.Stmt

	ins, err = db.Prepare("INSERT INTO `bigproject`.`project_table` (`name`,`surname`, `date_birth`,`email`,`password`,`updated_at`) VALUES(?, ?, ?, ?, ?, ?);")
	if err != nil {
		panic(err)
	}
	defer ins.Close()

	res, err := ins.Exec(&user_instance.Name, &user_instance.Surname, &user_instance.Date_birth, &user_instance.Email, &user_instance.Password, user_instance.Updated_at)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error inserting row:", err)
		return
	}

	// Returning json
	json.NewEncoder(w).Encode(user_instance)

}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var user_decode User

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user_decode); err != nil {
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	var user User
	var p string
	var i string
	var u []uint8

	row := db.QueryRow("SELECT * FROM bigproject.project_table WHERE name = ?;", user_decode.Name)
	err = row.Scan(&i, &user.Name, &user.Surname, &p, &user.Email, &user.Password, &u)
	if err != nil {
		http.Error(w, "Incorrect password or name"+err.Error(), http.StatusBadRequest)
		return

	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user_decode.Password)); err != nil {
		http.Error(w, "Incorrect password or name", http.StatusBadRequest)
		return
	}
	CreateTokenHandler_user(w, r, user_decode)
}

func getInfoUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	var p string
	var i string
	var u []uint8
	if verfied, claims := VerifyTokenHandler(w, r); verfied == true {
		row := db.QueryRow("SELECT * FROM bigproject.project_table WHERE name = ?;", claims["username"])

		if err := row.Scan(&i, &user.Name, &user.Surname, &p, &user.Email, &user.Password, &u); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//json.NewEncoder(w).Encode(user.Name)
		//json.NewEncoder(w).Encode(user.Surname)
		//json.NewEncoder(w).Encode(user.Date_birth)   NOT JSON
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

}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connection to the database

	//var user User
	var user_decodes User

	// Decode request body into the user_info struct
	err = json.NewDecoder(r.Body).Decode(&user_decodes)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//var i string
	var spc_uuu []byte
	if verfied, claims := VerifyTokenHandler(w, r); verfied == true {
		// checking last update, edit data possible once every 24 hours
		row := db.QueryRow("SELECT updated_at FROM bigproject.project_table WHERE name = ?;", claims["username"])
		row.Scan(&spc_uuu)
		// updating data
		timeString := string(spc_uuu)

		// Parse string to time.Time
		parsedTime, err := time.Parse("2006-01-02 15:04:05.9999999", timeString)
		if err != nil {
			fmt.Println(err)
		}
		timmm := time.Since(parsedTime)
		fmt.Println(parsedTime)
		fmt.Println(timmm)

		if timmm.Hours() > 24 {

			upStmt := "UPDATE `bigproject`.`project_table` SET `name` = ?, `surname` = ?, `date_birth` = ?, updated_at = ? WHERE (`name` = ?);"
			_, err := db.Exec(upStmt, user_decodes.Name, user_decodes.Surname, user_decodes.Date_birth, time.Now(), claims["username"])
			if err != nil {
				log.Println("Error executing update statement:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Too early", http.StatusBadRequest)
			return
		}

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

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
