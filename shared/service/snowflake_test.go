package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateId(t *testing.T) {
	generator, err := NewDefaultGenerator()
	assert.NoError(t, err)

	id1 := generator.NextID()
	assert.Len(t, EncodeBase62(id1), 10)

	id2 := generator.NextID()
	assert.Len(t, EncodeBase62(id2), 10)

	assert.NotEqual(t, id1, id2)
}
