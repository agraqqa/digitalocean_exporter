package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Token(t *testing.T) {
	config := Config{
		DigitalOceanToken: "test-token-123",
	}

	token, err := config.Token()

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "test-token-123", token.AccessToken)
	assert.Equal(t, "", token.TokenType) // Should be empty for oauth2.Token
}

func TestConfig_DefaultValues(t *testing.T) {
	config := Config{}
	
	// Test that default values work as expected
	assert.Equal(t, "", config.DigitalOceanToken)
	assert.Equal(t, "", config.SpacesAccessKeyID)
	assert.Equal(t, "", config.SpacesAccessKeySecret)
	assert.Equal(t, 0, config.HTTPTimeout) // Should be 0 before parsing
	assert.Equal(t, "", config.WebAddr)
	assert.Equal(t, "", config.WebPath)
}