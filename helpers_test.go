package logicclipping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPseodoRandomString(t *testing.T) {
	s := GetPseodoRandomString()
	assert.NotEmpty(t, s)
	length := len(s)
	assert.GreaterOrEqual(t, length, 8)
}
