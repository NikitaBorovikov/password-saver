package generation

import (
	"math/rand"
	"password-saver/pkg/dto"
)

const (
	letters  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits   = "0123456789"
	specials = "!@#$%^&*()_+"
)

func GenNewPassword(ps *dto.GeneratePasswordRequest) string {
	var charSet string

	if ps.UseSpecialSymbols {
		charSet = letters + digits + specials
	} else {
		charSet = letters + digits
	}

	password := make([]byte, ps.Length)
	for i := range password {
		password[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(password)
}
