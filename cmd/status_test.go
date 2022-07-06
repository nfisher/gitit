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
	// Create 2 workspaces and interleave a timeline of pushes as follows:
	// 1. A - initial push 001,002,003.
	// 2. B - pull 001, 002, 003.
	// 3. A - push change 002.
	// 4. B - push change 003.

	server, srvclose := LaunchServer(t)
	defer srvclose()

	repo1, r1close := CreateRepo(t)
	defer r1close()

	wt1 := WorkTree(t, repo1)

	CreateThreeLayerStack(t, repo1)
	CreateRemote(t, repo1, server)

	PushBranch(t, repo1, "kb1234/001_docs")
	PushBranch(t, repo1, "kb1234/002_api")
	PushBranch(t, repo1, "kb1234/003_ui")

	CheckoutBranch(t, wt1, "kb1234/002_api")
	Commit(t, wt1, map[string]string{"api.js": "function api() { return true; }"}, "Update api.js")

	repo2, r2close := CloneRepo(t, server, "kb1234/003_ui")
	defer r2close()

	FetchRefs(t, repo2, "kb1234/*")

	wt2 := WorkTree(t, repo2)
	Commit(t, wt2, map[string]string{"ui.js": "function ui() { return true; }"}, "Update ui.js")

	var buf2 bytes.Buffer
	i := Exec(Flags{SubCommand: "status"}, &buf2)
	assert.Int(t, i).Equals(Success)
	assert.String(t, string(buf2.Bytes())).Equals(`In stack kb1234
On branch kb1234/003_ui
Remote origin

Local Stack (+ ahead, = same, ∇ diverged):
    (=) 001_docs
    (=) 002_api
    (+) 003_ui
`)
	PushBranch(t, repo2, "kb1234/003_ui")

	Chdir(t, wt1)

	var buf1 bytes.Buffer
	i = Exec(Flags{SubCommand: "status"}, &buf1)
	assert.Int(t, i).Equals(Success)
	assert.String(t, string(buf1.Bytes())).Equals(`In stack kb1234
On branch kb1234/002_api
Remote origin

Local Stack (+ ahead, = same, ∇ diverged):
    (=) 001_docs
    (+) 002_api
    (∇) 003_ui
`)
}
