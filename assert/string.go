package assert

import "testing"
import "github.com/google/go-cmp/cmp"

func String(t *testing.T, s string) *str {
	return &str{t, s}
}

type str struct {
	t *testing.T
	a string
}

func (s *str) Equals(b string) {
	s.t.Helper()
	if diff := cmp.Diff(b, s.a); diff != "" {
		s.t.Error("mismatch (-want +got):\n", diff)
	}
}
