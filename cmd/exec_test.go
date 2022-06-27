package cmd_test

import (
	"github.com/go-git/go-git/v5"
	"github.com/nfisher/gitit/assert"
	. "github.com/nfisher/gitit/cmd"
	"io"
	"os"
	"testing"
)

func Test_no_args_returns_missing_subcommand(t *testing.T) {
	i := Exec(Flags{}, io.Discard)
	assert.Int(t, i).Equals(MissingSubCommand)
}

func Test_checkout_returns_success_with_branch_id(t *testing.T) {
	i := Exec(Flags{SubCommand: Checkout, BranchName: "001"}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func Test_checkout_returns_missing_args_with_no_branch_id(t *testing.T) {
	i := Exec(Flags{SubCommand: Checkout}, io.Discard)
	assert.Int(t, i).Equals(MissingArguments)
}

func Test_init_returns_missing_arguments_without_branch_arg(t *testing.T) {
	i := Exec(Flags{SubCommand: Init}, io.Discard)
	assert.Int(t, i).Equals(MissingArguments)
}

func Test_init_returns_success_with_branch(t *testing.T) {
	i := Exec(Flags{SubCommand: Init, BranchName: "123/migration"}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func Test_push_returns_success(t *testing.T) {
	i := Exec(Flags{SubCommand: Push}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func Test_rebase_returns_success(t *testing.T) {
	i := Exec(Flags{SubCommand: Rebase}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func Test_squash_returns_success(t *testing.T) {
	i := Exec(Flags{SubCommand: Squash}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func Test_status_returns_success(t *testing.T) {
	i := Exec(Flags{SubCommand: Status}, io.Discard)
	assert.Int(t, i).Equals(Success)
}

func CreateRepo(t *testing.T) (*git.Repository, func()) {
	t.Helper()
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Errorf("call=PlainInit err=`%v`\n", err)
	}

	return repo, func() {
		os.RemoveAll(dir)
	}
}
