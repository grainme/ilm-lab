package auth

import "github.com/alexedwards/argon2id"

func HashPassword(plainPassword string) (string, error) {
	hash, err := argon2id.CreateHash(plainPassword, argon2id.DefaultParams)
	return hash, err
}

func ComparePassword(plainPassword string, hashedPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plainPassword, hashedPassword)
	return match, err
}
