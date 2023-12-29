Edit in file sql.go this line: 
db, err := sql.Open("mysql", "user:password@tcp(localhost:xxxx)/database_bigproject")  with your data 
also you have to create schema "database_bigproject"
