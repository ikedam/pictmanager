package main

import (
	"fmt"

	"github.com/ikedam/pictmanager/cmd/server/cmd"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	cmd.SetVersion(fmt.Sprintf("%v:%v", version, commit))
	cmd.Execute()
}
