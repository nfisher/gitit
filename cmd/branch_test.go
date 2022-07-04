package cmd_test

import (
	"github.com/nfisher/gitit/assert"
	. "github.com/nfisher/gitit/cmd"
	"io"
	"testing"
)

func Test_branch_outside_repo_should_fail(t *testing.T) {
	tdclose := CreateBareDir(t)
	defer tdclose()

	i := Exec(Flags{SubCommand: "branch", BranchName: "ml_fairy"}, io.Discard)
	assert.Int(t, i).Equals(ErrNotRepository)
}

func Test_branch_on_existing_stack_returns_success(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	CreateThreeLayerStack(t, repo)

	i := Exec(Flags{SubCommand: "branch", BranchName: "ml_fairy"}, io.Discard)
	assert.Int(t, i).Equals(Success)
	assert.Repo(t, repo).Branch("kb1234/004_ml_fairy")
}

func Test_branch_on_invalid_stack_fails(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	InitialCommit(t, repo)

	i := Exec(Flags{SubCommand: "branch", BranchName: "ml_fairy"}, io.Discard)
	assert.Int(t, i).Equals(ErrInvalidStack)
	assert.Repo(t, repo).Branch("master")
}

func Test_branch_fails_with_no_name_specified(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	CreateThreeLayerStack(t, repo)

	i := Exec(Flags{SubCommand: "branch"}, io.Discard)
	assert.Int(t, i).Equals(ErrMissingArguments)
	assert.Repo(t, repo).Branch("kb1234/003_ui")
}

func Test_branch_returns_success_with_dirty_branch(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	CreateThreeLayerStack(t, repo)
	CreateFile(t, ".gitignore", "*.sw?\n.idea")

	i := Exec(Flags{SubCommand: "branch", BranchName: "update_ignore"}, io.Discard)

	assert.Int(t, i).Equals(Success)
	assert.Repo(t, repo).Branch("kb1234/004_update_ignore")
	assert.Exists(t, ".gitignore")
}
