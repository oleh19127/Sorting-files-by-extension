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
	"time"

	"github.com/mochi-co/autonamer"
	"github.com/theckman/yacspin"
)

type Data struct {
	folder, extension string
}

var (
	// Structure
	allData = []Data{
		{"Images", "png, jpg, webp, svg, gif, ico, jpeg, bmp, esp, jpeg 2000, heif, bat, cgm, tif, tiff, eps, raw, cr2, nef, orf, sr2"},
		{"Videos", "mp4, mov, wmv, fly, avi, mkv, flv, mpg, webm, oog, m4p, m4v, qt, swf, avchd, f4v, mpeg-2"},
		{"Music", "mp3, aac, flac, alac, wav, aiff, dsd, pcm, m4a, wma"},
		{"Documents", "txt, doc, docx, docx, odt, xls, xlsx, ppt, pptx"},
		{"Psd", "psd"},
		{"pdf", "pdf"},
		{"Archive", "zip, rar, 7z, tar"},
		{"Exe", "exe"},
		{"Torrent", "torrent"},
	}
)

// Sorted files folder
const sortedFilesFolder = "Sorted Files"

func sorting() (bool, int) {

	sortingSpinnerCfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          " Sorting",
		SuffixAutoColon: true,
		Message:         "files",
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}
	spinner, _ := yacspin.New(sortingSpinnerCfg)
	spinner.Start()

	var fileToSortExists bool
	var calcFolders int
	// Check path
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == sortedFilesFolder {
			return filepath.SkipDir
		}
		// Calculate folders
		if info.IsDir() && path != "." {
			calcFolders = calcFolders + 1
		}
		// Check files
		for _, data := range allData {
			extensions := strings.Split(data.extension, ", ")
			for _, extension := range extensions {
				fileExtname := strings.Trim(filepath.Ext(info.Name()), ".")
				if info.Name() != "sort_windows.exe" && info.Name() != sortedFilesFolder+".zip" && info.Name() != "sort.exe" && strings.ToLower(fileExtname) == strings.ToLower(extension) {
					// Get modification file time
					modTimeFolder := strconv.Itoa(info.ModTime().Year())
					// If folders not exist create
					if _, err := os.Stat(sortedFilesFolder); os.IsNotExist(err) {
						os.Mkdir(sortedFilesFolder, 0755)
					}
					if _, err := os.Stat(filepath.Join(sortedFilesFolder, modTimeFolder)); os.IsNotExist(err) {
						os.Mkdir(filepath.Join(sortedFilesFolder, modTimeFolder), 0755)
					}
					if _, err := os.Stat(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder)); os.IsNotExist(err) {
						os.Mkdir(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder), 0755)
					}
					newPath, err := autonamer.Pick(1000, filepath.Join(sortedFilesFolder, modTimeFolder, data.folder, info.Name()))
					if err != nil {
						fmt.Println(err)
					}
					// Move file
					os.Rename(path, newPath)
					// Update spinner
					spinner.Message(info.Name())
					fileToSortExists = true
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	// Stop spinner
	spinner.Stop()
	return fileToSortExists, calcFolders
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
	filesExist, folders := sorting()
	if filesExist {
		if folders > 0 {

			deleteEmptyFoldersSpinnerCfg := yacspin.Config{
				Frequency:       100 * time.Millisecond,
				CharSet:         yacspin.CharSets[59],
				Suffix:          " Delete empty folders",
				SuffixAutoColon: true,
				Message:         "files",
				StopCharacter:   "✓",
				StopColors:      []string{"fgGreen"},
			}
			spinner, _ := yacspin.New(deleteEmptyFoldersSpinnerCfg)
			spinner.Start()

			for i := 0; i < folders; i++ {
				err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						removeDir(path, info)
					}
					return nil
				})
				if err != nil {
					fmt.Println(err)
				}
			}

			time.Sleep(500 * time.Millisecond)
			spinner.Stop()

		}
	}
}

func zipIt(source, target string, needBaseDir bool) error {

	archiveSpinnerCfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[11],
		Suffix:          " Add files to archive",
		SuffixAutoColon: true,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}
	spinner, _ := yacspin.New(archiveSpinnerCfg)
	spinner.Start()

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

	spinner.Stop()

	return err
}

func archiveSortedFilesFolder() {
	if _, err := os.Stat(sortedFilesFolder); !os.IsNotExist(err) {
		var archiveInput string
		fmt.Println("Want archive files? (yes/no)")
		fmt.Scanln(&archiveInput)
		if strings.ToLower(archiveInput) == "yes" || strings.ToLower(archiveInput) == "y" {
			zipIt(sortedFilesFolder, sortedFilesFolder+".zip", false)
			os.RemoveAll(sortedFilesFolder)
		}
	}
}

func main() {
	start := time.Now()

	scanFolders()
	archiveSortedFilesFolder()

	duration := time.Since(start)
	fmt.Println("Work time:", duration.Round(time.Millisecond))
	
	// Uncomment if build for windows
	// var closeInput string
	// fmt.Println("Press enter to close!!!")
	// fmt.Scanln(&closeInput)
}
