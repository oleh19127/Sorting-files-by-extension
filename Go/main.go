package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	start := time.Now()
	archiveCmd := flag.NewFlagSet("archive", flag.ExitOnError)
	archiveAll := archiveCmd.Bool("all", false, "Archive all folders")
	if len(os.Args) == 1 {
		scanFolders()
		archiveFolders()
		if runtime.GOOS == "windows" {
			var closeInput string
			fmt.Println("Press enter to close!!!")
			fmt.Scanln(&closeInput)
		}
	} else if len(os.Args) == 3 {
		switch os.Args[1] {
		case "archive":
			scanFolders()
			handleGet(archiveCmd, archiveAll)
		default:
			fmt.Println("Don't understand input")
		}
	} else if len(os.Args) == 2 {
		fmt.Println("Expected subcommands, example: 'archive --all'")
		os.Exit(1)
	} else {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}
	duration := time.Since(start)
	fmt.Println("Work time:", duration.Round(time.Millisecond))
}
