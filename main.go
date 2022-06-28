package main

import (
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"log"
	"os"
	"sort"
	"strings"
)

func Execute() {
	status := flag.NewFlagSet("status", flag.ExitOnError)
	flags := []*flag.FlagSet{
		status,
	}
	flag.Parse()

	arg1 := flag.Arg(0)

	var sub string
	for _, f := range flags {
		if arg1 != f.Name() {
			continue
		}
		sub = arg1

		err := f.Parse(flag.Args())
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}

	var err error
	switch sub {
	case "status":
		err = Status()
	default:
		fmt.Printf("sub=%s err=`Unknown subcommand`\n", sub)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("sub=%s err=`%v`\n", sub, err)
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

	parts := strings.SplitN(ref.Name().String(), "/", 4)

	var prefix string
	var stack = "<< Undefined >>"
	var branch string
	if len(parts) > 3 {
		stack = parts[2]
		branch = strings.Join(parts[2:], "/")
		prefix = strings.Join(parts[:3], "/")
	} else {
		branch = parts[2]
	}

	iter, err := repo.Branches()
	if err != nil {
		log.Printf("call=Branches err=`%v`\n", err)
		return err
	}

	var stackBranches []string
	if prefix != "" {
		err := iter.ForEach(func(reference *plumbing.Reference) error {
			if strings.HasPrefix(reference.Name().String(), prefix) {
				stackBranches = append(stackBranches, reference.Name().String())
			}
			return nil
		})
		if err != nil {
			log.Printf("call=iter.ForEach err=`%v`\n", err)
		}
	}

	sort.Strings(stackBranches)

	fmt.Printf(`In stack %s
On branch %s

Stack:
`, stack, branch)
	for _, b := range stackBranches {
		var token = " "
		if ref.Name().String() == b {
			token = "*"
		}
		fmt.Printf("  %s %s\n", token, b)
	}

	return nil
}

var example = `In stack %s
On branch %s

Stack:
    001_migration   🚢 3629a61
  * 002_api         ✅ https://github.com/nfisher/gitit/pulls/110781
    003_ui          ❌ https://github.com/nfisher/gitit/pulls/110779

`

func main() {
	Execute()
}
