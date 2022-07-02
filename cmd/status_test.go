package cmd_test

import (
	"bytes"
	"github.com/nfisher/gitit/assert"
	. "github.com/nfisher/gitit/cmd"
	"io"
	"testing"
)

const simpleBranch = `Not in a stack
On branch master
`

func Test_status_on_branch_returns_success(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()
	InitialCommit(t, repo)

	var buf bytes.Buffer
	i := Exec(Flags{SubCommand: "status"}, &buf)
	assert.Int(t, i).Equals(Success)
	assert.String(t, string(buf.Bytes())).Equals(simpleBranch)
}

const smallStack = `In stack kb1234
On branch kb1234/003_ui

Stack:
    001_docs
    002_api
    003_ui
`

func Test_status_on_stack_returns_success(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	wt := WorkTree(t, repo)
	InitialCommit(t, repo)
	InitStack(t, repo, "kb3456", "001_migration")
	Commit(t, wt, map[string]string{"001_create.sql": "SELECT 1;"}, "Add 001_create.sql")
	InitStack(t, repo, "kb1234", "001_docs")
	Commit(t, wt, map[string]string{"README.md": "Hello world"}, "Add README.md")
	Branch(t, repo, "kb1234", "002_api")
	Commit(t, wt, map[string]string{"api.js": "function api() {}"}, "Add api.js")
	Branch(t, repo, "kb1234", "003_ui")
	Commit(t, wt, map[string]string{"ui.js": "function ui() {}"}, "Add ui.js")

	var buf bytes.Buffer
	i := Exec(Flags{SubCommand: "status"}, &buf)
	assert.Int(t, i).Equals(Success)
	assert.String(t, string(buf.Bytes())).Equals(smallStack)
}

func Test_status_returns_not_repository_with_empty_directory(t *testing.T) {
	tdclose := CreateBareDir(t)
	defer tdclose()

	i := Exec(Flags{SubCommand: "status"}, io.Discard)
	assert.Int(t, i).Equals(ErrNotRepository)
}
