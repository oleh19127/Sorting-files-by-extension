# List of extension

- Videos: mp4, mov, wmv, fly, avi, mkv, flv, mpg, webm, oog, m4p, m4v, qt, swf, avchd, f4v, mpeg-2

- Images: png, jpg, webp, svg, gif, ico, jpeg, bmp, esp, jpeg 2000, heif, bat, cgm, tif, tiff, eps, raw, cr2, nef, orf, sr2

- Music: mp3, aac, flac, alac, wav, aiff, dsd, pcm, m4a, wma

- Documents: txt, doc, docx, docx, odt, xls, xlsx, ppt, pptx

- Psd: psd

- Pdf: pdf

- Archive: zip, rar, 7z, tar

- Torrent: torrent

- Exe: exe

## The program will sort the files in all subfolders of the folder where you put it

## After start program

- All files located in **Sorted Files**

- All images files located in **Sorted Files/year created file/Images/**

- All documents files located in **Sorted Files/year created file/Documents/**

- All videos files located in **Sorted Files/year created file/Videos/**

- All music files located in **Sorted Files/year created file/Music/**

- All pdf files located in **Sorted Files/year created file/Pdf/**

- All psd files located in **Sorted Files/year created file/Psd/**

- All archive files located in **Sorted Files/year created file/Archive/**

- All torrent files located in **Sorted Files/year created file/Torrent/**

- All exe files located in **Sorted Files/year created file/Exe/**

- All other files located in **Other Files/extname/**

- All other files no extension located in **Other Files/**

- All empty folders will be deleted

- Possible to archive the Sorted Files and Other Files folders

## Demo

![Windows Demo](https://github.com/oleh312/Sorting-files-by-extension/blob/main/assets/windows_demo.gif)

# How to use

## Windows

1. Keep and download file(No virus) <a  href="https://github.com/oleh312/Sorting-files-by-extension/releases/download/v1.8/sort_windows.exe">sort_windows.exe</a>

2. Place the program file in the folder with the files you want to sort

3. Double click on file if not work <a  href="https://www.google.com/search?q=how+to+change+permissions+of+a+file+in+windows&sxsrf=ALeKk03ByQLIy_kPt0X2erLRnJHUqJrZDw%3A1628627435772&ei=6-ESYcvOLs3LrgTmwZ_oBA&oq=how+to+change+permissions+of+a+file+in+windows&gs_lcp=Cgdnd3Mtd2l6EAMyBQgAEMsBMgUIABDLATIGCAAQFhAeMgYIABAWEB4yBggAEBYQHjIGCAAQFhAeMgYIABAWEB4yBggAEBYQHjIGCAAQFhAeMgYIABAWEB46BwgAEEcQsAM6BwgAELADEENKBAhBGABQg9AbWLHZG2D73xtoAnACeACAAbIBiAHMB5IBAzAuOJgBAKABAcgBCcABAQ&sclient=gws-wiz&ved=0ahUKEwiL8J7-pafyAhXNpYsKHebgB00Q4dUDCA4&uact=5">Change file permission </a> or <a  href="https://support.microsoft.com/en-us/windows/turn-off-defender-antivirus-protection-in-windows-security-99e6004f-c54c-8509-773c-a4d776b77960">Off real time protection on windows</a>

4. If you want to sort the files and automatically archive them, open terminal in the folder and write to the terminal: sort_windows.exe archive --all

## Linux

1. Download file <a  href="https://github.com/oleh312/Sorting-files-by-extension/releases/download/v1.8/sort_linux">sort_linux</a>

2. Place the program file in the folder with the files you want to sort

3. Open terminal in the folder

4. Change file permission: sudo chmod +x "full path to file(sort_linux)"

5. Write in the terminal: ./sort_linux

6. If you want to sort the files and automatically archive them, write to the terminal: ./sort_linux archive --all

## If not work in you system

1. Download Golang in <a  href="https://golang.org/">offical site</a>

2. Download or copy files <a  href="https://github.com/oleh312/Sorting-files-by-extension/blob/main/Go/">files</a>

3. Place the program file in the folder with the files you want to sort

4. Open terminal or cmd in the folder

5. Write in the terminal or cmd: go mod init sorty && go mod tidy && go build sort.go

6. You can remove files: sort.go, linux.go, structureOfExtensions.go, windows.go, go mod, go.sum

5. You can move the compiled file where you want and run it on your system
