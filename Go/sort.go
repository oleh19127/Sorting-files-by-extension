package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/leaanthony/spinner"
)

type Data struct {
	folder    string
	extension string
}

var (
	// STRUCTURE OF TYPES
	images   = Data{"Images", "png jpg webp svg gif ico jpeg bmp esp jpeg 2000 heif bat cgm tif tiff eps raw cr2 nef orf sr2"}
	videos   = Data{"Videos", "mp4 mov wmv fly avi mkv flv mpg webm oog m4p m4v qt swf avchd f4v mpeg-2"}
	music    = Data{"Music", "mp3 aac flac alac wav aiff dsd pcm m4a wma"}
	document = Data{"Documents", "txt doc docx docx odt xls xlsx ppt pptx"}
	psd      = Data{"Psd", "psd"}
	pdf      = Data{"pdf", "pdf"}
	archive  = Data{"Archive", "zip rar 7z tar"}
	torrent  = Data{"Torrent", "torrent"}
	exe      = Data{"Exe", "exe"}
	allData  = []Data{images, videos, music, document, psd, pdf, archive, torrent, exe}
	// ALL FOLDERS
	folders int
)

// SORTED FILES FOLDER
const sortedFiles = "Sorted Files"

func sorting() bool {
	// IF FILE TO SORT EXIST = true, DEFAULT = false
	fileToSortExists := false
	// CHECK ALL PATH
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// SKIP DIR
		if path == sortedFiles {
			return fs.SkipDir
		}
		if err != nil {
			fmt.Println(err)
		}
		// CALCULATE FOLDERS
		if info.IsDir() {
			folders = folders + 1
		}
		// CHECK FILES
		for _, data := range allData {
			extensions := strings.Split(data.extension, " ")
			for _, extension := range extensions {
				fileExtname := strings.Trim(filepath.Ext(info.Name()), ".")
				// SORTING
				if info.Name() != "sort_windows.exe" && info.Name() != sortedFiles+".zip" && info.Name() != "sort.exe" && strings.ToLower(fileExtname) == extension {
					// GET MODIFICATION FILE TIME
					modTimeFolder := strconv.Itoa(info.ModTime().Year())
					// IF FOLDERS NOT EXIST CREATE
					if _, err := os.Stat(sortedFiles); os.IsNotExist(err) {
						os.Mkdir(sortedFiles, 0755)
						color.Yellow("Create " + sortedFiles + " folder")
					}
					if _, err := os.Stat(filepath.Join(sortedFiles, modTimeFolder)); os.IsNotExist(err) {
						os.Mkdir(filepath.Join(sortedFiles, modTimeFolder), 0755)
						color.Yellow("Create " + sortedFiles + "/" + modTimeFolder + " folder")
					}
					if _, err := os.Stat(filepath.Join(sortedFiles, modTimeFolder, data.folder)); os.IsNotExist(err) {
						os.Mkdir(filepath.Join(sortedFiles, modTimeFolder, data.folder), 0755)
						color.Yellow("Create " + sortedFiles + "/" + modTimeFolder + "/" + data.folder + " folder")
					}
					// MOVE FILE
					os.Rename(path, filepath.Join(sortedFiles, modTimeFolder, data.folder, info.Name()))
					color.Green(path + " moved >> " + filepath.Join(sortedFiles, modTimeFolder, data.folder, info.Name()))
					fileToSortExists = true
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return fileToSortExists
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

func zipIt(source, target string, needBaseDir bool) error {
	zipSpinner := spinner.New("ADD FILES TO ARCHIVE...")
	zipSpinner.Start()
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
	zipSpinner.Success("DONE!!!")
	return err
}

func ifSortedFilesFolderExist() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() == sortedFiles && file.IsDir() {
			var archiveInput string
			color.HiBlue("WANT ARCHIVE FILES? (yes/no)")
			fmt.Scanln(&archiveInput)
			if strings.ToLower(archiveInput) == "yes" || strings.ToLower(archiveInput) == "y" {
				zipIt(sortedFiles, sortedFiles+".zip", false)
			}
		}
	}
}

func scanFolders() {
	if sorting() {
		if folders-1 > 0 {
			removeFoldersSpinner := spinner.New("SCAN FOLDERS!!!")
			removeFoldersSpinner.Start()
			for i := 0; i < folders-1; i++ {
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
			removeFoldersSpinner.Success("EMPTY FOLDERS DELETED!!!")
		}
	} else {
		color.Cyan("NOTHING TO SORT!!!")
	}
	ifSortedFilesFolderExist()
}

func main() {
	scanFolders()

	// UNCOMMENT IF BUILD FOR WINDOWS
	// var closeInput string
	// color.White("PRESS ENTER TO CLOSE!!!")
	// fmt.Scanln(&closeInput)
}

// func unzip(archive, target string) error {
// 	reader, err := zip.OpenReader(archive)
// 	if err != nil {
// 		return err
// 	}
// 	defer reader.Close()

// 	if err := os.MkdirAll(target, 0755); err != nil {
// 		return err
// 	}

// 	for _, file := range reader.File {
// 		path := filepath.Join(target, file.Name)
// 		if file.FileInfo().IsDir() {
// 			os.MkdirAll(path, file.Mode())
// 			continue
// 		}

// 		fileReader, err := file.Open()
// 		if err != nil {
// 			return err
// 		}
// 		defer fileReader.Close()

// 		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
// 		if err != nil {
// 			return err
// 		}
// 		defer targetFile.Close()

// 		if _, err := io.Copy(targetFile, fileReader); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
