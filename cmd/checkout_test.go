package cmd_test

import (
	"github.com/nfisher/gitit/assert"
	. "github.com/nfisher/gitit/cmd"
	"io"
	"testing"
)

func Test_checkout_returns_unknown_branch_with_absent_branch_id(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	CreateThreeLayerStack(t, repo)

	i := Exec(Flags{SubCommand: "checkout", BranchName: "004"}, io.Discard)
	assert.Int(t, i).Equals(ErrUnknownBranch)
	assert.Repo(t, repo).Branch("kb1234/003_ui")
}

func Test_checkout_returns_success_with_known_branch_id(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	CreateThreeLayerStack(t, repo)

	i := Exec(Flags{SubCommand: "checkout", BranchName: "001"}, io.Discard)
	assert.Int(t, i).Equals(Success)
	assert.Repo(t, repo).Branch("kb1234/001_docs")
}

func Test_checkout_returns_not_found_with_invalid_stack(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	InitialCommit(t, repo)
	i := Exec(Flags{SubCommand: "checkout", BranchName: "002"}, io.Discard)
	assert.Int(t, i).Equals(ErrInvalidStack)
	assert.Repo(t, repo).Branch("master")
}

func Test_checkout_returns_missing_args_with_no_branch_id(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	InitialCommit(t, repo)
	i := Exec(Flags{SubCommand: "checkout"}, io.Discard)
	assert.Int(t, i).Equals(ErrMissingArguments)
}
