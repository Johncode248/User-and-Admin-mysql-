# Dokumentacja Mini -> User-and-Admin-mysql- 

## 1. Wprowadzenie
Aplikacja User-and-Admin-mysql-  to prosta aplikacja webowa napisana w języku Go. Celem tej mini dokumentacji jest przedstawienie podstawowych informacji dotyczących aplikacji oraz instrukcji dotyczących uruchamiania jej w środowisku lokalnym przy użyciu Docker.

## 2. Dostępność Aplikacji 

Aplikacja będzie dostępna pod adresem `http://localhost:7943`.

## 3. Dockeryzacja Aplikacji
Aby uruchomić aplikację przy użyciu Dockera, wykonaj poniższe kroki:

a. Zainstaluj Dockera na swoim komputerze zgodnie z instrukcjami dostępnymi na oficjalnej stronie: [Docker Installation](https://docs.docker.com/get-docker/)

b. Zbuduj obraz Dockera: `docker-compose build.`

c. Uruchom kontener: `docker-compose up -d`

## 4. Zapytania przy użyciu cURL - USER
Aby wykonać proste zapytania do aplikacji przy użyciu cURL, użyj poniższych przykładów:

a. Tworzenie konta:
   ```bash
   curl -X POST -d '{"name": "John", "surname": "Doe", "date_birth": "2022-01-01T11:00:00Z", "email": "john.doe@example.com", "password": "example_password"}' http://localhost:7943/create_account

   ```
b. Logowanie:
   ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"name": "John", "password": "example_password"}' http://localhost:7943/start/login/user
  ```
c. Profil uzytkownika
   ```bash
  curl -X GET -H "Authorization: TOKEN" http://localhost:7943/login/user
   ```
d. Aktualizacja danych o user
   ```bash
  curl -X PUT -H "Content-Type: application/json" -H "Authorization: TOKEN" -d '{"Name":"new_name","Surname":"new_surname","Date_birth":"2022-01-01T11:00:00Z"}' http://localhost:7943/update/user

  ```

## 5. Zapytania przy użyciu cURL - ADMIN  
a. Logowanie Admina:
   ```bash
    curl -X POST -d '{"name":"Admin","password":"12345"}' http://localhost:7943/admin/login

   ```
b. Przegląd uzytkownikow:
   ```bash
   curl -X GET -H "Authorization: TOKEN" http://localhost:7943/admin/get_users
   ```
c. Usuwanie uzytkownika:
   ```bash
   curl -X DELETE -H "Authorization: TOKEN" http://localhost:7943/admin/delete/{id}
   ```      
   id -> User's email
