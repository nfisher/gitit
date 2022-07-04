package assert

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"sort"
	"strings"
	"testing"
)

type gitremote struct {
	t *testing.T
	s string
}

func Remote(t *testing.T, s string) *gitremote {
	return &gitremote{t, s}
}

func (rem *gitremote) IncludesBranches(branches ...string) {
	rem.t.Helper()
	remoteBranches := rem.remoteBranches(branches)

	sort.Strings(remoteBranches)
	remoteList := strings.Join(remoteBranches, "\n")

	sort.Strings(branches)
	branchList := strings.Join(branches, "\n")

	if !strings.Contains(remoteList, branchList) {
		diff := cmp.Diff(branchList, remoteList, cmpopts.AcyclicTransformer("multiline", func(s string) []string {
			return strings.Split(s, "\n")
		}))
		rem.t.Errorf("remote-ls (-want +got):%s\n", diff)
	}
}

func (rem *gitremote) ExcludesBranches(branches ...string) {
	rem.t.Helper()
	remoteList := rem.remoteBranches(branches)

	var m = map[string]bool{}
	for _, r := range remoteList {
		m[r] = true
	}

	for _, b := range branches {
		if m[b] {
			rem.t.Errorf("remote-ls should not contain: %v\n", b)
		}
	}
}

func (rem *gitremote) remoteBranches(branches []string) []string {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{rem.s},
	})

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		rem.t.Fatalf("call=List err=`%v`\n", err)
	}
	var prefix = "refs/heads/"
	var a []string
	for _, r := range refs {
		s := r.Name().String()
		if strings.HasPrefix(s, prefix) {
			a = append(a, s[len(prefix):])
		}
	}

	return a
}

type gitrepo struct {
	t *testing.T
	g *git.Repository
}

func Repo(t *testing.T, g *git.Repository) *gitrepo {
	return &gitrepo{t, g}
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
