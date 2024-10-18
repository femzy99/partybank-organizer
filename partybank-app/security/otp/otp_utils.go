package otp

import (
	"crypto/rand"
	"io"
	"time"
)

const (
	OTP_EXPIRES_AT   = 5
	DEFAULT_OTP_SIZE = 6
)

type OneTimePassword struct {
	Code      string
	ExpiresAt time.Time
}

func GenerateOtp() *OneTimePassword {
	return &OneTimePassword{
		buildPassword(DEFAULT_OTP_SIZE),
		time.Now().Add(OTP_EXPIRES_AT * time.Minute),
	}
}

func buildPassword(passwordLength int) string {
	store := make([]byte, passwordLength)
	numberOfBytesRead, err := io.ReadAtLeast(rand.Reader, store, passwordLength)
	if numberOfBytesRead != passwordLength {
		panic(err)
	}
	for index := 0; index < len(store); index++ {
		store[index] = table[int(store[index])%len(table)]
	}
	return string(store)
}

var table = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
