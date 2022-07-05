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

Local Stack:
    001_docs
    002_api
    003_ui
`

func Test_status_on_stack_returns_success(t *testing.T) {
	repo, repoclose := CreateRepo(t)
	defer repoclose()

	CreateThreeLayerStack(t, repo)

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

func Test_status_with_default_remote_successful(t *testing.T) {
	server, srvclose := LaunchServer(t)
	defer srvclose()

	repo, repoclose := CreateRepo(t)
	defer repoclose()

	CreateThreeLayerStack(t, repo)
	CreateRemote(t, repo, server)

	PushBranch(t, repo, "kb1234/001_docs")
	PushBranch(t, repo, "kb1234/002_api")

	var buf bytes.Buffer
	i := Exec(Flags{SubCommand: "status"}, &buf)
	assert.Int(t, i).Equals(Success)
	assert.String(t, string(buf.Bytes())).Equals(`In stack kb1234
On branch kb1234/003_ui
Remote origin

Local Stack (+ ahead, = same, ∇ diverged):
    (=) 001_docs
    (=) 002_api
    (+) 003_ui
`)
}

func Test_status_with_diverged_remote_successful(t *testing.T) {
	SkipWIP(t, *runWip)
	server, srvclose := LaunchServer(t)
	defer srvclose()

	repo, repoclose := CreateRepo(t)
	defer repoclose()

	// TODO: Create 2 repos and interleave a timeline of pushes as follows:
	// 1. A - initial push 001,002,003.
	// 2. B - pull 001, 002, 003.
	// 3. A - push change 002.
	// 4. B - push change 003.

	CreateThreeLayerStack(t, repo)
	CreateRemote(t, repo, server)

	PushBranch(t, repo, "kb1234/001_docs")
	PushBranch(t, repo, "kb1234/002_api")

	var buf bytes.Buffer
	i := Exec(Flags{SubCommand: "status"}, &buf)
	assert.Int(t, i).Equals(Success)
	assert.String(t, string(buf.Bytes())).Equals(`In stack kb1234
On branch kb1234/003_ui
Remote origin

Local Stack (+ ahead, = same, ∇ diverged):
    (=) 001_docs
    (+) 002_api
    (∇) 003_ui
`)
}
