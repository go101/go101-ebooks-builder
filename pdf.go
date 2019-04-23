package main

import (
	"log"
	"os"

	"github.com/bmaupin/go-epub"
)

/*
https://chromedevtools.github.io/devtools-protocol/tot/Page/#method-printToPDF

landscape
displayHeaderFooter
printBackground
preferCSSPageSize

paperWidth | paperHeight

marginTop | marginBottom | marginLeft | marginRight

headerTemplate | footerTemplate

    HTML template for the print header. Should be valid HTML markup
    with following classes used to inject printing values into them:
    * date: formatted print date
    * title: document title
    * url: document location
    * pageNumber: current page number
    * totalPages: total pages in the document
    For example, <span class=title></span> would generate span containing the title.

*/

// --preferCSSPageSize doesn't work?
// --scale=0.67 doesn't work? Neither --scale=67, --scale=67%
// Looks all above don't work? Why?

// chromium --headless --disable-gpu --preferCSSPageSize --scale=67% --print-to-pdf=Go101-desktop.pdf http://localhost:55555/article/pdf-book101

// chromium --headless --disable-gpu --preferCSSPageSize --scale=67% --print-to-pdf=Golang101-desktop.pdf http://localhost:12345/article/pdf-book101


/*
It looks many options are not supported when isntall wkhhtmltopdf by apt.
The download version from the official website is ok.

wkhtmltopdf --javascript-delay 3000 --no-stop-slow-scripts --title "Go 101" --enable-internal-links --footer-center [page] --header-center [section] --header-line http://localhost:55555/article/pdf-book101 Go101-desktop-with-index.pdf

wkhtmltopdf --javascript-delay 3000 --no-stop-slow-scripts --title "Go 101" --enable-internal-links --footer-center [page] --header-center [section] --header-line http://localhost:55555/article/print-book101 Go101-desktop-with-index-2.pdf
		
wkhtmltopdf --javascript-delay 3000 --no-stop-slow-scripts --title "Go语言101" --enable-internal-links --footer-center [page] --header-center [section] --header-line http://localhost:12345/article/pdf-book101 Golang101-desktop-with-index.pdf

wkhtmltopdf --javascript-delay 3000 --no-stop-slow-scripts --title "Go语言101" --enable-internal-links --footer-center [page] --header-center [section] --header-line http://localhost:12345/article/print-book101 Golang101-desktop-with-index-2.pdf
*/

/*
func genetatePdfFile_Chrome(bookProjectDir, bookVersion string) {
	var outFilename string
	var pdfHtmlPage string
	
	projectName := confirmBookProjectName(bookProjectDir)
	switch projectName {
	default:
		log.Fatal("unknow book porject: ", projectName)
	case "Go101":
		outFilename = "Go101-" + bookVersion + ".pdf"
		pdfHtmlPage = "http://localhost:55555/article/pdf-book101"
	case "Golang101":
		outFilename = "Golang101-" + bookVersion + ".pdf"
		pdfHtmlPage = "http://localhost:12345/article/pdf-book101"
	}
	
	runShellCommand(".", "chromium", "--headless", "--disable-gpu",
			"--preferCSSPageSize", "--scale=67%",
			"--print-to-pdf=" + outFilename,
			pdfHtmlPage)
	log.Println("Create", outFilename, "done!")
}
*/

/*

// _Calibre
func genetatePdfFile(bookProjectDir, bookVersion string) {
	//epubFilename := genetateEpubFile(bookProjectDir, bookVersion)
	
	var epubFilename string
	var outFilename string
	
	projectName := confirmBookProjectName(bookProjectDir)
	switch projectName {
	default:
		log.Fatal("unknow book porject: ", projectName)
	case "Go101":
		epubFilename = "Go101-" + bookVersion + ".epub"
		outFilename = "Go101-" + bookVersion + ".pdf"
	case "Golang101":
		epubFilename = "Golang101-" + bookVersion + ".epub"
		outFilename = "Golang101-" + bookVersion + ".pdf"
	}
	
	// --prefer-metadata-cover
	// --remove-first-image
	// --pdf-add-toc
	
	runShellCommand(".", "ebook-convert", epubFilename, outFilename,
			"--prefer-metadata-cover",
			"--pdf-footer-template", `'<p style="text-align:center; font-size: small;">_PAGENUM_</p>'`,
			"--pdf-header-template", `'<p style="text-align:center; font-size: small;">_SECTION_</p>'`,
			"--paper-size", "a4",
			"--pdf-default-font-size", `14`,
			"--pdf-mono-font-size", `13`,
			"--pdf-page-margin-top", `36.0`,
			"--pdf-page-margin-bottom", `36.0`,
			"--pdf-page-margin-right", `36.0`,
			"--pdf-page-margin-left", `36.`,
			//"--pdf-page-number-map", `'if (n < 2) 0; else n - 2;'`,
	)
	log.Println("Create", outFilename, "done!")
}
*/

/*
<p style="text-align:center; font-size: small;">_SECTION_</p>

<p style="text-align:center; font-size: small;">_PAGENUM_</p>
*/



func genetatePdfFile(bookProjectDir, bookVersion string, forPrint bool) string {
	var e *epub.Epub
	var outFilename string
	var indexArticleTitle string
	var bookWebsite string
	var engVersion bool
	
	target := "pdf"
	css := PdfCSS
	ext := ".pdf.epub"
	if forPrint {
		target = "print"
		ext = ".print" + ext
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

	writeEpub_Go101(outFilename, e, -1, bookWebsite, projectName, indexArticleTitle, bookProjectDir, cssPath, target, engVersion)
	log.Println("Create", outFilename, "done!")
	
	return outFilename
}