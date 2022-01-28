package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func incrementFileName(path string) string {
	folder := filepath.Dir(path)
	file := filepath.Base(path)
	re := regexp.MustCompile(`^(.+?)(\((\d.*)\)|)(\..*|$)(.*?)$`)
	matches := re.FindAllStringSubmatch(file, -1)
	label := strings.TrimSpace(matches[0][1])
	ext := strings.TrimSpace(matches[0][4])
	number := 0
	if pathExist(path) {
		for pathExist(path) {
			number++
			newPath := filepath.Join(folder, fmt.Sprintf("%s (%d)%s", label, number, ext))
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				return newPath
			}
		}
	}
	return path
}
