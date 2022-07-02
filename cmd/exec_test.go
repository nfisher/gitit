package cmd_test

import (
	"bytes"
	"flag"
	"github.com/nfisher/gitit/assert"
	. "github.com/nfisher/gitit/cmd"
	"io"
	"testing"
)

var runWip = flag.Bool("runwip", false, "Run WIP tests")

func Test_no_args_returns_missing_subcommand(t *testing.T) {
	i := Exec(Flags{}, io.Discard)
	assert.Int(t, i).Equals(ErrMissingSubCommand)
}

func Test_checkout_returns_success_with_branch_id(t *testing.T) {
	i := Exec(Flags{SubCommand: "checkout", BranchName: "001"}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func Test_checkout_returns_missing_args_with_no_branch_id(t *testing.T) {
	i := Exec(Flags{SubCommand: "checkout"}, io.Discard)
	assert.Int(t, i).Equals(ErrMissingArguments)
}

func Test_init_returns_missing_arguments_without_branch_arg(t *testing.T) {
	i := Exec(Flags{SubCommand: "init"}, io.Discard)
	assert.Int(t, i).Equals(ErrMissingArguments)
}

func Test_init_returns_success_with_branch_specified(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()
	InitialCommit(t, repo)

	i := Exec(Flags{SubCommand: "init", BranchName: "123/migration"}, io.Discard)
	assert.Int(t, i).Equals(Success)
	assert.Repo(t, repo).Branch("123/001_migration")
}

func Test_init_returns_failure_with_invalid_branch_specification(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()
	InitialCommit(t, repo)

	i := Exec(Flags{SubCommand: "init", BranchName: "migration"}, io.Discard)
	assert.Int(t, i).Equals(ErrInvalidArgument)
}

func Test_push_returns_success(t *testing.T) {
	i := Exec(Flags{SubCommand: "push"}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func Test_rebase_returns_success(t *testing.T) {
	i := Exec(Flags{SubCommand: "rebase"}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func Test_squash_returns_success(t *testing.T) {
	i := Exec(Flags{SubCommand: "squash"}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

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

	InitialCommit(t, repo)
	wt, err := repo.Worktree()
	if err != nil {
		t.Errorf("call=WorkTree err=`%v`", err)
	}
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
