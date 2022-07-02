package assert

import (
	"github.com/go-git/go-git/v5"
	"strings"
	"testing"
)

func Repo(t *testing.T, g *git.Repository) *gitrepo {
	return &gitrepo{t, g}
}

type gitrepo struct {
	t *testing.T
	g *git.Repository
}

func (gr *gitrepo) Branch(b string) {
	gr.t.Helper()
	head, err := gr.g.Head()
	if err != nil {
		gr.t.Fatalf("call=Head err=`%v`\n", err)
	}

	a := strings.Replace(head.Name().String(), "refs/heads/", "", 1)
	if a != b {
		gr.t.Fatalf("want %v, got %v\n", b, a)
	}
}
