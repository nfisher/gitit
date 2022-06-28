package cmd

import (
	"github.com/go-git/go-git/v5"
	"os"
	"testing"
)

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

func InitialCommit(t *testing.T, repo *git.Repository) {
	w, err := os.Create("hello.txt")
	if err != nil {
		t.Fatalf("call=os.Create err=`%v`\n", err)
	}
	defer w.Close()

	_, err = w.WriteString("hello")
	if err != nil {
		t.Fatalf("call=w.WriteString err=`%v`\n", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("call=CreateBranch err=`%v`\n", err)
	}
	_, err = wt.Add("hello.txt")
	if err != nil {
		t.Fatalf("call=Add err=`%v`\n", err)
	}
	_, err = wt.Commit("Add hello.txt", &git.CommitOptions{})
	if err != nil {
		t.Fatalf("call=Commit err=`%v`\n", err)
	}
}
