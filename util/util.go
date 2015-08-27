package util

import (
    "math/rand"
)

const (
    length = 32
    letters = "abcdefghijkmnopqrstuvxyzABCDEFGHIJKMNOPQRSTUVXYZ0123456789"
)

func RandomString(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func RandomID() string {
    return RandomString(length)
}
