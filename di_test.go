package depinject

import (
	"testing"
)

func TestDI(t *testing.T) {
	type X struct {
		A string
	}

	var target struct {
		X  *X
		X2 *X `depinject:"skip"`
	}

	x := &X{}

	var di DI
	di.MustRegisterComponent(x)
	di.MustInject(&target)

	if target.X != x {
		t.Fatal("expected target.X to be equal to x")
	}

	if target.X2 != nil {
		t.Fatal("expected target.X2 to be nil")
	}
}
