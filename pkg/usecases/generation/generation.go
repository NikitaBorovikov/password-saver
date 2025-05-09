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
	charSet := makeCharSet(ps.UseSpecialSymbols)

	password := make([]byte, ps.Length)
	for i := range password {
		password[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(password)
}

func makeCharSet(useSpecialSymbols bool) string {
	var charSet string
	if useSpecialSymbols {
		charSet = letters + digits + specials
	} else {
		charSet = letters + digits
	}

	return charSet
}
