package cmd_test

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"net"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func CreateFile(t *testing.T, filename, contents string) {
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

type server struct {
	Port int
	Root string
	Repo *git.Repository
}

func (s *server) Address() string {
	return fmt.Sprintf("http://localhost:%d", s.Port)
}

func LaunchServer(t *testing.T) (*server, func()) {
	// inspired by go-git tests.
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("call=Listen err=`%v`\n", err)
	}
	d := t.TempDir()

	cmd := exec.Command("git", "--exec-path")
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("call=CombinedOutput err=`%v`\n", err)
	}

	backend := filepath.Join(strings.TrimSpace(string(b)), "git-http-backend")

	srv := http.Server{
		Handler: &cgi.Handler{
			Path: backend,
			Env:  []string{"GIT_HTTP_EXPORT_ALL=true", fmt.Sprintf("GIT_PROJECT_ROOT=%s", d)},
		},
	}

	repo, err := git.PlainInit(d, true)
	if err != nil {
		t.Fatalf("call=PlainInit err=`%v`\n", err)
	}

	err = os.WriteFile(filepath.Join(d, "config"), []byte(configContents), os.ModePerm)
	if err != nil {
		t.Fatalf("call=WriteFile err=`%v`\n", err)
	}

	s := &server{
		Root: d,
		Port: l.Addr().(*net.TCPAddr).Port,
		Repo: repo,
	}

	go func() {
		srv.Serve(l)
	}()

	return s, func() {
		srv.Shutdown(context.Background())
	}
}

var configContents = `[core]
  bare = true
[http]
  receivepack = true
`

func CreateRemote(t *testing.T, repo *git.Repository, s *server) {
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{fmt.Sprintf("http://localhost:%d", s.Port)},
	})
	if err != nil {
		t.Fatalf("call=CreateRemote err=`%v`\n", err)
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
