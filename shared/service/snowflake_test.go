package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateId(t *testing.T) {
	generator, err := NewDefaultGenerator()
	assert.NoError(t, err)

	id1 := EncodeBase62(generator.NextID())
	assert.Len(t, id1, 10)

	id2 := EncodeBase62(generator.NextID())
	assert.Len(t, id2, 10)

	assert.NotEqual(t, id1, id2)
}

func TestGenerateBase62Id(t *testing.T) {
	generator, err := NewDefaultGenerator()
	assert.NoError(t, err)

	id1 := generator.NextBase62ID()
	assert.Len(t, id1, 10)

	id2 := generator.NextBase62ID()
	assert.Len(t, id2, 10)

	assert.NotEqual(t, id1, id2)
}
