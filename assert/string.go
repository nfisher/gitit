package assert

import "testing"

func String(t *testing.T, s string) *str {
	return &str{t, s}
}

type str struct {
	t *testing.T
	a string
}

func (s *str) Equals(b string) {
	s.t.Helper()
	if s.a != b {
		s.t.Helper()
		s.t.Errorf("want `%s`, got `%s`\n", b, s.a)
	}
}
