package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type shit struct {
	f1 string
	f2 int
}

func (s *shit) fieldsPtrs() []any { return []any{&s.f1, &s.f2} }

func TestInst(t *testing.T) {
	shitPtr := inst[*shit]()
	assert.NotNil(t, shitPtr)
	assert.NotNil(t, shitPtr.fieldsPtrs())
}
