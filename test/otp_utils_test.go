package test

import (
	"github.com/djfemz/rave/partybank-app/security/otp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateOtp(t *testing.T) {
	password := otp.GenerateOtp()
	assert.NotNil(t, password)
	assert.NotEmpty(t, password)
}
