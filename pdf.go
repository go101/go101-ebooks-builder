package main

import (
	"log"
	"os"
	"strings"

	"github.com/bmaupin/go-epub"
)

func genetatePdfFile(bookProjectDir, bookVersion, coverImagePath string, forPrint bool) string {
	var e *epub.Epub
	var outFilename string
	var indexArticleTitle string
	var bookWebsite string
	var engVersion bool

	target := "pdf"
	css := PdfCSS
	ext := ".pdf"
	if forPrint {
		target = "print"
		css = PrintCSS
	}

	projectName := confirmBookProjectName(bookProjectDir)
	switch projectName {
	default:
		log.Fatal("unknow book porject: ", projectName)
	case "Go101":
		e = epub.NewEpub("Go 101")
		e.SetAuthor("Tapir Liu")
		indexArticleTitle = "Contents"
		bookWebsite = "https://go101.org"
		engVersion = true
		outFilename = "Go101-" + bookVersion + ext
	case "Golang101":
		e = epub.NewEpub("Go语言101")
		e.SetAuthor("老貘")
		indexArticleTitle = "目录"
		bookWebsite = "https://gfw.go101.org"
		engVersion = false
		outFilename = "Golang101-" + bookVersion + ext
	}

	cssFilename := "all.css"
	tempCssFile := mustCreateTempFile("all*.css", []byte(css))
	defer os.Remove(tempCssFile)
	cssPath, err := e.AddCSS(tempCssFile, cssFilename)
	if err != nil {
		log.Fatalln("add css", cssFilename, "failed:", err)
	}

	// ...
	tempOutFilename := outFilename + "*.epub"
	tempOutFilename = mustCreateTempFile(tempOutFilename, nil)
	defer os.Remove(tempOutFilename)
	//tempOutFilename := outFilename + ".epub"

	writeEpub_Go101(tempOutFilename, e, -1, bookWebsite, projectName, indexArticleTitle, bookProjectDir, coverImagePath, cssPath, target, engVersion)

	removePagesFromEpub(tempOutFilename, "EPUB/xhtml/cover.xhtml")

	epub2pdf := func(serifFont, fontSize, inputFilename, outputFilename string) {
		conversionParameters := make([]string, 0, 32)
		pushParams := func(params ...string) {
			conversionParameters = append(conversionParameters, params...)
		}
		pushParams(inputFilename, outputFilename)
		pushParams("--toc-title", indexArticleTitle)
		pushParams("--pdf-header-template", `'<div style="text-align: center; padding-top: -9px; font-size: small;">_SECTION_</div>'`)
		pushParams("--pdf-footer-template", `'<div style="text-align: center; padding-top: -9px; font-size: small;">_PAGENUM_</div>'`)
		//pushParams("--pdf-page-numbers")
		pushParams("--paper-size", "a4")
		pushParams("--pdf-serif-family", serifFont)
		//pushParams("--pdf-sans-family", serifFont)
		pushParams("--pdf-mono-family", "Liberation Mono")
		pushParams("--pdf-default-font-size", fontSize)
		pushParams("--pdf-mono-font-size", "15")
		pushParams("--pdf-page-margin-top", "36")
		pushParams("--pdf-page-margin-bottom", "36")
		if forPrint {
			pushParams("--pdf-add-toc")
			pushParams("--pdf-page-margin-left", "72")
			pushParams("--pdf-page-margin-right", "72")
		} else {
			pushParams("--pdf-page-margin-left", "36")
			pushParams("--pdf-page-margin-right", "36")
		}
		pushParams("--preserve-cover-aspect-ratio")

		runShellCommand(".", "ebook-convert", conversionParameters...)

		log.Println("Create", outputFilename, "done!")
	}

	if forPrint {
		outFilenameForPrinting := strings.Replace(outFilename, ".pdf", ".pdf-ForPrinting.pdf", 1)
		if projectName == "Go101" {
			epub2pdf("Liberation Serif", "17", tempOutFilename, outFilenameForPrinting)
		} else if projectName == "Golang101" {
			epub2pdf("AR PL SungtiL GB", "16", tempOutFilename, outFilenameForPrinting)
		}
	} else {
		if projectName == "Go101" {
			epub2pdf("Liberation Serif", "17", tempOutFilename, outFilename)
		} else if projectName == "Golang101" {
			outFilenameKaiTi := strings.Replace(outFilename, ".pdf", ".pdf-KaiTi.pdf", 1)
			epub2pdf("AR PL KaitiM GB", "16", tempOutFilename, outFilenameKaiTi)

			outFilenameSongTi := strings.Replace(outFilename, ".pdf", ".pdf-SongTi.pdf", 1)
			epub2pdf("AR PL SungtiL GB", "16", tempOutFilename, outFilenameSongTi)
		}
	}

	return outFilename
}
