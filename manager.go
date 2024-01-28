package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserManager struct {
	userRepository UserRepository
}

func (m *UserManager) CreateUser(name, surname, email, password string, date_birth time.Time) (*User, error) {

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	// Wywołanie repozytorium w celu utworzenia użytkownika
	user, err := m.userRepository.Create(name, surname, email, hashedPassword, date_birth)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *UserManager) loginUserManager(*User) error {

	var user_instance *User

	user_instance, err = m.userRepository.GetUser(user_instance.Name)

	if err := bcrypt.CompareHashAndPassword([]byte(user_instance.Password), []byte(user_instance.Password)); err != nil {
		return err
	}

	return nil
}

func (m *UserManager) getInfoUser(claims jwt.MapClaims) (User, error) {

	user, err := m.userRepository.GetRow(claims)
	if err != nil {
		log.Println(err)
	}

	return user, err

}

func (m *UserManager) updateUser(user_decodes User, claims jwt.MapClaims) error {
	// checking last update, edit data possible once every 24 hours
	//var user_decodes User

	user, err := m.userRepository.GetRow(claims)
	if err != nil {
		return err
	}
	timeString := user.Updated_at

	timmm := time.Since(timeString)
	fmt.Println(timmm)

	if timmm.Hours() > 24 {
		err = m.userRepository.Update(user_decodes, claims)

		if err != nil {
			log.Println("Error executing update statement:", err)

			return err
		}
	} else {
		fmt.Println("Too early")
	}

	return nil
}

// Admin  functions

type Page struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Users    []User `json:"users"`
}

/*
	func (m *UserManager) getUsersAdmin(page int) (int, []User) {
		// Construct the SQL query with LIMIT and OFFSET
		var users []User
		pageSize := 10
		query := fmt.Sprintf("SELECT name, surname, date_birth, email, password, updated_at FROM bigproject.project_table LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)

		// Execute the query
		rows, err := db.Query(query)
		if err != nil {
			log.Println(err)

			return 0, nil
		}
		defer rows.Close()

		var u string
		var u2 string
		var user User
		for rows.Next() {

			if err := rows.Scan(&user.Name, &user.Surname, &u2, &user.Email, &user.Password, &u); err == nil {

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

				return 0, nil
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {

			log.Println(err)

			return 0, nil
		}
		return page, users
	}
*/
func (m *UserManager) getUsersAdmin(page int) (int, []User) {

	users, err := m.userRepository.List(page)
	if err != nil {
		log.Println(err)
		return 0, nil
	}

	return page, users
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
