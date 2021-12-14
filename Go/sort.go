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

type Data struct {
	folder     string
	extensions []string
}

var (
	// Structure
	images = Data{
		folder:     "Images",
		extensions: []string{"png", "jpg", "webp", "svg", "gif", "ico", "jpeg", "bmp", "esp", "jpeg 2000", "heif", "bat", "cgm", "tif", "tiff", "eps", "raw", "cr2", "nef", "orf", "sr2"},
	}
	videos = Data{
		folder:     "Videos",
		extensions: []string{"mp4", "mov", "wmv", "fly", "avi", "mkv", "flv", "mpg", "webm", "oog", "m4p", "m4v", "qt", "swf", "avchd", "f4v", "mpeg-2"},
	}
	music = Data{
		folder:     "Music",
		extensions: []string{"mp3", "aac", "flac", "alac", "wav", "aiff", "dsd", "pcm", "m4a", "wma"},
	}
	documents = Data{
		folder:     "Documents",
		extensions: []string{"txt", "doc", "docx", "docx", "odt", "xls", "xlsx", "ppt", "pptx"},
	}
	psd = Data{
		folder:     "Psd",
		extensions: []string{"psd"},
	}
	pdf = Data{
		folder:     "Pdf",
		extensions: []string{"pdf"},
	}
	archive = Data{
		folder:     "Archive",
		extensions: []string{"zip", "rar", "7z", "tar"},
	}
	exe = Data{
		folder:     "Exe",
		extensions: []string{"exe"},
	}
	torrent = Data{
		folder:     "Torrent",
		extensions: []string{"torrent"},
	}
	allData = []Data{images, videos, music, documents, psd, pdf, archive, exe, torrent}
)

// Sorted files folder
const sortedFilesFolder = "Sorted Files"

func sorting() (bool, int) {
	sortingSpinner := spinner.New("Sorting...")
	sortingSpinner.Start()
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
			for _, extension := range data.extensions {
				fileExtname := strings.Trim(filepath.Ext(info.Name()), ".")
				if !strings.HasPrefix(info.Name(), "sort") && !strings.HasPrefix(info.Name(), sortedFilesFolder) && strings.EqualFold(fileExtname, extension) {
					// Get modification file time
					modTimeFolder := strconv.Itoa(info.ModTime().Year())
					// If folders not exist create
					if _, err := os.Stat(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder)); os.IsNotExist(err) {
						os.MkdirAll(filepath.Join(sortedFilesFolder, modTimeFolder, data.folder), 0755)
					}
					newPath, err := autonamer.Pick(1000, filepath.Join(sortedFilesFolder, modTimeFolder, data.folder, info.Name()))
					if err != nil {
						fmt.Println(err)
					}
					// Move file
					os.Rename(path, newPath)
					// Update spinner
					fileToSortExists = true
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	sortingSpinner.Success()
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
			deleteEmptyFoldersSpinner := spinner.New("Delete empty folders...")
			deleteEmptyFoldersSpinner.Start()
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
			deleteEmptyFoldersSpinner.Success()
		}
	}
}

func zipIt(source, target string, needBaseDir bool) error {

	// Uncomment "controllCpuPriorityWindows()" if build for windows(Windows only)
	// controllCpuPriorityWindows()

	// Uncomment "controllCpuPriorityLinux()" if build for linux(Linux only)
	// controllCpuPriorityLinux()

	archiveSpinner := spinner.New("Archive files...")
	archiveSpinner.Start()
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
	archiveSpinner.Success()
	return err
}

func archiveSortedFilesFolder() {
	if _, err := os.Stat(sortedFilesFolder); !os.IsNotExist(err) {
		var archiveInput string
		fmt.Println("Want archive files? (yes/no)")
		fmt.Scanln(&archiveInput)
		if strings.ToLower(archiveInput) == "yes" || strings.ToLower(archiveInput) == "y" {
			newPath, err := autonamer.Pick(1000, sortedFilesFolder+".zip")
			if err != nil {
				fmt.Println(err)
			}
			zipIt(sortedFilesFolder, newPath, false)
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
	if runtime.GOOS == "windows" {
		var closeInput string
		fmt.Println("Press enter to close!!!")
		fmt.Scanln(&closeInput)
	}
}
