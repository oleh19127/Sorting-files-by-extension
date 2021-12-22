package main

import (
	"flag"
	"fmt"
	"os"
)

func handleGet(getCmd *flag.FlagSet, all *bool) {
	getCmd.Parse(os.Args[2:])
	if !*all {
		fmt.Print("id is require or specify --all for all folders")
		getCmd.PrintDefaults()
		os.Exit(1)
	}
	if *all {
		ch := make(chan string)
		go archiveSortedFilesFolder(ch)
		go archiveOtherFilesFolder(ch)
		s1 := <-ch
		fmt.Println(s1)
		s := <-ch
		fmt.Println(s)
		return
	}
}
