package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func confirmBookProjectName(bookProjectDir string) string {
	checkFileExistence := func(filename string) bool {
		info, err := os.Stat(filepath.Join(bookProjectDir, filename))
		return err == nil && !info.IsDir()
	}
	if checkFileExistence("go101.go") {
		return "Go101"
	}
	if checkFileExistence("golang101.go") {
		return "Golang101"
	}
	return ""
}

func mustCreateTempFile(pattern string, content []byte) string {
	tmpfile, err := ioutil.TempFile("", pattern)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}

func mustParseImageData(s string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}

func runShellCommand(rootPath, cmd string, args ...string) {
	runShellCommand2(false, rootPath, cmd, args...)
}

func runShellCommand2(needOutput bool, rootPath, cmd string, args ...string) []byte {
	log.Println(append([]string{cmd}, args...))

	command := exec.Command(cmd, args...)
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Dir = rootPath

	if needOutput {
		o, err := command.Output()
		if err != nil {
			log.Fatal(err)
		}
		return o
	} else {
		command.Stdout = os.Stdout
		if err := command.Run(); err != nil {
			log.Fatal(err)
		}
		return nil
	}
}

type Article struct {
	Filename string
	Title    string
	Content  []byte

	chapter, chapter2 string
	internalFilename  []byte
}

//const ArticlesFolder = "articles"
const ArticlesFolder = "fundamentals"

func mustArticles(root string, engVersion bool) (index *Article, articles []*Article, chapterMapping map[string]*Article) {
	index = mustArticle(engVersion, -1, root, "pages", ArticlesFolder, "101.html")
	articles, chapterMapping = must101Articles(root, index, engVersion)
	//for _, a := range articles {
	//	log.Println(a.Title)
	//}
	return
}

// The last path token is the filename.
func mustArticle(engVersion bool, chapterNumber int, root string, pathTokens ...string) *Article {
	path := filepath.Join(root, filepath.Join(pathTokens...))
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("read file("+path+") error:", err)
	}

	title := retrieveArticleTitle(content)
	if title == "" {
		log.Fatalln("title not found in file(" + path + ")")
	}

	chapter, chapter2 := "", ""
	if engVersion {
		chapter = fmt.Sprintf(" (§%d)", chapterNumber)
		chapter2 = fmt.Sprintf("§%d. ", chapterNumber)
		title = fmt.Sprintf("§%d. %s", chapterNumber, title)
	} else {
		chapter = fmt.Sprintf("（第%d章）", chapterNumber)
		chapter2 = fmt.Sprintf("第%d章：", chapterNumber)
		title = fmt.Sprintf("第%d章：%s", chapterNumber, title)
	}

	return &Article{
		Filename: pathTokens[len(pathTokens)-1],
		Title:    title,
		Content:  content,

		chapter:  chapter,
		chapter2: chapter2,
	}
}

const MaxLen = 256

var H1, _H1 = []byte("<h1"), []byte("</h1>")
var TagSigns = [2]rune{'<', '>'}

func retrieveArticleTitle(content []byte) string {
	j, i := -1, bytes.Index(content, H1)
	if i < 0 {
		return ""
	}

	i += len(H1)
	i2 := bytes.IndexByte(content[i:i+MaxLen], '>')
	if i2 < 0 {
		return ""
	}
	i += i2 + 1

	j = bytes.Index(content[i:i+MaxLen], _H1)
	if j < 0 {
		return ""
	}

	//return string(content[i-len(H1) : i+j+len(_H1)])
	//return string(content[i : i+j])

	title := string(bytes.TrimSpace(content[i : i+j]))
	k, s := 0, make([]rune, 0, MaxLen)
	for _, r := range title {
		if r == TagSigns[k] {
			k = (k + 1) & 1
		} else if k == 0 {
			s = append(s, r)
		}
	}
	return string(s)
}

var Anchor, _Anchor, LineToRemoveTag, endl = []byte(`<li><a class="index" href="`), []byte(`">`), []byte(`(to remove)`), []byte("\n")
var IndexContentStart, IndexContentEnd = []byte(`<!-- index starts (don't remove) -->`), []byte(`<!-- index ends (don't remove) -->`)

func must101Articles(root string, indexArticle *Article, engVersion bool) (articles []*Article, chapterMapping map[string]*Article) {
	chapterMapping = make(map[string]*Article)

	content := indexArticle.Content
	i := bytes.Index(content, IndexContentStart)
	if i < 0 {
		log.Fatalf("%s not found", IndexContentStart)
	}
	i += len(IndexContentStart)
	content = content[i:]
	i = bytes.Index(content, IndexContentEnd)
	if i >= 0 {
		content = content[:i]
	}

	var buf bytes.Buffer
	for range [1000]struct{}{} {
		i = bytes.Index(content, LineToRemoveTag)
		if i < 0 {
			break
		}
		start := bytes.LastIndex(content[:i], endl)
		if start >= 0 {
			buf.Write(content[:start])
		}
		end := bytes.Index(content[i:], endl)
		content = content[i:]
		if end < 0 {
			end = len(content)
		}
		content = content[end:]
	}
	buf.Write(content)

	// modify index content
	indexArticle.Content = buf.Bytes()

	// find all articles from links
	content = indexArticle.Content
	chapter := 0
	for range [1000]struct{}{} {
		i = bytes.Index(content, Anchor)
		if i < 0 {
			break
		}
		content = content[i+len(Anchor):]
		i = bytes.Index(content, _Anchor)
		if i < 0 {
			break
		}

		article := mustArticle(engVersion, chapter, root, "pages", ArticlesFolder, string(content[:i]))
		articles = append(articles, article)
		chapter++

		if i := strings.Index(article.Filename, ".html"); i >= 0 {
			filename := article.Filename
			internalFilename := []byte(strings.ReplaceAll(filename, ".html", ".xhtml"))
			//internalFilename := href[:i]
			if internalFilename[0] >= '0' && internalFilename[0] <= '9' {
				internalFilename = append([]byte("go"), internalFilename...)
			}

			article.internalFilename = internalFilename

			chapterMapping[article.Filename] = article
			//log.Println(article.Filename, ":", string(article.internalFilename))
		}

		content = content[i+len(_Anchor):]
	}

	return
}

/*
func collectInternalArticles(content []byte) map[string][]byte {
	var m = make(map[string][]byte)

	for range [1000]struct{}{} {
		aStart := find(content, 0, A)
		if aStart < 0 {
			break
		}
		index := aStart+len(A)

		end := find(content, index, []byte(">"))
		if end < 0 {
			fatalError("a tag is incomplete", content[aStart:])
		}
		end++

		hrefStart := find(content[:end], index, Href)
		if hrefStart < 0 {
			//fatalError("a tag has not href", content[aStart:])
			content = content[end:]
			continue
		}
		hrefStart += len(Href)

		quotaStart := find(content[:end], hrefStart, []byte(`"`))
		if quotaStart < 0 {
			fatalError("a href is incomplete", content[aStart:])
		}
		quotaStart++

		quotaEnd := find(content[:end], quotaStart, []byte(`"`))
		if quotaEnd < 0 {
			fatalError("a href is incomplete", content[aStart:])
		}

		href := bytes.TrimSpace(content[quotaStart:quotaEnd])
		if bytes.HasPrefix(href, []byte("http")) {

		} else if i := bytes.Index(href, []byte(".html")); i >= 0 {
			filename := string(href[:i+len(".html")])
			internalFilename := []byte(strings.ReplaceAll(filename, ".html", ".xhtml"))
			//internalFilename := href[:i]
			if internalFilename[0] >= '0' && internalFilename[0] <= '9' {
				internalFilename = append([]byte("go"), internalFilename...)
			}
			m[filename] = internalFilename
		}

		content = content[end:]
	}

	return m
}
*/

func find(content []byte, start int, s []byte) int {
	i := bytes.Index(content[start:], s)
	if i >= 0 {
		i += start
	}
	return i
}

func fatalError(err string, content []byte) {
	n := len(content)
	if n > 100 {
		n = 100
	}
	log.Fatalln(err, ":", string(content[:n]))
}

func wrapContentDiv(articles []*Article) {
	for _, article := range articles {
		var buf = bytes.NewBuffer(make([]byte, 0, len(article.Content)+100))
		buf.WriteString(`<div class="content">`)
		buf.Write(article.Content)
		buf.WriteString(`</div>`)

		article.Content = buf.Bytes()
	}
}

func calLineWidth(line string) int {
	var n int
	var lastI int
	for i, r := range line {
		if r == '&' {
			n -= 3
		}
		if k := i - lastI; k > 1 {
			n += 2
		} else if k == 1 {
			n++
		}
		lastI = i
	}
	if k := len(line) - lastI; k > 1 {
		n += 2
	} else if k == 1 {
		n++
	}
	return n
}

var Pre, _Pre = []byte(`<pre`), []byte(`</pre>`)
var Code, _Code = []byte(`<code`), []byte(`</code>`)
var tabSpaces = strings.Repeat(" ", 3)

type Commend struct {
	slashStart int
	numSpaces  int
}

var commends [1024]Commend

func escapeCharactorWithinCodeTags(articles []*Article, target string) {
	for _, article := range articles {
		content := article.Content
		var buf = bytes.NewBuffer(make([]byte, 0, len(content)+10000))
		for range [1000]struct{}{} {
			preStart := find(content, 0, Pre)
			if preStart < 0 {
				break
			}

			index := preStart + len(Pre)
			preClose := find(content, index, _Pre)
			if preClose < 0 {
				fatalError("pre tag is incomplete in article:"+article.Filename, content[index:])
			}
			preEnd := preClose + len(_Pre)

			codeStart := find(content[:preClose], index, Code)
			if codeStart < 0 {
				buf.Write(content[:preEnd])
				content = content[preEnd:]
				continue
			}

			programStart := find(content[:preClose], codeStart, []byte(">"))
			if programStart < 0 {
				fatalError("code start tag is incomplete in article:"+article.Filename, content[codeStart:])
			}
			programStart++

			codeClose := find(content[:preClose], programStart, _Code)
			if codeClose < 0 {
				fatalError("code tag doesn't match in article:"+article.Filename, content[:programStart])
			}

			codeCloseEnd := find(content[:preClose], codeClose, []byte(">"))
			if codeCloseEnd < 0 {
				fatalError("code close tag is incomplete in article:"+article.Filename, content[:preClose])
			}
			codeCloseEnd++

			temp := string(content[programStart:codeClose])

			// Cancelled for the experience of copy-paste from pdf is bad.
			// At least, it is a little better to paste spaces than to paste nothing.
			//if target != "pdf" && target != "print" {
			// for pdf, keep tabs
			temp = strings.ReplaceAll(temp, "\t", tabSpaces)
			//}

			temp = strings.ReplaceAll(temp, "&amp;", "&")
			temp = strings.ReplaceAll(temp, "&lt;", "<")
			temp = strings.ReplaceAll(temp, "&gt;", ">")

			temp = strings.ReplaceAll(temp, "&", "&amp;")
			temp = strings.ReplaceAll(temp, "<", "&lt;")
			temp = strings.ReplaceAll(temp, ">", "&gt;")

			var mustLineNumbers = bytes.Index(content[preStart:codeStart], []byte("must-line-numbers")) >= 0
			var disableLineNumbers = bytes.Index(content[preStart:codeStart], []byte("disable-line-numbers111")) >= 0 ||
				bytes.Index(content[preStart:codeStart], []byte("must-not-line-numbers-on-kindle")) >= 0
			if disableLineNumbers {
				buf.Write(content[:preStart])
				buf.Write(bytes.ReplaceAll(content[preStart:codeStart], []byte("line-numbers"), []byte("xxx-yyy")))
				buf.Write(content[codeStart:programStart])
				buf.WriteString(temp)
				buf.Write(content[codeClose:preEnd])
			} else {
				switch target {
				case "epub", "apple", "pdf", "print":
					mustLineNumbers = true
					fallthrough
				case "azw3":
					mustLineNumbers = true // still better to have line numbers
					if mustLineNumbers {
						buf.Write(content[:codeStart])

						lines := strings.Split(temp, "\n")
						if strings.TrimSpace(lines[len(lines)-1]) == "" {
							lines = lines[:len(lines)-1]
						}
						if strings.TrimSpace(lines[0]) == "" {
							lines = lines[1:]
						}
						for _, line := range lines {
							buf.WriteString("<code>")
							if len(line) > 0 && line[len(line)-1] == '\r' {
								line = line[:len(line)-1]
							}
							buf.WriteString(line)
							buf.WriteString("</code>\n")

							if n := calLineWidth(line); n > 62 {
								log.Println("  ", n, ":", line)
							}
						}

						buf.Write(content[codeCloseEnd:preEnd])
					} else {
						buf.Write(content[:preStart])
						buf.Write(bytes.ReplaceAll(content[preStart:codeStart], []byte("line-numbers"), []byte("xxx-yyy")))
						buf.Write(content[codeStart:programStart])
						buf.WriteString(temp)
						buf.Write(content[codeClose:preEnd])
					}

				case "mobi":
					buf.Write(content[:programStart])
					lines := strings.Split(temp, "\n")
					if strings.TrimSpace(lines[len(lines)-1]) == "" {
						lines = lines[:len(lines)-1]
					}
					if strings.TrimSpace(lines[0]) == "" {
						lines = lines[1:]
					}
					for i, line := range lines {
						if mustLineNumbers {
							fmt.Fprintf(buf, "%3d. ", i+1)
						}
						if len(line) > 0 && line[len(line)-1] == '\r' {
							line = line[:len(line)-1]
						}
						buf.WriteString(line)
						buf.WriteString("\n")

						if n := calLineWidth(line); n > 62 {
							log.Println("  ", n, ":", line)
						}
					}
					buf.Write(content[codeClose:preEnd])
					/*
						default: // LienNumbers_Manually // epub
							linecount := strings.Count(temp, "\n")
							if linecount > 0 && temp[len(temp)-1] != '\n' {
								linecount++
							}
							var b strings.Builder
							for i := 1; i <= linecount; i++ {
								fmt.Fprintf(&b, "%d.", i)
								if i < linecount {
									fmt.Fprint(&b, "\n")
								}
							}
							buf.Write(content[:preStart])
							buf.WriteString(`
								<table class="table code" style="border-spacing:1px; border-collapse: collapse; width: 100%;">
								<tr>
								<td style="text-align: right; width: auto;">
								<pre class="line-numbers"><code class="language-go">`)
							buf.WriteString(b.String())
							buf.WriteString(`</code></pre>
								</td>
								<td style="width: 100%;">
							`)
							buf.Write(content[preStart:programStart])
							buf.WriteString(temp)
							buf.Write(content[codeClose:preEnd])
							buf.WriteString(`
								</td>
								</tr>
								</table>
							`)
					*/
				}
			}

			content = content[preEnd:]
		}
		buf.Write(content)
		article.Content = buf.Bytes()
	}
}

var Img, Src = []byte(`<img`), []byte(`src`)

func replaceImageSources(articles []*Article, imagePaths map[string]string, rewardImage string) {
	for _, article := range articles {
		content := article.Content
		var buf = bytes.NewBuffer(make([]byte, 0, len(content)+10000))
		for range [1000]struct{}{} {
			imgStart := find(content, 0, Img)
			if imgStart < 0 {
				break
			}
			index := imgStart + len(Img)

			end := find(content, index, []byte(">"))
			if end < 0 {
				fatalError("img tag is incomplete in article:"+article.Filename, content[imgStart:])
			}
			end++

			srcStart := find(content[:end], index, Src)
			if srcStart < 0 {
				fatalError("img tag has not src in article:"+article.Filename, content[imgStart:])
			}
			srcStart += len(Src)

			quotaStart := find(content[:end], srcStart, []byte(`"`))
			if quotaStart < 0 {
				fatalError("img src is incomplete in article:"+article.Filename, content[imgStart:])
			}
			quotaStart++

			quotaEnd := find(content[:end], quotaStart, []byte(`"`))
			if quotaEnd < 0 {
				fatalError("img src is incomplete in article:"+article.Filename, content[imgStart:])
			}

			src := bytes.TrimSpace(content[quotaStart:quotaEnd])
			newSrc := imagePaths[string(src)]
			if newSrc == "" {
				log.Fatalf("%s has no image path", src)
			}

			buf.Write(content[:quotaStart])
			buf.WriteString(newSrc)
			buf.Write(content[quotaEnd:end])

			content = content[end:]
		}
		buf.Write(content)

		if rewardImage != "" { // Go语言101
			fmt.Fprintf(buf, `
				<hr/>
				<div style="margin: 16px 0px; text-align: center;">
				<div>本书由<a href="https://gfw.tapirgames.com">老貘</a>历时三年写成。目前本书仍在不断改进和增容中。你的赞赏是本书和Go101.org网站不断增容和维护的动力。</div>
				<img src="%s" alt="赞赏"></img>
				<div>（请搜索关注微信公众号“Go 101”或者访问<a href="https://github.com/golang101/golang101">github.com/golang101/golang101</a>获取本书最新版）</div>
				</div>
			`, imagePaths[rewardImage])
		} else { // Go 101
			fmt.Fprintf(buf, `
				<hr/>
				<div style="margin: 16px 50px; text-align: center;">
				<div>(The <b>Go 101</b> book is still being improved frequently from time to time.
				Please visit <a href="https://go101.org">go101.org</a> or follow
				<a href="https://x.com/zigo_101">@zigo_101</a>
				to get the latest news of this book. BTW, Tapir,
				the author of the book, has developed several fun games.
				You can visit <a href="https://www.tapirgames.com/">tapirgames.com</a>
				to get more information about these games. Hope you enjoy them.)</div>
				</div>
			`)
		}

		article.Content = buf.Bytes()
	}

}

var A, _A, Href = []byte(`<a`), []byte(`</a>`), []byte(`href`)

var linkFmtPatterns = map[bool]string{true: " (%s)", false: "（%s）"}

func replaceInternalLinks(articles []*Article, chapterMapping map[string]*Article, bookWebsite string, forPrint, engVersion bool) {
	for _, article := range articles {
		content := article.Content
		var buf = bytes.NewBuffer(make([]byte, 0, len(content)+10000))
		for range [1000]struct{}{} {
			aStart := find(content, 0, A)
			if aStart < 0 {
				break
			}
			index := aStart + len(A)

			aClose := find(content, index, _A)
			if aClose < 0 {
				fatalError("a href is not closed in article:"+article.Filename, content[aStart:])
			}
			aEnd := aClose + len(_A)

			//openEnd := find(content, index, []byte(">"))
			//if openEnd < 0 {
			//	fatalError("a tag is incomplete in article:" + article.Filename, content[aStart:])
			//}
			//openEnd++

			hrefStart := find(content[:aEnd], index, Href)
			if hrefStart < 0 {
				//fatalError("a tag has not href in article:" + article.Filename, content[aStart:])
				buf.Write(content[:aEnd])
				content = content[aEnd:]
				continue
			}
			hrefStart += len(Href)

			quotaStart := find(content[:aEnd], hrefStart, []byte(`"`))
			if quotaStart < 0 {
				fatalError("a href is incomplete in article:"+article.Filename, content[aStart:])
			}
			quotaStart++

			quotaEnd := find(content[:aEnd], quotaStart, []byte(`"`))
			if quotaEnd < 0 {
				fatalError("a href is incomplete in article:"+article.Filename, content[aStart:])
			}

			href := bytes.TrimSpace(content[quotaStart:quotaEnd])
			done := false
			if bytes.HasPrefix(href, []byte("http")) {
				//buf.Write(content[:aEnd])
			} else if i := bytes.Index(href, []byte(".html")); i >= 0 {
				done = true

				var newHref []byte
				filename := string(href[:i+len(".html")])

				linkArticle := chapterMapping[filename]
				if linkArticle != nil {
					//newHref = bytes.ReplaceAll(href, []byte(".html"), []byte(internalName))
					newHref = linkArticle.internalFilename
				} else {
					//log.Println("internal url path", filename, "not found!")
					panic("internal url path " + filename + " not found!")
					newHref = append([]byte(bookWebsite+"/article/"), href...)
				}

				if article.Filename == "101.html" {
					buf.Write(content[:aStart])
					buf.WriteString(linkArticle.chapter2)
					buf.Write(content[aStart:quotaStart])
					buf.Write(newHref)
					buf.Write(content[quotaEnd:aEnd])
				} else {
					buf.Write(content[:quotaStart])
					buf.Write(newHref)
					buf.Write(content[quotaEnd:aEnd])
					buf.WriteString(linkArticle.chapter)
				}
			} else {
				//buf.Write(content[:aEnd])
			}
			if !done {
				if forPrint {
					buf.Write(content[:aEnd])
					fmt.Fprintf(buf, linkFmtPatterns[engVersion], href)
				} else {
					buf.Write(content[:aEnd])
				}
			}

			content = content[aEnd:]
		}
		buf.Write(content)
		article.Content = buf.Bytes()
	}

}

/*
  :not(pre) > code {
  =>
  <code bgcolor=#dddddd vspace=1 hspace=1></code>


  //pre {
  //=>
  //<pre vspace=5></pre>
*/
func setHtml32Atributes(articles []*Article) {

	for _, article := range articles {
		content := article.Content
		var buf = bytes.NewBuffer(make([]byte, 0, len(content)+10000))
		for range [1000]struct{}{} {
			codeStart := find(content, 0, Code)
			if codeStart < 0 {
				break
			}
			preStart := find(content, 0, Pre)
			if preStart >= 0 && preStart < codeStart {
				index := preStart + len(Pre)
				preClose := find(content, index, _Pre)
				if preClose < 0 {
					fatalError("pre tag is incomplete in article:"+article.Filename, content[index:])
				}
				preEnd := preClose + len(_Pre)

				buf.Write(content[:index])
				//buf.WriteString(` vspace=5`)
				buf.Write(content[index:preEnd])

				content = content[preEnd:]
				continue
			}

			index := codeStart + len(Code)

			codeClose := find(content, index, _Code)
			if codeClose < 0 {
				fatalError("code tag doesn't match in article:"+article.Filename, content[:codeStart])
			}
			codeEnd := codeClose + len(_Code)

			buf.Write(content[:index])
			buf.WriteString(` bgcolor="#dddddd" vspace="1" hspace="1"`)
			buf.Write(content[index:codeEnd])

			content = content[codeEnd:]
		}
		buf.Write(content)
		article.Content = buf.Bytes()
	}
}

var Kindle, _Kindle, CommentEnd = []byte("kindle starts:"), []byte("kindle ends:"), []byte("-->")

func filterArticles(content []byte, bookId int) []byte {
	var buf = bytes.NewBuffer(make([]byte, 0, len(content)+10000))

	for range [1000]struct{}{} {
		kindleStart := bytes.Index(content, Kindle)
		if kindleStart < 0 {
			break
		}

		index := kindleStart + len(Kindle)

		startEnd := find(content, index, CommentEnd)
		if startEnd < 0 {
			fatalError("kindle tag is imcomplete", content[:kindleStart])
		}

		kindleEnd := find(content, startEnd+len(CommentEnd), _Kindle)
		if kindleEnd < 0 {
			fatalError("kindle tag doesn't match a", content[:startEnd])
		}

		endEnd := find(content, kindleEnd+len(_Kindle), CommentEnd)
		if endEnd < 0 {
			fatalError("kindle tag doesn't match b", content[:kindleEnd])
		}
		endEnd += len(CommentEnd)

		idStr := string(bytes.TrimSpace(content[index:startEnd]))
		n, err := strconv.Atoi(idStr)
		if err != nil {
			fatalError("bad kindle book id: "+idStr+". "+err.Error(), content[:endEnd])
		}

		if n == bookId {
			buf.Write(content[:endEnd])
		} else {
			buf.Write(content[:kindleStart])
		}
		content = content[endEnd:]
	}
	buf.Write(content)
	return buf.Bytes()
}

func hintExternalLinks(articles []*Article, externalLinkPngPath string) {
	img := `<img src="` + externalLinkPngPath + `" width="20" height="20"></img>`

	for _, article := range articles {
		content := article.Content
		var buf = bytes.NewBuffer(make([]byte, 0, len(content)+10000))
		for range [1000]struct{}{} {
			aStart := find(content, 0, A)
			if aStart < 0 {
				break
			}
			index := aStart + len(A)

			end := find(content, index, []byte(">"))
			if end < 0 {
				fatalError("a tag is incomplete in article:"+article.Filename, content[aStart:])
			}
			end++

			hrefStart := find(content[:end], index, Href)
			if hrefStart < 0 {
				//fatalError("a tag has not href in article:" + article.Filename, content[aStart:])
				buf.Write(content[:end])
				content = content[end:]
				continue
			}
			hrefStart += len(Href)

			quotaStart := find(content[:end], hrefStart, []byte(`"`))
			if quotaStart < 0 {
				fatalError("a href is incomplete in article:"+article.Filename, content[aStart:])
			}
			quotaStart++

			quotaEnd := find(content[:end], quotaStart, []byte(`"`))
			if quotaEnd < 0 {
				fatalError("a href is incomplete in article:"+article.Filename, content[aStart:])
			}

			aEnd := find(content, index, _A)
			if aEnd < 0 {
				fatalError("a tag doesn't match in article:"+article.Filename, content[index:])
			}
			endEnd := aEnd + len(_A)

			href := bytes.TrimSpace(content[quotaStart:quotaEnd])
			if bytes.HasPrefix(href, []byte("http")) {
				buf.Write(content[:aEnd])
				buf.WriteString(img)
				buf.Write(content[aEnd:endEnd])
			} else {
				buf.Write(content[:endEnd])
			}

			content = content[endEnd:]
		}
		buf.Write(content)
		article.Content = buf.Bytes()
	}
}

var (
	attribtuesTobeRemoved = [][]byte{
		[]byte(` valign="bottom"`),
		[]byte(` valign="middle"`),
		[]byte(` align="center"`),
		[]byte(` align="left"`),
		[]byte(` border="1"`),
		[]byte(` scope="row"`),
	}
)

func removeXhtmlAttributes(articles []*Article) {

	for _, article := range articles {
		content := article.Content
		log.Println("===========", article.Title, len(content))
		for _, attr := range attribtuesTobeRemoved {
			content = bytes.ReplaceAll(content, attr, []byte{})
			log.Printf("%s: %d", attr, len(content))
		}
		article.Content = content
	}
}
