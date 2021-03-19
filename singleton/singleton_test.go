package singleton

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnce(t *testing.T) {
	a := 0
	fn := func() {
		a++
	}
	Once(11, fn)
	assert.Equal(t, 1, a)

	Once(11, fn)
	Once(11, fn)
	Once(11, fn)
	assert.Equal(t, 1, a)

	Cancel(11)
	a = 0
	Once(11, fn)
	assert.Equal(t, 1, a)
}
