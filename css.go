package main

/*
const EpubCSS = `
.content {
       font-size: 15px;
       line-height: 150%;
}

:not(pre) > code {
       padding: 1px 2px;
}

code {
       background-color: #dddddd;
}

pre {
       background-color: #dddddd;
       padding: 3px 6px;
       margin-left: 0px;
       margin-right: 0px;
}

table.table-bordered {
       border-collapse: collapse;
       border: 1px solid #999999;
}

table.table-bordered td {
      border: 1px solid black;
}

table.table-bordered th {
      border: 1px solid black;
}

.text-center {
       text-align: center;
}

.text-left {
       text-align: left;
}

a[href*='//']::after {
	content: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAQElEQVR42qXKwQkAIAxDUUdxtO6/RBQkQZvSi8I/pL4BoGw/XPkh4XigPmsUgh0626AjRsgxHTkUThsG2T/sIlzdTsp52kSS1wAAAABJRU5ErkJggg==);
	margin: 0 3px 0 5px;
}
`
*/

const CommonCSS = `

body {
       line-height: 1.1;
}

.text-center {
       text-align: center;
}

.text-left {
       text-align: left;
}

table.table-bordered {
       border-collapse: collapse;
       border: 1px solid #999999;
}

table.table-bordered td {
      border: 1px solid black;
}

table.table-bordered th {
      border: 1px solid black;
}

div.alert {
       background-color: #dddddd;
       margin-left: 12pt;
       margin-right: 12pt;
	padding: 3pt 10pt;
}

blockquote {
       background-color: #dddddd;
       margin-left: 12pt;
       margin-right: 12pt;
	padding: 5pt 10pt 5pt 16pt;
}

:not(pre) > code {
       padding: 1px 2px;
}

code {
       background-color: #dddddd;
}

pre {
       background-color: #dddddd;
       line-height: 1;
}

pre.line-numbers {
	counter-reset: line;
}

pre.line-numbers > code {
	counter-increment: line;
}

pre.line-numbers > code:before {
	content: counter(line);
	display: inline-block;
	text-align:right;
	width: 21pt;
	padding: 0 2pt 0 0;
	margin: 0 4pt 0 2pt;
	border-right: 1px solid #333;
	user-select: none;
	-webkit-user-select: none;
	-moz-user-select: none;
	-ms-user-select: none;
}
`




const Awz3CSS = CommonCSS + `

pre > code {
       font-family: Futura, "Caecilia Condensed", Courier;
       font-size: 7pt;
       text-align: left;
       line-height: 1;
}

pre.fixed-width > code {
       font-family: Courier, Futura, "Caecilia Condensed";
}

pre.fixed-width {
       font-family: Courier, Futura, "Caecilia Condensed";
}

pre {
       padding: 3px 3px;
       margin-left: 0px;
       margin-right: 0px;
}

span.invisible {
	visibility:hidden
}

`

const EpubCSS = CommonCSS + `

pre > code {
       text-align: left;
       line-height: 150%;
	tab-size: 7;
	-moz-tab-size: 7;
}

pre {
       padding: 3px 6px;
       margin-left: 0px;
       margin-right: 0px;
}

a[href*='//']::after {
	content: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAQElEQVR42qXKwQkAIAxDUUdxtO6/RBQkQZvSi8I/pL4BoGw/XPkh4XigPmsUgh0626AjRsgxHTkUThsG2T/sIlzdTsp52kSS1wAAAABJRU5ErkJggg==);
	margin: 0 3px 0 5px;
}

h1 {
	font-size: 300%;
}

h3 {
	font-size: 182%;
}

h4 {
	font-size: 128%;
	border-left: 3px solid #333;
	padding-left: 3px;
}

.content {
       font-size: 15px;
       line-height: 150%;
}
`

const PdfCommonCSS = CommonCSS + `

pre > code {
       text-align: left;
       line-height: 150%;
	tab-size: 7;
	-moz-tab-size: 7;
}

pre {
       padding: 3px 6px;
       margin-left: 0px;
       margin-right: 0px;
}


h1, h3, h4 {
	font-weight: bold;
       line-height: 110%;
}

h1 {
	font-size: 300%;
}

h3 {
	font-size: 182%;
}

h4 {
	font-size: 128%;
	border-left: 3px solid #333;
	padding-left: 3px;
}

.content {
       font-size: 15px;
       line-height: 150%;
}
`

const PdfCSS = PdfCommonCSS + `

a[href*='//']::after {
	content: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAQElEQVR42qXKwQkAIAxDUUdxtO6/RBQkQZvSi8I/pL4BoGw/XPkh4XigPmsUgh0626AjRsgxHTkUThsG2T/sIlzdTsp52kSS1wAAAABJRU5ErkJggg==);
	margin: 0 3px 0 5px;
}
`

const PrintCSS = PdfCommonCSS + `

h3 {
	padding-bottom: 2px;
	border-bottom: 2px solid #333;
	padding-left: 3px;
	border-left: 6px solid #333;
}
`

const MobiCSS = `

.text-center {
       text-align: center;
}

.text-left {
       text-align: left;
}

`
