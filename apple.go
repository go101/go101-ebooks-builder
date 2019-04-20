package main

import (
	"log"
	"os"

	"github.com/bmaupin/go-epub"
)

func genetateAppleFile(bookProjectDir, bookVersion string) string {
	var e *epub.Epub
	var outFilename string
	var indexArticleTitle string
	var bookWebsite string
	
	projectName := confirmBookProjectName(bookProjectDir)
	switch projectName {
	default:
		log.Fatal("unknow book porject: ", projectName)
	case "Go101":
		e = epub.NewEpub("Go 101")
		e.SetAuthor("Tapir Liu")
		indexArticleTitle = "Contents"
		bookWebsite = "https://go101.org"
		outFilename = "Go101-" + bookVersion + ".apple.epub"
	case "Golang101":
		e = epub.NewEpub("Go语言101")
		e.SetAuthor("老貘")
		indexArticleTitle = "目录"
		bookWebsite = "https://gfw.go101.org"
		outFilename = "Golang101-" + bookVersion + ".apple.epub"
	}
	
	cssFilename := "all.css"
	tempCssFile := mustCreateTempFile("all*.css", []byte(PdfCSS))
	defer os.Remove(tempCssFile)
	cssPath, err := e.AddCSS(tempCssFile, cssFilename)
	if err != nil {
		log.Fatalln("add css", cssFilename, "failed:", err)
	}

	writeEpub_Go101(outFilename, e, -1, bookWebsite, projectName, indexArticleTitle, bookProjectDir, cssPath, "apple")
	log.Println("Create", outFilename, "done!")
	
	return outFilename
}
