package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leaanthony/spinner"
)

func removeEmptyFolders() (bool, bool) {
	filesByExtensionExist, folders := sortingByExt()
	otherFilesExist := sortingOthersFiles()
	if filesByExtensionExist || otherFilesExist {
		deleteEmptyFoldersSpinner := spinner.New("Delete empty folders...")
		deleteEmptyFoldersSpinner.Start()
		if folders > 0 {
			for i := 0; i < folders; i++ {
				err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
					if info.IsDir() && info.Name() == sortedFilesFolder || info.Name() == otherFilesFolder {
						return filepath.SkipDir
					}
					if info.IsDir() {
						removeDir(path, info)
					}
					return nil
				})
				if err != nil {
					fmt.Println(err)
				}
			}
			deleteEmptyFoldersSpinner.Success()
		} else {
			deleteEmptyFoldersSpinner.Error("Empty folders not exist")
		}
	}
	return filesByExtensionExist, otherFilesExist
}
