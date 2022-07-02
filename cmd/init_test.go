package cmd_test

import (
	"github.com/nfisher/gitit/assert"
	. "github.com/nfisher/gitit/cmd"
	"io"
	"testing"
)

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
