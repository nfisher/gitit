package cmd_test

import (
	"bytes"
	"github.com/nfisher/gitit/assert"
	. "github.com/nfisher/gitit/cmd"
	"io"
	"testing"
)

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

func Test_init_returns_success_with_branch(t *testing.T) {
	i := Exec(Flags{SubCommand: "init", BranchName: "123/migration"}, io.Discard)
	assert.Int(t, i).Equals(Success)
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
    001_migration
  * 002_api
    003_ui
`

func Test_status_on_stack_returns_success(t *testing.T) {
	t.Skip("WIP")
	_, repoclose := CreateRepo(t)
	defer repoclose()

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
