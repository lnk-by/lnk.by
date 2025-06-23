package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func GenerateId(t *testing.T) {
	generator, err := NewDefaultGenerator()
	assert.Nil(t, err)
	id1 := generator.NextID()
	assert.Len(t, id1, 10)

	id2 := generator.NextID()
	assert.Len(t, id2, 10)

	assert.NotEqual(t, id1, id2)
}
