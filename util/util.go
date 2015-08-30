package util

import (
    "time"
    "math/rand"
)

const (
    length = 32
    letters = "abcdefghijkmnopqrstuvxyzABCDEFGHIJKMNOPQRSTUVXYZ0123456789"
)

func RandomString(n int) string {
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func RandomID() string {
    return RandomString(length)
}
