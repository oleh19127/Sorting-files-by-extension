package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/leaanthony/spinner"
)

func sortingByExt() (bool, int) {
	sortingByExt := spinner.New("Sorting files by extensions...")
	sortingByExt.Start()
	var fileToSortExists bool
	var calcFolders int
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == sortedFilesFolder || info.Name() == otherFilesFolder {
			return filepath.SkipDir
		}
		// Calculate folders
		if info.IsDir() && path != "." {
			calcFolders = calcFolders + 1
		}
		// Check files
		for _, data := range allData {
			for _, extension := range data.extensions {
				fileExtname := strings.Trim(filepath.Ext(info.Name()), ".")
				if !strings.HasPrefix(info.Name(), "sort") && !strings.HasPrefix(info.Name(), sortedFilesFolder) && strings.EqualFold(fileExtname, extension) && !strings.HasPrefix(info.Name(), otherFilesFolder) {
					// Get modification file time
					modTimeFolder := strconv.Itoa(info.ModTime().Year())
					// If folders not exist create
					if !pathExist(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder)) {
						os.MkdirAll(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder), 0755)
					}
					// If file already exists, increment filename: name.txt -> name(1).txt
					newPath := incrementFileName(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder, info.Name()))
					// Move file
					os.Rename(path, newPath)
					fileToSortExists = true
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	if fileToSortExists {
		sortingByExt.Success()
	} else {
		sortingByExt.Error("Files to sort by extension not exist")
	}
	return fileToSortExists, calcFolders
}
