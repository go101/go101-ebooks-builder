This project is used to produce Go 101 eBooks.
The code is this repository is ugly and full of bad practices.
Don't expect to learn something from the code.

```
Program options:
-book-project-dir : the path to Go 101 book project.
-book-version : the version of the book, presented in the names of output files.
-target : output format. [epub | azw3 | mobi | apple | pdf | print | all]
```

For any format, an epub file will be produced firstly,
then the epub version will be converted to the target format
by using the calibre GUI or command line tools.
So please install the calibre software before run this program.


Some examples:

```
$ export BookVersion=v1.12.c.7
$ export BookProjectDir=/home/go/src/github.com/go101/go101

$ go run *.go -target=epub -book-version=$BookVersion -book-project-dir=$BookProjectDir
$ go run *.go -target=azw3 -book-version=$BookVersion -book-project-dir=$BookProjectDir
$ go run *.go -target=apple -book-version=$BookVersion -book-project-dir=$BookProjectDir
$ go run *.go -target=print -book-version=$BookVersion -book-project-dir=$BookProjectDir
```
