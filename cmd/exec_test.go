package cmd_test

import (
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
