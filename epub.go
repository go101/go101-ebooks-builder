package main

import (
	"archive/zip"
	//"bytes"
	"io"
	"log"
	"os"
	"strings"
	"path/filepath"

	"github.com/bmaupin/go-epub"
)

/*
.mobi file
* bgcolor
* <div class="text-center">
  =>
  <div align=center>
* vspace=1 hspace=1

:not(pre) > code {
       padding: 1px 2px;
       background-color: #dbdbdb;
}

pre {
       padding: 3px 6px;
       margin-left: 0px;
       margin-right: 0px;
}

ebook-convert "book.epub" "book.mobi"
*/

func genetateEpubFile(bookProjectDir, bookVersion string) string {
	var e *epub.Epub
	var outFilename string
	var indexArticleTitle string
	var bookWebsite string
	var engVersion bool
	
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
		outFilename = "Go101-" + bookVersion + ".epub"
	case "Golang101":
		e = epub.NewEpub("Go语言101")
		e.SetAuthor("老貘")
		indexArticleTitle = "目录"
		bookWebsite = "https://gfw.go101.org"
		engVersion = false
		outFilename = "Golang101-" + bookVersion + ".epub"
	}
	
	cssFilename := "all.css"
	tempCssFile := mustCreateTempFile("all*.css", []byte(EpubCSS))
	defer os.Remove(tempCssFile)
	cssPath, err := e.AddCSS(tempCssFile, cssFilename)
	if err != nil {
		log.Fatalln("add css", cssFilename, "failed:", err)
	}

	writeEpub_Go101(outFilename, e, -1, bookWebsite, projectName, indexArticleTitle, bookProjectDir, cssPath, "epub", engVersion)
	log.Println("Create", outFilename, "done!")
	
	return outFilename
}

const (
	LienNumbers_Manually      = iota
	LienNumbers_Unchange
	LienNumbers_Selectable
	LienNumbers_Automatically
)

// zero bookId means all.
func writeEpub_Go101(outputFilename string, e *epub.Epub, bookId int, bookWebsite, projectName, indexArticleTitle, bookProjectDir, cssPath, target string, engVersion bool) {
	imagePaths := addImages(e, bookProjectDir)
	var rewardImage string
	if projectName == "Golang101" {
		rewardImage = "res/101-reward-qrcode-2.png"
	}
	
	index, articles, chapterMapping := mustArticles(bookProjectDir, engVersion)
	index.Title = indexArticleTitle
	index.Content = append([]byte("<h1>" + index.Title + "</h1>"), index.Content...)
	
	if bookId > 0 {
		index.Content = filterArticles(index.Content, bookId)
	}
	//internalArticles := collectInternalArticles(index.Content)
	//log.Println("internalArticles:", internalArticles)	
	
	oldArticles := articles
	articles = nil
	for _, article := range oldArticles {
		if _, present := chapterMapping[article.Filename]; present {
			articles = append(articles, article)
		}
	}
	articles = append([]*Article{index}, articles...)

	escapeCharactorWithinCodeTags(articles, target)
	replaceInternalLinks(articles, chapterMapping, bookWebsite, target == "print", engVersion)
	replaceImageSources(articles, imagePaths, rewardImage)

	switch target {
	case "azw3":
		fallthrough
	case "mobi": // mobi
		setHtml32Atributes(articles)
		
		pngFilename := "external-link.png"
		tempPngFile := mustCreateTempFile("external-link*.png", mustParseImageData(ExternalLinkPNG))
		defer os.Remove(tempPngFile)
		imgpath, err := e.AddImage(tempPngFile, pngFilename)
		if err != nil {
			log.Fatalln("add image", pngFilename, "failed:", err)
		}
		imagePaths[pngFilename] = imgpath
		
		hintExternalLinks(articles, imgpath)
	case "apple":
		removeXhtmlAttributes(articles)
	default:
	}
	
	wrapContentDiv(articles)

	// ...
	e.SetCover(imagePaths["res/101-front-cover-1400x.jpg"], "")

	for _, article := range articles {
		internalFilename := string(article.internalFilename)
		e.AddSection(string(article.Content), article.Title, internalFilename, cssPath)
	}

	if err := e.Write(outputFilename); err != nil {
		log.Fatalln("write epub failed:", err)
	}
}

func addImages(e *epub.Epub, bookProjectDir string) map[string]string {
	imagePaths := make(map[string]string)

	root := filepath.Join(bookProjectDir, ArticlesFolder, "res")
	f := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != root {
			return filepath.SkipDir
		}
		//log.Printf("visited file or dir: %q\n", path)
		if !info.IsDir() {
			index := strings.Index(path, "res/")
			imgsrc := path[index:]
			filename := filepath.Base(path)
			lower := strings.ToLower(imgsrc)
			if strings.HasSuffix(lower, ".png") ||
				strings.HasSuffix(lower, ".gif") ||
				strings.HasSuffix(lower, ".jpg") ||
				strings.HasSuffix(lower, ".jpeg") {
				imgpath, err := e.AddImage(path, filename)
				if err != nil {
					log.Fatalln("add image", filename, "failed:", err)
				}
				imagePaths[imgsrc] = imgpath
				//log.Println(imgsrc, filename, imgpath)
			}
		}

		return nil
	}

	if err := filepath.Walk(root, f); err != nil {
		log.Fatalln("list article res image files error:", err)
	}

	return imagePaths
}

func removePageFromEpub(epubFilename string, pagesToRemove ...string) {
	r, err := zip.OpenReader(epubFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	
	os.Remove(epubFilename)
	
	outputFile, err := os.Create(epubFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	w := zip.NewWriter(outputFile)

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		//log.Printf("Contents of %s:\n", f.Name)
		if f.Name == "EPUB/xhtml/cover.xhtml" {
			continue
		}
		
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		
		of, err := w.Create(f.Name)
		if err != nil {
			log.Fatal(err)
		}
		
		_, err = io.Copy(of, rc)
		if err != nil {
			log.Fatal(err)
		}
		
		rc.Close()
		log.Println()
	}
	
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
	
	outputFile.Sync()
}
