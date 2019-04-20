package main

import (
	"log"
	"os"

	"github.com/bmaupin/go-epub"
)

func genetateMobiFile(bookProjectDir, bookVersion string) {
	genetateMobiFileForBook(bookProjectDir, bookVersion, 0)
	
	//genetateMobiFileForBook(bookProjectDir, bookVersion, 1)
	//genetateMobiFileForBook(bookProjectDir, bookVersion, 2)
}

// zero bookId means all.
func genetateMobiFileForBook(bookProjectDir, bookVersion string, bookId int) {
	var e *epub.Epub
	var outFilename string
	var indexArticleTitle string
	var bookWebsite string
	
	projectName := confirmBookProjectName(bookProjectDir)
	switch projectName {
	default:
		log.Fatal("unknow book porject: ", projectName)
	case "Go101":
		if bookId == 0 {
			e = epub.NewEpub("Go 101")
			outFilename = "Go101-" + bookVersion + ".mobi"
		} else if bookId == 1 {
			e = epub.NewEpub("Go 101 (Type System)")
			outFilename = "Go101-" + bookVersion + "-types.mobi"
		} else if bookId == 2 {
			e = epub.NewEpub("Go 101 (Extended)")
			outFilename = "Go101-" + bookVersion + "-extended.mobi"
		} else {
			log.Fatal("unknown book id: ", bookId)
		}
		e.SetAuthor("Tapir Liu")
		bookWebsite = "https://go101.org"
		indexArticleTitle = "Contents"
	case "Golang101":
		if bookId == 0 {
			e = epub.NewEpub("Go语言101")
			outFilename = "Golang101-" + bookVersion + ".mobi"
		} else if bookId == 1 {
			e = epub.NewEpub("Go语言101（类型系统）")
			outFilename = "Golang101" + bookVersion + "-types.mobi"
		} else if bookId == 2 {
			e = epub.NewEpub("Go语言101（扩展阅读）")
			outFilename = "Golang101-" + bookVersion + "-extended.mobi"
		} else {
			log.Fatal("unknown book id: ", bookId)
		}
		e.SetAuthor("老貘")
		bookWebsite = "https://gfw.go101.org"
		indexArticleTitle = "目录"
	}
	
	cssFilename := "all.css"
	tempCssFile := mustCreateTempFile("all*.css", []byte(MobiCSS))
	defer os.Remove(tempCssFile)
	cssPath, err := e.AddCSS(tempCssFile, cssFilename)
	if err != nil {
		log.Fatalln("add css", cssFilename, "failed:", err)
	}
	
	//tempOutFilename := outFilename + "*.epub"
	//tempOutFilename = mustCreateTempFile(tempOutFilename, nil)
	//defer os.Remove(tempOutFilename)
	tempOutFilename := outFilename + ".epub"

	writeEpub_Go101(tempOutFilename, e, bookId, bookWebsite, projectName, indexArticleTitle, bookProjectDir, cssPath, "mobi")
	println("ebook-convert", tempOutFilename, outFilename)
	runShellCommand(".", "ebook-convert", tempOutFilename, outFilename)
	log.Println("Create", outFilename, "done!")
}
