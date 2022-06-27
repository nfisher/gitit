package assert

import "testing"

func Int(t *testing.T, i int) *integer {
	return &integer{i, t}
}

type integer struct {
	a int
	t *testing.T
}

func (i *integer) Equals(b int) {
	if i.a != b {
		i.t.Helper()
		i.t.Errorf("expected %v, got %v\n", b, i.a)
	}
}
