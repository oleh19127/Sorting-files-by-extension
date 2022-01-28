package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/leaanthony/spinner"
)

const (
	sortedFilesFolder = "Sorted Files"
	otherFilesFolder  = "Other Files"
)

func pathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

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

func removeDir(path string, info os.FileInfo) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	if len(files) != 0 {
		return
	}
	err = os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func scanFolders() (bool, bool) {
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

func getUserInput() string {
	var archiveInput string
	fmt.Println("Want archive files? (yes or y/any key to not)")
	fmt.Scanln(&archiveInput)
	return archiveInput
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
	sortByExt, sortOtherFiles := scanFolders()
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
