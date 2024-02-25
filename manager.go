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

func (m *UserManager) loginUserManager(userDecode *User) error {

	var user_instance *User
	fmt.Println("imie : ", userDecode.Name)

	user_instance, err = m.userRepository.GetUser(userDecode.Name)
	//fmt.Println("user1: ", user_instance)

	//userDecode = user_instance
	//fmt.Println("user2: ", userDecode)

	if err := bcrypt.CompareHashAndPassword([]byte(user_instance.Password), []byte(userDecode.Password)); err != nil {
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
