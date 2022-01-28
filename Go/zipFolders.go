package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/leaanthony/spinner"
)

func zipIt(source, target string, needBaseDir bool) error {
	// Uncomment "setCpuPriorityWindows()" if build for windows(Windows only), set priority: below normal, normal, above normal, hight
	// setCpuPriorityWindows("below normal")
	// Uncomment "setCpuPriorityLinux()" if build for linux(Linux only), set priority from 20 to 9: very low, from 10 to 1: low, from 0 to -9: normal, from -10 to -20: hight
	// setCpuPriorityLinux(5)
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if baseDir != "" {
			if needBaseDir {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			} else {
				path := strings.TrimPrefix(path, source)
				if len(path) > 0 && (path[0] == '/' || path[0] == '\\') {
					path = path[1:]
				}
				if len(path) == 0 {
					return nil
				}
				header.Name = path
			}
		}
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}

func archiveSortedFilesFolder(exit chan string) {
	if pathExist(sortedFilesFolder) {
		sortedFilesFolderArchiveSpinner := spinner.New("Archive sorted files...")
		sortedFilesFolderArchiveSpinner.Start()
		// If file already exists, increment filename: Sorted Files.zip -> Sorted Files(1).zip
		newPathForSortedFilesFolder := incrementFileName(sortedFilesFolder + ".zip")
		zipIt(sortedFilesFolder, newPathForSortedFilesFolder, false)
		os.RemoveAll(sortedFilesFolder)
		sortedFilesFolderArchiveSpinner.Success()
	}
	exit <- ""
}

func archiveOtherFilesFolder(exit chan string) {
	if pathExist(otherFilesFolder) {
		otherFilesFolderArchiveSpinner := spinner.New("Archive other files...")
		otherFilesFolderArchiveSpinner.Start()
		// If file already exists, increment filename: Other Files.zip -> Other Files(1).zip
		newPathForOtherFilesFolder := incrementFileName(otherFilesFolder + ".zip")
		zipIt(otherFilesFolder, newPathForOtherFilesFolder, false)
		os.RemoveAll(otherFilesFolder)
		otherFilesFolderArchiveSpinner.Success()
	}
	exit <- ""
}

func archiveFolders() {
	sortByExt, sortOtherFiles := removeEmptyFolders()
	if sortByExt || sortOtherFiles || pathExist(otherFilesFolder) || pathExist(sortedFilesFolder) {
		userInput := getUserInput()
		if strings.ToLower(userInput) == "yes" || strings.ToLower(userInput) == "y" {
			ch := make(chan string)
			if sortByExt && sortOtherFiles || pathExist(otherFilesFolder) || pathExist(sortedFilesFolder) {
				go archiveSortedFilesFolder(ch)
				go archiveOtherFilesFolder(ch)
				s := <-ch
				fmt.Println(s)
				s1 := <-ch
				fmt.Println(s1)
			}
			if sortByExt || pathExist(sortedFilesFolder) {
				go archiveSortedFilesFolder(ch)
				s1 := <-ch
				fmt.Println(s1)
			}
			if sortOtherFiles || pathExist(otherFilesFolder) {
				go archiveOtherFilesFolder(ch)
				s := <-ch
				fmt.Println(s)
			}
		}
	}
}
