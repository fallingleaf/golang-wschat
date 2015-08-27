package tests

import (
    "testing"
    db "ocean/database"
)

func TestConnect(t *testing.T) {
    db.Connect()
    if db.Dbconn == nil {
        t.Error("connect fail")
    }
}

func TestQuery(t *testing.T) {
    user, err := db.GetUserById(1)
    if err != nil {
        t.Errorf("query fail %v", user)
    }
}