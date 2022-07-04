package cmd_test

import (
	"github.com/nfisher/gitit/assert"
	. "github.com/nfisher/gitit/cmd"
	"io"
	"testing"
)

func Test_push_returns_success(t *testing.T) {
	server, srvclose := LaunchServer(t)
	defer srvclose()

	repo, repoclose := CreateRepo(t)
	defer repoclose()

	CreateThreeLayerStack(t, repo)
	CreateRemote(t, repo, server)

	i := Exec(Flags{SubCommand: "push"}, io.Discard)
	assert.Int(t, i).Equals(Success)
	assert.Remote(t, server.Address()).IncludesBranches(
		"kb1234/001_docs",
		"kb1234/002_api",
		"kb1234/003_ui")
	assert.Remote(t, server.Address()).ExcludesBranches(
		"kb3456/001_migration")
}
