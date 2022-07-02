package cmd_test

import (
	"github.com/nfisher/gitit/assert"
	. "github.com/nfisher/gitit/cmd"
	"io"
	"testing"
)

func Test_checkout_returns_success_with_known_branch_id(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	InitialCommit(t, repo)
	InitStack(t, repo, "kb3456", "001_migration")
	wt := WorkTree(t, repo)
	Commit(t, wt, map[string]string{"001_create.sql": "SELECT 1;"}, "Add 001_create.sql")

	i := Exec(Flags{SubCommand: "checkout", BranchName: "001"}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func Test_checkout_returns_not_found_with_invalid_stack(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	InitialCommit(t, repo)
	i := Exec(Flags{SubCommand: "checkout", BranchName: "002"}, io.Discard)
	assert.Int(t, i).Equals(ErrInvalidStack)
}

func Test_checkout_returns_missing_args_with_no_branch_id(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	InitialCommit(t, repo)
	i := Exec(Flags{SubCommand: "checkout"}, io.Discard)
	assert.Int(t, i).Equals(ErrMissingArguments)
}
