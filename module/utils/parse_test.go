package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseUUID(t *testing.T) {
	parentURL, uuid := ParseUUID("/api/test/uuid")
	assert.Equal(t, parentURL, "/api/test")
	assert.Equal(t, uuid, "uuid")
}
