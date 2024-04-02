package models

import (
    "database/sql"

    _ "github.com/lib/pq"
)



func ConnectDB()(*sql.DB, error){

    connectionString := "postgresql://docker:docker@localhost:5432/polls?sslmode=disable&search_path=public"

    db, err := sql.Open("postgres", connectionString)

    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db,nil
}
