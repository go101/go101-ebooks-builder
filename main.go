package main

import (
	"flag"
	"log"
	"strings"
)

// The code of this project is very ugly. It just makes the job done.

var bookProjectDirFlag = flag.String("book-project-dir", "", "the path to the book project")
var bookVersionFlag = flag.String("book-version", "", "the version of the book")
var targetFlag = flag.String("target", "all", "output target (epub | azw3 | mobi | apple | pdf | print | all)")

// Go 101:
// export BookVersion=v1.12.a
// go run *.go -target=epub -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/go101/go101
// go run *.go -target=azw3 -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/go101/go101
// go run *.go -target=mobi -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/go101/go101
// go run *.go -target=pdf -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/go101/go101

// Golang 101:
// export BookVersion=v1.12.a
// go run *.go -target=epub -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/golang101/golang101
// go run *.go -target=azw3 -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/golang101/golang101
// go run *.go -target=mobi -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/golang101/golang101
// go run *.go -target=pdf -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/golang101/golang101

/*
   go run *.go -target=all -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/go101/go101
   go run *.go -target=all -book-version=$BookVersion -book-project-dir=/home/lx/rsync/projects/go101.org/golang101/golang101
*/


func main() {
	log.SetFlags(0)
	flag.Parse()
	
	bookProjectDir := strings.TrimSpace(*bookProjectDirFlag)
	if bookProjectDir == "" {
		log.Fatal("-book-project-dir is required.")
	}
	
	bookVersion := strings.TrimSpace(*bookVersionFlag)
	if bookVersion == "" {
		log.Fatal("-book-version is required.")
	}
	
	switch target := strings.TrimSpace(*targetFlag); target {
	case "epub":
		genetateEpubFile(bookProjectDir, bookVersion)
	case "azw3":
		genetateAzw3File(bookProjectDir, bookVersion)
	case "mobi":
		genetateMobiFile(bookProjectDir, bookVersion)
	case "apple":
		genetateAppleFile(bookProjectDir, bookVersion)
	case "pdf":
		genetatePdfFile(bookProjectDir, bookVersion, false)
	case "print":
		genetatePdfFile(bookProjectDir, bookVersion, true)
	case "all", "":
		genetateAzw3File(bookProjectDir, bookVersion)
		genetateEpubFile(bookProjectDir, bookVersion)
		genetateMobiFile(bookProjectDir, bookVersion)
		genetateAppleFile(bookProjectDir, bookVersion)
		genetatePdfFile(bookProjectDir, bookVersion, false)
		genetatePdfFile(bookProjectDir, bookVersion, true)
	default:
		log.Fatal("Unknown target:", target)
	}
}

