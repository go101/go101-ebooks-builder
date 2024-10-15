package main

import (
	"archive/zip"
	//"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmaupin/go-epub"
)

func genetateEpubFile(bookProjectDir, bookVersion, coverImagePath string) string {
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

	writeEpub_Go101(outFilename, e, -1, bookWebsite, projectName, indexArticleTitle, bookProjectDir, coverImagePath, cssPath, "epub", engVersion)
	log.Println("Create", outFilename, "done!")

	return outFilename
}

const (
	LienNumbers_Manually = iota
	LienNumbers_Unchange
	LienNumbers_Selectable
	LienNumbers_Automatically
)

// zero bookId means all.
func writeEpub_Go101(outputFilename string, e *epub.Epub, bookId int, bookWebsite, projectName, indexArticleTitle, bookProjectDir, coverImagePath, cssPath, target string, engVersion bool) {
	imagePaths, coverImagePathInEpub := addImages(e, bookProjectDir, coverImagePath)
	var rewardImage string
	if projectName == "Golang101" {
		rewardImage = "res/101-reward-qrcode-8.png"
	}

	index, articles, chapterMapping := mustArticles(bookProjectDir, engVersion)
	index.Title = indexArticleTitle
	index.Content = append([]byte("<h1>"+index.Title+"</h1>"), index.Content...)

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
	e.SetCover(coverImagePathInEpub, "")

	for _, article := range articles {
		internalFilename := string(article.internalFilename)
		e.AddSection(string(article.Content), article.Title, internalFilename, cssPath)
	}

	if err := e.Write(outputFilename); err != nil {
		log.Fatalln("write epub failed:", err)
	}
}

func addImages(e *epub.Epub, bookProjectDir, coverImagePath string) (map[string]string, string) {
	imagePaths := make(map[string]string)

	root := filepath.Join(bookProjectDir, "pages", ArticlesFolder, "res")
	f := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != root {
			return filepath.SkipDir
		}
		//log.Printf("visited file or dir: %q\n", path)
		if !info.IsDir() {
			index := strings.Index(path, "res"+string(filepath.Separator))
			imgsrc := path[index:]
			filename := filepath.Base(path)
			lower := strings.ToLower(imgsrc)
			if strings.Index(filename, "front-cover") < 0 &&
				(strings.HasSuffix(lower, ".png") ||
					strings.HasSuffix(lower, ".gif") ||
					strings.HasSuffix(lower, ".jpg") ||
					strings.HasSuffix(lower, ".jpeg")) {
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

	// Cover image
	var coverImagePathInEpub string
	{
		filename := filepath.Base(coverImagePath)
		imgpath, err := e.AddImage(coverImagePath, filename)
		if err != nil {
			log.Fatalln("add cover image", filename, "failed:", err)
		}
		coverImagePathInEpub = imgpath
	}

	return imagePaths, coverImagePathInEpub
}

func removePagesFromEpub(epubFilename string, pagesToRemove ...string) {
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

	shouldRemove := map[string]bool{}
	for _, page := range pagesToRemove {
		shouldRemove[page] = true
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		//log.Printf("Contents of %s:\n", f.Name)
		if shouldRemove[f.Name] {
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
