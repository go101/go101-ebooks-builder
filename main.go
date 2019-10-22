package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"strings"
)

// The code of this project is very ugly. It just makes the job done.

var bookProjectDirFlag = flag.String("book-project-dir", "", "the path to the book project")
var bookVersionFlag = flag.String("book-version", "", "the version of the book")
var targetFlag = flag.String("target", "all", "output target (epub | azw3 | mobi | apple | pdf | print | all)")

func main() {
	log.SetFlags(0)
	flag.Parse()

	bookProjectDirs := make([]string, 0, 2)
	projectDir := strings.TrimSpace(*bookProjectDirFlag)
	if projectDir != "" {
		bookProjectDirs = append(bookProjectDirs, projectDir)
	} else {
		if path := os.Getenv("Go101Path"); path != "" {
			bookProjectDirs = append(bookProjectDirs, path)
		}
		if path := os.Getenv("Golang101Path"); path != "" {
			bookProjectDirs = append(bookProjectDirs, path)
		}
	}
	
	if len(bookProjectDirs) == 0 {
		log.Fatal("-book-project-dir is required.")
	}

	for _, bookProjectDir := range bookProjectDirs {
		bookVersion := strings.TrimSpace(*bookVersionFlag)
		if bookVersion == "" {
			tag := runShellCommand2(true, bookProjectDir, "git", "describe", "--tags", "--abbrev=0")
			tag = bytes.TrimSpace(tag)
			if len(tag) > 0 {
				bookVersion = string(tag)
			}
		}
		
		if bookVersion == "" {	
			log.Fatal("-book-version is required.")
		}

		coverImagePath := buildCoverImage(bookProjectDir, bookVersion)

		switch target := strings.TrimSpace(*targetFlag); target {
		case "epub":
			genetateEpubFile(bookProjectDir, bookVersion, coverImagePath)
		case "azw3":
			genetateAzw3File(bookProjectDir, bookVersion, coverImagePath)
		case "mobi":
			genetateMobiFile(bookProjectDir, bookVersion, coverImagePath)
		case "apple":
			genetateAppleFile(bookProjectDir, bookVersion, coverImagePath)
		case "pdf":
			genetatePdfFile(bookProjectDir, bookVersion, coverImagePath, false)
		case "print":
			genetatePdfFile(bookProjectDir, bookVersion, coverImagePath, true)
		case "all", "":
			genetateAzw3File(bookProjectDir, bookVersion, coverImagePath)
			genetateEpubFile(bookProjectDir, bookVersion, coverImagePath)
			genetateMobiFile(bookProjectDir, bookVersion, coverImagePath)
			genetateAppleFile(bookProjectDir, bookVersion, coverImagePath)
			genetatePdfFile(bookProjectDir, bookVersion, coverImagePath, false)
			genetatePdfFile(bookProjectDir, bookVersion, coverImagePath, true)
		default:
			log.Fatal("Unknown target:", target)
		}
	}
}
