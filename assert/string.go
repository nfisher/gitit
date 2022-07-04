package assert

import (
	"github.com/google/go-cmp/cmp/cmpopts"
	"strings"
	"testing"
)
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
	if diff := cmp.Diff(b, s.a, cmpopts.AcyclicTransformer("multiline", func(s string) []string {
		return strings.Split(s, "\n")
	})); diff != "" {
		s.t.Error("mismatch (-want +got):\n", diff)
	}
}
