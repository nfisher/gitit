package main

import (
	"flag"
	"github.com/nfisher/gitit/cmd"
	"os"
)

var example = `In stack %s
On branch %s

Stack:
    001_migration   🚢 3629a61
  * 002_api         ✅ https://github.com/nfisher/gitit/pulls/110781
    003_ui          ❌ https://github.com/nfisher/gitit/pulls/110779

`

func main() {
	var input cmd.Flags
	flag.Parse()

	input.SubCommand = flag.Arg(0)

	os.Exit(cmd.Exec(input, os.Stdout))
}
