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
	i.t.Helper()
	if i.a != b {
		i.t.Errorf("want %v, got %v\n", b, i.a)
	}
}
