package main

import (
	"github.com/nfisher/gitit/cmd"
	"os"
)

func main() {
	var input cmd.Flags

	if len(os.Args) > 1 {
		input.SubCommand = os.Args[1]
	}

	if len(os.Args) > 2 {
		input.BranchName = os.Args[2]
	}

	os.Exit(cmd.Exec(input, os.Stdout))
}

var example = `In stack %s
On branch %s

Stack:
    001_migration   🚢 3629a61
  * 002_api         ✅ https://github.com/nfisher/gitit/pulls/110781
    003_ui          ❌ https://github.com/nfisher/gitit/pulls/110779

`
