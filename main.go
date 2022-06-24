package main

import (
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"log"
	"os"
	"strings"
)

func Execute() {
	status := flag.NewFlagSet("status", flag.ExitOnError)

	cmds := []*flag.FlagSet{status}

	var subcmd string
	flag.Parse()
	arg1 := flag.Arg(0)

	for _, f := range cmds {
		if arg1 != f.Name() {
			continue
		}
		subcmd = arg1

		err := f.Parse(flag.Args())
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}

	switch subcmd {
	case "status":
		Status()
	default:
		fmt.Printf("subcmd=%s err=`Unknown subcommand`\n", subcmd)
		os.Exit(1)
	}
}

func Status() error {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		log.Printf("call=PlainOpen err=`%v`\n", err)
		return err
	}

	ref, err := repo.Head()
	if err != nil {
		log.Printf("call=Head err=`%v`\n", err)
		return err
	}

	parts := strings.SplitN(ref.Name().String(), "/", 3)
	branch := parts[2]
	stack := "fix123"

	iter, err := repo.Branches()
	if err != nil {
		log.Printf("call=Branches err=`%v`\n", err)
		return err
	}

	iter.ForEach(func(reference *plumbing.Reference) error {
		reference.Target().IsBranch()
		fmt.Println(reference.Target().IsBranch(), reference.Type(), reference.Name())
		return nil
	})

	fmt.Printf(`In stack %s
On branch %s

Stack:
    001_migration   🚢 3629a61
  * 002_api         ✅ https://github.com/nfisher/gitit/pulls/110781
    003_ui          ❌ https://github.com/nfisher/gitit/pulls/110779

`, stack, branch)

	return nil
}

func main() {
	Execute()
}
