package main

import (
	"golang.org/x/crypto/bcrypt"
)

// Hashes the given password and returns the hashed value
func hashPassword(password string) (string, error) {
	// Generate a salted hash for the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compares the provided password with the hashed password
func comparePasswords(password, hashedPassword string) error {
	// Compare the password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
