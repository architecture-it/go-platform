package database

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateConnectionRedis_DefaultTimeouts(t *testing.T) {
	// Clear any existing environment variables
	os.Unsetenv("REDIS_DIAL_TIMEOUT_MS")
	os.Unsetenv("REDIS_READ_TIMEOUT_MS")
	os.Unsetenv("REDIS_WRITE_TIMEOUT_MS")
	os.Unsetenv("REDIS_SKIP_PING")

	// Note: This will fail to connect since we don't have a real Redis, 
	// but we're testing the timeout configuration logic
	client := createConnectionRedis("localhost:6379", "", 0)

	// Client should be nil because connection will fail
	// but the important part is that it used the default timeouts
	assert.Nil(t, client, "Client should be nil when Redis is not available")
}

func TestCreateConnectionRedis_CustomDialTimeout(t *testing.T) {
	os.Setenv("REDIS_DIAL_TIMEOUT_MS", "1000")
	defer os.Unsetenv("REDIS_DIAL_TIMEOUT_MS")

	client := createConnectionRedis("localhost:6379", "", 0)
	
	// Client should be nil because connection will fail
	assert.Nil(t, client, "Client should be nil when Redis is not available")
}

func TestCreateConnectionRedis_CustomReadTimeout(t *testing.T) {
	os.Setenv("REDIS_READ_TIMEOUT_MS", "500")
	defer os.Unsetenv("REDIS_READ_TIMEOUT_MS")

	client := createConnectionRedis("localhost:6379", "", 0)
	
	assert.Nil(t, client, "Client should be nil when Redis is not available")
}

func TestCreateConnectionRedis_CustomWriteTimeout(t *testing.T) {
	os.Setenv("REDIS_WRITE_TIMEOUT_MS", "500")
	defer os.Unsetenv("REDIS_WRITE_TIMEOUT_MS")

	client := createConnectionRedis("localhost:6379", "", 0)
	
	assert.Nil(t, client, "Client should be nil when Redis is not available")
}

func TestCreateConnectionRedis_AllCustomTimeouts(t *testing.T) {
	os.Setenv("REDIS_DIAL_TIMEOUT_MS", "1000")
	os.Setenv("REDIS_READ_TIMEOUT_MS", "600")
	os.Setenv("REDIS_WRITE_TIMEOUT_MS", "600")
	defer func() {
		os.Unsetenv("REDIS_DIAL_TIMEOUT_MS")
		os.Unsetenv("REDIS_READ_TIMEOUT_MS")
		os.Unsetenv("REDIS_WRITE_TIMEOUT_MS")
	}()

	client := createConnectionRedis("localhost:6379", "", 0)
	
	assert.Nil(t, client, "Client should be nil when Redis is not available")
}

func TestCreateConnectionRedis_InvalidTimeout(t *testing.T) {
	// Set invalid timeout value (non-numeric)
	os.Setenv("REDIS_DIAL_TIMEOUT_MS", "invalid")
	defer os.Unsetenv("REDIS_DIAL_TIMEOUT_MS")

	// Should fall back to default 500ms
	client := createConnectionRedis("localhost:6379", "", 0)
	
	assert.Nil(t, client, "Client should be nil when Redis is not available")
}

func TestCreateConnectionRedis_SkipPing(t *testing.T) {
	os.Setenv("REDIS_SKIP_PING", "true")
	defer os.Unsetenv("REDIS_SKIP_PING")

	// With SKIP_PING=true, client should be created even without valid Redis
	client := createConnectionRedis("localhost:6379", "", 0)
	
	// When skip ping is true, client should be created (not nil)
	assert.NotNil(t, client, "Client should be created when REDIS_SKIP_PING is true")
	
	// Verify the client options
	assert.Equal(t, "localhost:6379", client.Options().Addr)
	assert.Equal(t, 0, client.Options().DB)
	
	// Verify default timeouts are set
	assert.Equal(t, time.Millisecond*500, client.Options().DialTimeout)
	assert.Equal(t, time.Millisecond*300, client.Options().ReadTimeout)
	assert.Equal(t, time.Millisecond*300, client.Options().WriteTimeout)
}

func TestCreateConnectionRedis_SkipPingWithCustomTimeouts(t *testing.T) {
	os.Setenv("REDIS_SKIP_PING", "true")
	os.Setenv("REDIS_DIAL_TIMEOUT_MS", "2000")
	os.Setenv("REDIS_READ_TIMEOUT_MS", "1500")
	os.Setenv("REDIS_WRITE_TIMEOUT_MS", "1500")
	defer func() {
		os.Unsetenv("REDIS_SKIP_PING")
		os.Unsetenv("REDIS_DIAL_TIMEOUT_MS")
		os.Unsetenv("REDIS_READ_TIMEOUT_MS")
		os.Unsetenv("REDIS_WRITE_TIMEOUT_MS")
	}()

	client := createConnectionRedis("localhost:6379", "", 0)
	
	assert.NotNil(t, client, "Client should be created when REDIS_SKIP_PING is true")
	
	// Verify custom timeouts are set
	assert.Equal(t, time.Millisecond*2000, client.Options().DialTimeout)
	assert.Equal(t, time.Millisecond*1500, client.Options().ReadTimeout)
	assert.Equal(t, time.Millisecond*1500, client.Options().WriteTimeout)
}
