package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/0x00-ketsu/taskcli/internal/cmd"
)


func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("Can't execute this on a windows machine")
		os.Exit(0)
	}

	cmd.Execute()
}
