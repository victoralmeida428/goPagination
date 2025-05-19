package pagination

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestEnvelope(t *testing.T) {
	env := Envelope{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}
	
	assert.Equal(t, "value1", env["key1"])
	assert.Equal(t, 42, env["key2"])
	assert.Equal(t, true, env["key3"])
}
