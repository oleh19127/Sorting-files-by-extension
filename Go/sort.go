package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/leaanthony/spinner"
	"github.com/mochi-co/autonamer"
)

const (
	sortedFilesFolder = "Sorted Files"
	otherFilesFolder  = "Other Files"
)

func folderExist(folder string) bool {
	if _, err := os.Stat(folder); !os.IsNotExist(err) {
		return true
	}
	return false
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
					if !folderExist(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder)) {
						os.MkdirAll(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder), 0755)
					}
					// If file already exists, increment filename: name.txt -> name(1).txt
					newPath, err := autonamer.Pick(1000, filepath.Join(sortedFilesFolder, modTimeFolder, data.folder, info.Name()))
					if err != nil {
						fmt.Println(err)
					}
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
	sortingByExt.Success()
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
		if !strings.HasPrefix(info.Name(), "sort") && !strings.HasPrefix(info.Name(), sortedFilesFolder) && !strings.HasPrefix(info.Name(), otherFilesFolder) && !info.IsDir() {
			// If folders not exist create
			if !folderExist(filepath.Join(otherFilesFolder, fileExtname)) {
				os.MkdirAll(filepath.Join(otherFilesFolder, fileExtname), 0755)
			}
			// If file already exists, increment filename: name.txt -> name(1).txt
			newPath, err := autonamer.Pick(1000, filepath.Join(otherFilesFolder, fileExtname, info.Name()))
			if err != nil {
				fmt.Println(err)
			}
			// Move file
			os.Rename(path, newPath)
			othersFilesExist = true
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	sortingOthersFiles.Success()
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

func scanFolders() {
	filesByExtensionExist, folders := sortingByExt()
	otherFilesExist := sortingOthersFiles()
	if filesByExtensionExist || otherFilesExist {
		if folders > 0 {
			deleteEmptyFoldersSpinner := spinner.New("Delete empty folders...")
			deleteEmptyFoldersSpinner.Start()
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
		}
	}
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

func archiveSortedFilesFolder() {
	if folderExist(sortedFilesFolder) || folderExist(otherFilesFolder) {
		var archiveInput string
		fmt.Println("Want archive files? (yes or y/any key to not)")
		fmt.Scanln(&archiveInput)
		if strings.ToLower(archiveInput) == "yes" || strings.ToLower(archiveInput) == "y" {
			if folderExist(sortedFilesFolder) {
				sortedFilesFolderArchiveSpinner := spinner.New("Archive sorted files...")
				sortedFilesFolderArchiveSpinner.Start()
				// If file already exists, increment filename: Sorted Files.zip -> Sorted Files(1).zip
				newPathForSortedFilesFolder, err := autonamer.Pick(1000, sortedFilesFolder+".zip")
				if err != nil {
					fmt.Println(err)
				}
				zipIt(sortedFilesFolder, newPathForSortedFilesFolder, false)
				os.RemoveAll(sortedFilesFolder)
				sortedFilesFolderArchiveSpinner.Success()
			}
			if folderExist(otherFilesFolder) {
				otherFilesFolderArchiveSpinner := spinner.New("Archive other files...")
				otherFilesFolderArchiveSpinner.Start()
				// If file already exists, increment filename: Other Files.zip -> Other Files(1).zip
				newPathForOtherFilesFolder, err := autonamer.Pick(1000, otherFilesFolder+".zip")
				if err != nil {
					fmt.Println(err)
				}
				zipIt(otherFilesFolder, newPathForOtherFilesFolder, false)
				os.RemoveAll(otherFilesFolder)
				otherFilesFolderArchiveSpinner.Success()
			}
		}
	}
}

func main() {
	start := time.Now()
	scanFolders()
	archiveSortedFilesFolder()
	duration := time.Since(start)
	fmt.Println("Work time:", duration.Round(time.Millisecond))
	if runtime.GOOS == "windows" {
		var closeInput string
		fmt.Println("Press enter to close!!!")
		fmt.Scanln(&closeInput)
	}
}
