package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"html/template"
	"io"
	"log"
	"sort"
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

type Stack struct {
	Name     string
	Branch   string
	Branches []string
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

	branchName := ref.Name().String()
	parts := strings.Split(branchName, "/")
	if len(parts) == 4 {
		iter, err := repo.Branches()
		if err != nil {
			log.Printf("call=Branches err=`%v`\n", err)
			return ErrOutputWriter
		}
		var a []string
		err = iter.ForEach(func(reference *plumbing.Reference) error {
			s := reference.Name().String()
			p := strings.Split(s, "/")
			if len(p) == 4 && p[2] == parts[2] {
				a = append(a, p[3])
			}
			return nil
		})
		if err != nil {
			log.Printf("call=Branches err=`%v`\n", err)
			return ErrOutputWriter
		}
		sort.Strings(a)
		stack := &Stack{
			Name:     parts[2],
			Branch:   parts[3],
			Branches: a,
		}
		err = stackTpl.Execute(w, stack)
		if err != nil {
			// TODO: if w is stdout this is likely to fail as well.
			log.Printf("call=tpl.Execute err=`%v`\n", err)
			return ErrOutputWriter
		}
	} else if len(parts) == 3 {
		branch := parts[2]
		_, err = fmt.Fprintf(w, simpleBranch, branch)
		if err != nil {
			// TODO: if w is stdout this is likely to fail as well.
			log.Printf("call=Fprintf err=`%v`\n", err)
			return ErrOutputWriter
		}
	}

	return Success
}

var stackTpl = template.Must(template.New("stack").Parse(`In stack {{ .Name }}
On branch {{ .Name }}/{{ .Branch }}

Stack:
{{- range .Branches }}
    {{ . }}{{ end }}
`))

const simpleBranch = `Not in a stack
On branch %s
`
