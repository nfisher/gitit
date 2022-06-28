package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"io"
	"log"
	"strings"
)

type Flags struct {
	SubCommand string
	BranchName string
}

const (
	Success = iota
	ErrHead
	ErrMissingArguments
	ErrMissingSubCommand
	ErrNotRepository
	ErrOutputWriter
)

func Exec(input Flags, w io.Writer) int {
	switch input.SubCommand {
	case "checkout":
		return Checkout(input)

	case "init":
		return Init(input)

	case "push":
		return Push(input)

	case "rebase":
		return Rebase(input)

	case "squash":
		return Squash(input)

	case "status":
		return Status(input, w)

	default:
		return ErrMissingSubCommand
	}
}

func Squash(_ Flags) int {
	return Success
}

func Rebase(_ Flags) int {
	return Success
}

func Push(_ Flags) int {
	return Success
}

func Checkout(input Flags) int {
	if input.BranchName == "" {
		return ErrMissingArguments
	}
	return Success
}

func Init(input Flags) int {
	if input.BranchName == "" {
		return ErrMissingArguments
	}
	return Success
}

func Status(_ Flags, w io.Writer) int {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})

	if err == git.ErrRepositoryNotExists {
		log.Printf("call=PlainOpen err=`%v`\n", err)
		return ErrNotRepository
	}

	ref, err := repo.Head()
	if err != nil {
		// TODO: how do we get here? Detached head?
		log.Printf("call=Head err=`%v`\n", err)
		return ErrHead
	}

	branch := splitRef(ref.Name().String())
	_, err = fmt.Fprintf(w, simpleBranch, branch)
	if err != nil {
		// TODO: if w is stdout this is likely to fail as well.
		log.Printf("call=Fprintf err=`%v`\n", err)
		return ErrOutputWriter
	}

	return Success
}

func splitRef(s string) string {
	a := strings.SplitN(s, "/", 3)
	if len(a) < 3 {
		return ""
	}
	return a[2]
}

const simpleBranch = `Not in a stack
On branch %s
`
