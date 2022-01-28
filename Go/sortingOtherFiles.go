package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leaanthony/spinner"
)

func sortingOthersFiles() bool {
	sortingOthersFiles := spinner.New("Sorting other files...")
	sortingOthersFiles.Start()
	var othersFilesExist bool
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == sortedFilesFolder || info.Name() == otherFilesFolder {
			return filepath.SkipDir
		}
		fileExtname := strings.Trim(filepath.Ext(info.Name()), ".")
		if !strings.HasPrefix(info.Name(), "sort") && !strings.HasSuffix(info.Name(), "go") && !strings.HasSuffix(info.Name(), "mod") && !strings.HasSuffix(info.Name(), "sum") && !strings.HasPrefix(info.Name(), sortedFilesFolder) && !strings.HasPrefix(info.Name(), otherFilesFolder) && !info.IsDir() {
			// If folders not exist create
			if !pathExist(filepath.Join(otherFilesFolder, fileExtname)) {
				os.MkdirAll(filepath.Join(otherFilesFolder, fileExtname), 0755)
			}
			// If file already exists, increment filename: name.txt -> name(1).txt
			newPath := incrementFileName(filepath.Join(otherFilesFolder, fileExtname, info.Name()))
			// Move file
			os.Rename(path, newPath)
			othersFilesExist = true
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	if othersFilesExist {
		sortingOthersFiles.Success()
	} else {
		sortingOthersFiles.Error("Other files not exist")
	}
	return othersFilesExist
}
