package asp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithEnvPrefix(t *testing.T) {
	a := &aspBase{}

	err := WithEnvPrefix("PREFIX")(a)
	assert.NoError(t, err)
	assert.Equal(t, "PREFIX", a.envPrefix)
}

func TestWithConfigFlag(t *testing.T) {
	a := &aspBase{}

	err := WithConfigFlag(a)
	assert.NoError(t, err)
	assert.True(t, a.withConfigFlag)
}

func TestWithoutConfigFlag(t *testing.T) {
	a := &aspBase{withConfigFlag: true}

	err := WithoutConfigFlag(a)
	assert.NoError(t, err)
	assert.False(t, a.withConfigFlag)
}

func TestWithDecodeHook(t *testing.T) {
	a := &aspBase{}

	dummyHook := "HOOK"

	err := WithDecodeHook(dummyHook)(a)
	assert.NoError(t, err)
	assert.Equal(t, dummyHook, a.decodeHook)

}
