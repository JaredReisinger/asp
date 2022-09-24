package asp

import (
	"testing"
)

func TestWithEnvPrefix(t *testing.T) {
	a := &asp[TestConfig]{}

	err := WithEnvPrefix[TestConfig]("PREFIX")(a)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if a.envPrefix != "PREFIX" {
		t.Fail()
	}
}

func TestWithConfigFlag(t *testing.T) {
	a := &asp[TestConfig]{}

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
	a := &asp[TestConfig]{withConfigFlag: true}

	err := WithoutConfigFlag(a)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if a.withConfigFlag {
		t.Fail()
	}
}
