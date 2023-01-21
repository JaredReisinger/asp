package asp

import (
	"testing"
)

func TestWithEnvPrefix(t *testing.T) {
	a := &aspBase{}

	err := WithEnvPrefix("PREFIX")(a)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if a.envPrefix != "PREFIX" {
		t.Fail()
	}
}

func TestWithConfigFlag(t *testing.T) {
	a := &aspBase{}

	err := WithConfigFlag(a)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if !a.withConfigFlag {
		t.Fail()
	}
}

func TestWithoutConfigFlag(t *testing.T) {
	a := &aspBase{withConfigFlag: true}

	err := WithoutConfigFlag(a)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if a.withConfigFlag {
		t.Fail()
	}
}
