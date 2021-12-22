package main

type Data struct {
	folder     string
	extensions []string
}

var (
	// List of extension
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
