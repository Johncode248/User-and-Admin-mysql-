package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserRepository interface {
	//DatabaseConnect()
	Create(name, surname, email, password string, date_birth time.Time) (*User, error)

	GetUser(name string) (*User, error)
	GetRow(claims jwt.MapClaims) (User, error)

	List(page int) ([]User, error)
	Update(user_decodes User, claims jwt.MapClaims) error
	Delate(email string) error
}

type DatabaseUserRepository struct {
	// Implementacja połączenia z bazą danych
}

func (r *DatabaseUserRepository) Delate(email string) error {
	_, err := db.Exec("DELETE FROM `bigproject`.`project_table` WHERE (`email` = ?)", email)
	return err
}

func (r *DatabaseUserRepository) Update(user_decodes User, claims jwt.MapClaims) error {
	upStmt := "UPDATE `bigproject`.`project_table` SET `name` = ?, `surname` = ?, `date_birth` = ?, updated_at = ? WHERE (`name` = ?);"
	_, err := db.Exec(upStmt, user_decodes.Name, user_decodes.Surname, user_decodes.Date_birth, time.Now(), claims["username"])
	return err
}

func (r *DatabaseUserRepository) GetRow(claims jwt.MapClaims) (User, error) {
	var user User
	var p string
	var i string
	var u []uint8
	row := db.QueryRow("SELECT * FROM bigproject.project_table WHERE name = ?;", claims["username"])

	err := row.Scan(&i, &user.Name, &user.Surname, &p, &user.Email, &user.Password, &u)

	return user, err
}
func DatabaseConnect() {
	err := createTableIfNotExists()
	if err != nil {

		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 2)
			err = createTableIfNotExists()
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Fatal("Failed to connect to database", err)
		}
	}
}

func (r *DatabaseUserRepository) Create(name, surname, email, password string, date_birth time.Time) (*User, error) {
	// Implementacja operacji tworzenia użytkownika w bazie danych
	updated_at := time.Now().Add(-24 * time.Hour) //updated_at := time.Now().Add(-24 * time.Hour)
	user := &User{
		Name:       name,
		Surname:    surname,
		Email:      email,
		Password:   password,
		Date_birth: date_birth,
		Updated_at: updated_at,
	}

	var ins *sql.Stmt

	ins, err = db.Prepare("INSERT INTO `bigproject`.`project_table` (`name`,`surname`, `date_birth`,`email`,`password`,`updated_at`) VALUES(?, ?, ?, ?, ?, ?);")
	if err != nil {
		panic(err)
	}
	defer ins.Close()

	res, err := ins.Exec(name, surname, date_birth, email, password, updated_at)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		log.Fatal(err)
		return nil, err
	}

	return user, nil
}

func (r *DatabaseUserRepository) GetUser(name string) (*User, error) {
	user_instance := &User{}
	var p string
	var i string
	var u []uint8
	row := db.QueryRow("SELECT * FROM bigproject.project_table WHERE name = ?;", name)
	err = row.Scan(&i, &user_instance.Name, &user_instance.Surname, &p, &user_instance.Email, &user_instance.Password, &u)
	fmt.Println("user_repo: ", user_instance)

	return user_instance, err
}

func (r *DatabaseUserRepository) List(page int) ([]User, error) {
	var users []User
	pageSize := 10
	query := fmt.Sprintf("SELECT name, surname, date_birth, email, password, updated_at FROM bigproject.project_table LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)
	//query := fmt.Sprintf("SELECT * FROM bigproject.project_table LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)

	// Wykonaj zapytanie do bazy danych
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u, u2 string
	//var u2 string
	//var i string
	//var p string
	var user User
	for rows.Next() {
		if err := rows.Scan(&user.Name, &user.Surname, &u2, &user.Email, &user.Password, &u); err == nil {
			// err := rows.Scan(&user.Name, &user.Surname, &u2, &user.Email, &user.Password, &u); err == nil {
			user.Date_birth, err = time.Parse("2006-01-02 15:00:00", u2)
			if err != nil {
				fmt.Println("Error parsing Date_birth:", err)
			}

			user.Updated_at, err = time.Parse("2006-01-02 15:04:05.9999999", u)
			if err != nil {
				fmt.Println("Error parsing Updated_at:", err)
			}

		} else {
			fmt.Println(err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
