package main

import (
	"fmt"

	"github.com/ikedam/pictmanager/cmd/random/cmd"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	cmd.SetVersion(fmt.Sprintf("%v:%v", version, commit))
	cmd.Execute()
}
