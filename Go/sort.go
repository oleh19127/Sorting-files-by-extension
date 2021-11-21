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
						fmt.Println("Create", sortedFilesFolder, "folder")
					}
					if _, err := os.Stat(filepath.Join(sortedFilesFolder, modTimeFolder)); os.IsNotExist(err) {
						os.Mkdir(filepath.Join(sortedFilesFolder, modTimeFolder), 0755)
						fmt.Println("Create", filepath.Join(sortedFilesFolder, modTimeFolder), "folder")
					}
					if _, err := os.Stat(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder)); os.IsNotExist(err) {
						os.Mkdir(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder), 0755)
						fmt.Println("Create", filepath.Join(sortedFilesFolder, modTimeFolder, data.folder), "folder")
					}
					// Move file
					os.Rename(path, filepath.Join(sortedFilesFolder, modTimeFolder, data.folder, info.Name()))
					fmt.Println(path, ">>", filepath.Join(sortedFilesFolder, modTimeFolder, data.folder, info.Name()))
					fileToSortExists = true
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
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
			fmt.Println("Checking folders...")
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
			fmt.Println("Empty folders deleted!")
		}
	} else {
		fmt.Println("Nothing to sort!")
	}
}

func zipIt(source, target string, needBaseDir bool) error {
	fmt.Println("Add files to archive...")
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

	fmt.Println("Done!")
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
	scanFolders()
	archiveSortedFilesFolder()
	// Uncomment if build for windows
	// var closeInput string
	// fmt.Println("Press enter to close!!!")
	// fmt.Scanln(&closeInput)
}
