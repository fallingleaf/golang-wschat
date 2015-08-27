package database

import (
    "log"
)

const FIND_USER_Q = "SELECT * FROM tea_user WHERE id = $1"

func GetUserById(id int) (*User, error){
    var user User
    params := []interface{} {id}
    row := Dbconn.QueryRow(FIND_USER_Q, params)
    err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Token)
    if err != nil {
        log.Fatal(err)
        return &User{}, err
    }
    return &user, nil   
}