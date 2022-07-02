package cmd_test

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"testing"
)

func WorkTree(t *testing.T, repo *git.Repository) *git.Worktree {
	t.Helper()
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("call=WorkTree err=`%v`", err)
	}
	return wt
}

func CreateBareDir(t *testing.T) func() {
	t.Helper()

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("call=Getwd err=`%v`\n", err)
	}

	dir := t.TempDir()
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("call=Chdir err=`%v`\n", err)
	}

	return func() {
		os.Chdir(pwd)
		os.RemoveAll(dir)
	}
}

func CreateRepo(t *testing.T) (*git.Repository, func()) {
	t.Helper()

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("call=Getwd err=`%v`\n", err)
	}

	dir := t.TempDir()
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("call=Chdir err=`%v`\n", err)
	}

	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("call=PlainInit err=`%v`\n", err)
	}

	return repo, func() {
		os.Chdir(pwd)
		os.RemoveAll(dir)
	}
}

func Commit(t *testing.T, wt *git.Worktree, files map[string]string, msg string) {
	t.Helper()
	for n, c := range files {
		AddFile(t, wt, n, c)
	}

	_, err := wt.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{Email: "nate@fisher.com", Name: "Nate Fisher"},
	})
	if err != nil {
		t.Fatalf("call=Commit err=`%v`\n", err)
	}
}

func AddFile(t *testing.T, wt *git.Worktree, filename, contents string) {
	t.Helper()
	w, err := os.Create(filename)
	if err != nil {
		t.Fatalf("call=os.Create err=`%v`\n", err)
	}
	defer w.Close()

	_, err = w.WriteString(contents)
	if err != nil {
		t.Fatalf("call=w.WriteString err=`%v`\n", err)
	}

	_, err = wt.Add(filename)
	if err != nil {
		t.Fatalf("call=Add err=`%v`\n", err)
	}
}

func InitialCommit(t *testing.T, repo *git.Repository) {
	t.Helper()
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("call=CreateBranch err=`%v`\n", err)
	}

	Commit(t, wt, map[string]string{
		".gitignore": "*.sw?",
	}, "Add .gitignore")
}

func CreateBranch(t *testing.T, repo *git.Repository, stack, branch string) {
	t.Helper()
	name := fmt.Sprintf("%s/%s", stack, branch)
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("call=Worktree err=`%v`\n", err)
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
		Create: true,
	})
	if err != nil {
		t.Fatalf("call=Checkout err=`%v`\n", err)
	}
}

func InitStack(t *testing.T, repo *git.Repository, stack, branch string) {
	t.Helper()
	CreateBranch(t, repo, stack, branch)
}

func SkipWIP(t *testing.T, runWip bool) {
	if !runWip {
		t.Skip("WIP")
	}
}

func CreateThreeLayerStack(t *testing.T, repo *git.Repository) {
	wt := WorkTree(t, repo)
	InitialCommit(t, repo)
	InitStack(t, repo, "kb3456", "001_migration")
	Commit(t, wt, map[string]string{"001_create.sql": "SELECT 1;"}, "Add 001_create.sql")
	InitStack(t, repo, "kb1234", "001_docs")
	Commit(t, wt, map[string]string{"README.md": "Hello world"}, "Add README.md")
	CreateBranch(t, repo, "kb1234", "002_api")
	Commit(t, wt, map[string]string{"api.js": "function api() {}"}, "Add api.js")
	CreateBranch(t, repo, "kb1234", "003_ui")
	Commit(t, wt, map[string]string{"ui.js": "function ui() {}"}, "Add ui.js")
}
