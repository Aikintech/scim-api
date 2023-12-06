package utils

import (
	"errors"
	"github.com/matthewhartstonge/argon2"
)

func MakePasswordHash(password string) (string, error) {
	argon := argon2.DefaultConfig()

	encoded, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", errors.New("Error hashing password")
	}

	return string(encoded), nil
}

func VerifyPasswordHash(password string, hashed string) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(hashed))
	if err != nil {
		return false, errors.New("Error verifying password")
	}

	return ok, nil
}
