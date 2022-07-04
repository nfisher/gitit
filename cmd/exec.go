package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"html/template"
	"io"
	"log"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
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
	ErrInvalidArgument
	ErrInvalidStack
	ErrUnknownBranch
	ErrNotRepository
	ErrOutputWriter
	ErrInvalidSequence
	ErrCreatingBranch
	ErrPushingStack
)

const (
	stackName   = 2
	stackBranch = 3
)

func Exec(input Flags, w io.Writer) int {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	switch input.SubCommand {
	case "branch":
		return Branch(input)

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

	case "version":
		return Version(w)

	default:
		usage(w)
		return ErrMissingSubCommand
	}
}

func Version(w io.Writer) int {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		os.Exit(1)
	}
	isDirty := false
	rev := "devel"
	for _, s := range buildInfo.Settings {
		switch s.Key {
		case "vcs.modified":
			isDirty = s.Value == "true"
		case "vcs.revision":
			rev = s.Value
		}
	}
	log.Printf("gitit@%v isDirty=%v\n", rev, isDirty)
	return Success
}

func Branch(input Flags) int {
	if input.BranchName == "" {
		log.Printf("call=BranchName err=`branch name is empty, must be specified`\n")
		return ErrMissingArguments
	}

	repo, wt, err := openWorkTree()
	if err != nil {
		log.Printf("call=openWorkTree err=`%v`\n", err)
		return ErrNotRepository
	}

	parts, err := headParts(repo)
	if err != nil {
		log.Printf("call=headParts err=`%v`\n", err)
		return ErrHead
	}

	if len(parts) != 4 {
		log.Printf("call=Split err=`want 4 parts, got %d`\n", len(parts))
		return ErrInvalidStack
	}

	iter, err := repo.Branches()
	if err != nil {
		log.Printf("call=Branches err=`%v`\n", err)
		return ErrOutputWriter
	}

	var a []string
	err = iter.ForEach(func(reference *plumbing.Reference) error {
		p := splitRef(reference)
		if len(p) == 4 && p[stackName] == parts[stackName] {
			a = append(a, p[stackBranch])
		}
		return nil
	})
	if err != nil {
		log.Printf("call=Branches err=`%v`\n", err)
		return ErrOutputWriter
	}
	sort.Strings(a)

	last := a[len(a)-1][:3]
	i, err := strconv.Atoi(last)
	if err != nil {
		log.Printf("call=Atoi err=`%v`\n", err)
		return ErrInvalidSequence
	}

	name := fmt.Sprintf("%s/%03d_%s", parts[stackName], i+1, input.BranchName)

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
		Create: true,
	})
	if err != nil {
		log.Printf("call=Checkout err=`%v`\n", err)
		return ErrCreatingBranch
	}

	fmt.Println("Created branch", name)
	return Success
}

func splitRef(reference *plumbing.Reference) []string {
	s := reference.Name().String()
	return strings.Split(s, "/")
}

func usage(w io.Writer) {
	w.Write([]byte(`usage: git stack <command> [<name>]

These are common Stack commands used in various situations:

start a new stack
   init       Create a new stack

examine the state
   status     Show the stack status
`))
}

func Squash(_ Flags) int {
	return Success
}

func Rebase(_ Flags) int {
	return Success
}

func Push(_ Flags) int {
	repo, _, err := openWorkTree()
	if err != nil {
		log.Printf("call=openWorkTree err=`%v`\n", err)
		return ErrNotRepository
	}

	parts, err := headParts(repo)
	if err != nil {
		log.Printf("call=headParts err=`%v`\n", err)
		return ErrHead
	}
	if len(parts) != 4 {
		log.Printf("call=Split err=`want 4 parts, got %d`\n", len(parts))
		return ErrInvalidStack
	}

	remotes, err := repo.Remotes()
	if err != nil {
		log.Printf("call=Remotes err=`%v`\n", err)
		return ErrInvalidStack
	}

	for _, r := range remotes {
		log.Printf("%v\n", r)
	}

	authcb, err := ssh.NewSSHAgentAuth("git")
	if err != nil {
		log.Printf("call=Split err=`want 4 parts, got %d`\n", len(parts))
		return ErrInvalidStack
	}

	spec := config.RefSpec(fmt.Sprintf("refs/heads/%[1]s/*:refs/heads/%[1]s/*", parts[stackName]))
	err = repo.Push(&git.PushOptions{
		Auth:       authcb,
		Progress:   os.Stdout,
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{spec},
	})
	if err != nil {
		log.Printf("call=Push spec=%v err=`%v`\n", spec, err)
		return ErrPushingStack
	}
	return Success
}

func Checkout(input Flags) int {
	if input.BranchName == "" {
		return ErrMissingArguments
	}

	repo, wt, err := openWorkTree()
	if err != nil {
		return ErrNotRepository
	}

	parts, err := headParts(repo)
	if err != nil {
		return ErrHead
	}
	if len(parts) != 4 {
		log.Printf("call=Split err=`want 4 parts, got %d`\n", len(parts))
		return ErrInvalidStack
	}

	iter, err := repo.Branches()
	if err != nil {
		log.Printf("call=Branches err=`%v`\n", err)
		return ErrOutputWriter
	}
	var target = ""
	err = iter.ForEach(func(reference *plumbing.Reference) error {
		p := splitRef(reference)
		if len(p) == 4 && p[stackName] == parts[stackName] && strings.HasPrefix(p[stackBranch], input.BranchName) {
			target = strings.Join(p[stackName:], "/")
		}
		return nil
	})
	if err != nil {
		log.Printf("call=ForEach err=`%v`\n", err)
		return ErrOutputWriter
	}

	if target == "" {
		log.Printf("call=ForEach err=`%v not found`\n", input.BranchName)
		return ErrUnknownBranch
	}

	err = wt.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(target), Keep: true})
	if err != nil {
		log.Printf("call=Checkout err=`%v`\n", err)
		return ErrUnknownBranch
	}

	return Success
}

func headParts(repo *git.Repository) ([]string, error) {
	ref, err := repo.Head()
	if err != nil {
		// TODO: how do we get here? Detached head?
		log.Printf("call=Head err=`%v`\n", err)
		return nil, err
	}
	return splitRef(ref), nil
}

func openWorkTree() (*git.Repository, *git.Worktree, error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err == git.ErrRepositoryNotExists {
		log.Printf("call=PlainOpen err=`%v`\n", err)
		return nil, nil, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		log.Printf("call=WorkTree err=`%v`\n", err)
		return nil, nil, err
	}

	return repo, wt, nil
}

func Init(input Flags) int {
	if input.BranchName == "" {
		return ErrMissingArguments
	}

	_, wt, err := openWorkTree()
	if err != nil {
		return ErrNotRepository
	}

	parts := strings.Split(input.BranchName, "/")
	if len(parts) != 2 {
		log.Printf("call=Split err=`%v`\n", err)
		return ErrInvalidArgument
	}
	name := fmt.Sprintf("%s/%03d_%s", parts[0], 1, parts[1])

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
		Create: true,
		Keep:   true,
	})
	if err != nil {
		log.Printf("call=Checkout err=`%v`\n", err)
		return ErrNotRepository
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

	parts, err := headParts(repo)
	if err != nil {
		return ErrHead
	}

	if len(parts) == 4 {
		iter, err := repo.Branches()
		if err != nil {
			log.Printf("call=Branches err=`%v`\n", err)
			return ErrOutputWriter
		}
		var a []string
		err = iter.ForEach(func(reference *plumbing.Reference) error {
			p := splitRef(reference)
			if isCurrentStack(p, parts) {
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

func isCurrentStack(p []string, cur []string) bool {
	return len(p) == 4 && p[stackName] == cur[stackName]
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
