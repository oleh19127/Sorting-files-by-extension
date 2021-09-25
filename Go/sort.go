package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
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
	// GET ALL FOLDERS
	folders int
	// GET ALL SORTED FILES
	files int64
)

func sorting() bool {
	// IF FILE TO SORT EXIST = true, DEFAULT = false
	fileToSortExists := false
	// CHECK ALL PATH
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// GLOBAL FOLDER
		const sortedFiles = "Sorted Files"
		// SKIP SORTED FILES DIR
		SkipDir := fs.SkipDir
		if path == sortedFiles {
			return SkipDir
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
				if info.Name() != "sort_windows.exe" && info.Name() != "sort.exe" && strings.ToLower(fileExtname) == extension {
					// CALCULATE FILES
					files = files + 1
					// GET MODIFICATION FILE TIME
					modTimeFolder := info.ModTime().Format("2006")
					// IF FOLDERS NOT EXIST CREATE
					if _, err := os.Stat(sortedFiles); os.IsNotExist(err) {
						os.Mkdir(sortedFiles, 0755)
					}
					if _, err := os.Stat(filepath.Join(sortedFiles, modTimeFolder)); os.IsNotExist(err) {
						os.Mkdir(filepath.Join(sortedFiles, modTimeFolder), 0755)
					}
					if _, err := os.Stat(filepath.Join(sortedFiles, modTimeFolder, data.folder)); os.IsNotExist(err) {
						os.Mkdir(filepath.Join(sortedFiles, modTimeFolder, data.folder), 0755)
					}
					// MOVE FILE
					os.Rename(path, filepath.Join(sortedFiles, modTimeFolder, data.folder, info.Name()))
					fileToSortExists = true

				}
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	// PROGRESSBAR
	if files > 0 {
		fmt.Println("SORTING!!!")
		bar := progressbar.Default(files)
		for i := files; i > 0; i-- {
			bar.Add(1)
			time.Sleep(time.Millisecond)
		}
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

func main() {
	if sorting() {
		if folders-1 > 0 {
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
			fmt.Println("EMPTY FOLDERS DELETED!!!")
		}
	} else {
		fmt.Println("NOTHING TO SORT!!!")
	}
	// UNCOMMENT IF BUILD FOR WINDOWS
	var input string
	fmt.Println("PRESS ENTER TO CLOSE!!!")
	fmt.Scanln(&input)
}
