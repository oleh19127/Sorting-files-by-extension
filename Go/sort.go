package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
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
	// EMPTY FOLDERS
	emptyFolders int
)

// GLOBAL FOLDER
const sortedFiles = "Sorted Files"

func sorting() bool {
	// IF FILE TO SORT EXIST = true, DEFAULT = false
	fileToSortExists := false
	// CHECK ALL PATH
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// SKIP SORTED FILES DIR
		SkipDir := fs.SkipDir
		if path == sortedFiles {
			return SkipDir
		}
		if err != nil {
			fmt.Println(err)
		}
		// CALCULATE EMPTY FOLDERS
		if info.IsDir() {
			emptyFolders = emptyFolders + 1
		}
		// GET MODIFICATION FILE TIME
		modTimeFolder := info.ModTime().Format("2006")
		// CHECK FILES
		for _, data := range allData {
			extensions := strings.Split(data.extension, " ")
			for _, extension := range extensions {
				fileExtname := strings.Trim(filepath.Ext(info.Name()), ".")
				// SORTING
				if info.Name() != "sort_windows.exe" && info.Name() != "sort.exe" && strings.ToLower(fileExtname) == extension {
					// IF FOLDERS NOT EXIST CREATE
					if _, err := os.Stat(sortedFiles); os.IsNotExist(err) {
						os.Mkdir(sortedFiles, 0755)
						color.Yellow("Create " + sortedFiles + " folder")
					}
					if _, err := os.Stat(filepath.Join(sortedFiles, modTimeFolder)); os.IsNotExist(err) {
						os.MkdirAll(filepath.Join(sortedFiles, modTimeFolder), 0755)
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
	color.Red(path + " folder removed!!!")
}

func main() {
	if sorting() {
		for i := 0; i < emptyFolders-1; i++ {
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
	} else {
		color.Cyan("NOTHING TO SORT!!!")
	}
	// UNCOMMENT IF BUILD FOR WINDOWS
	// var input string
	// color.Cyan("PRESS ENTER TO CLOSE!!!")
	// fmt.Scanln(&input)
}
