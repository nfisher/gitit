package cmd

import "io"

type Flags struct {
	SubCommand string
	BranchName string
}

const (
	Checkout = "checkout"
	Init     = "init"
	Push     = "push"
	Rebase   = "rebase"
	Status   = "status"
	Squash   = "squash"
)

const (
	Success = iota
	MissingArguments
	MissingSubCommand
)

func Exec(input Flags, _ io.Writer) int {
	switch input.SubCommand {
	case Checkout:
		if input.BranchName == "" {
			return MissingArguments
		}
		return Success
	case Init:
		if input.BranchName == "" {
			return MissingArguments
		}
		return Success
	case Push:
		return Success
	case Rebase:
		return Success
	case Squash:
		return Success
	case Status:
		return Success
	default:
		return MissingSubCommand
	}
}
