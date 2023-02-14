package convertor

import (
	"html"
	"regexp"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/interfaces"
)

const ImageLinkType = "image for download"
const ImagePathType = "local image"

var reImageExt = regexp.MustCompile(`(?i)\.(bmp|jpg|jpeg|gif|png|tiff)$`)

var reComment = regexp.MustCompile(`<!--.*?(-->|$)`)
var reTitle = regexp.MustCompile(`(?ims)<title.*?>(.*?)</title>`)
var reHead = regexp.MustCompile(`(?ims)<head>.*?</head>`)
var reScript = regexp.MustCompile(`(?ims)<script.*?</script>`)
var reCode = regexp.MustCompile(`(?ims)<code[ >].*?</code>`)
var reStyle = regexp.MustCompile(`(?ims)<style.*?</style>`)
var reBody = regexp.MustCompile(`(?ims).*?<body.*?>|</body>.*`)
var reTag = regexp.MustCompile(`(?ms)<(\w+)[ >]`)

var reEnter = regexp.MustCompile(`(?ms)\s*\n\s*`)
var reEnters = regexp.MustCompile(`(?ms)\n+`)
var ReSpaces = regexp.MustCompile(`(?ms) +`)
var reClearAllowTags = regexp.MustCompile(`(?ims)\s*</?(div|p|ul|br)>\s*`)

var reImage = regexp.MustCompile(`(?ims) src=["']([^"']+)`)
var reHTTP = regexp.MustCompile(`(?i)^http`)
var reSlash = regexp.MustCompile(`^/`)
var reDomain = regexp.MustCompile(`^[^/\?$]+`)
var reURL = regexp.MustCompile(`^(.*)[/\?$]`)
var reFilePath = regexp.MustCompile(`^.*/`)

var allowTagsWithAttributes = map[string]string{
	"p": "nothing",
	"a": "href",

	"br": "nothing",
	"li": "nothing",
	"ul": "nothing",
	"ol": "nothing",
	"tr": "nothing",
	"th": "nothing",
	"td": "nothing",
	"h1": "nothing",
	"h2": "nothing",
	"h3": "nothing",
	"h4": "nothing",
	"h5": "nothing",

	"img": "src",
	"div": "nothing",

	"table": "nothing",
}

func ClearHTML(url, file, content string) (string, string, []interfaces.Link) {

	title := file
	arr := reTitle.FindStringSubmatch(content)
	if len(arr) > 1 {
		title = arr[1]
	}
	content = reComment.ReplaceAllString(content, "")
	content = reHead.ReplaceAllString(content, "")
	content = reScript.ReplaceAllString(content, "")
	content = reCode.ReplaceAllString(content, "")
	content = reStyle.ReplaceAllString(content, "")
	content = reBody.ReplaceAllString(content, "")

	content = clearNotAllowTags(content)
	content = clearNotAllowAttributes(content)
	content = clearSpaces(content)
	content = html.UnescapeString(content)
	files := getPictures(url, file, content)
	return title, content, files

}

func getPictures(url, file, content string) []interfaces.Link {

	var files []interfaces.Link
	var domain, baseURL, filePath string
	ar := reDomain.FindAllString(url, 1)
	if len(ar) > 0 {
		domain = ar[0]
	}
	arr := reURL.FindAllStringSubmatch(url, 1)
	if len(arr) > 0 {
		baseURL = arr[0][1]
	}
	ar = reFilePath.FindAllString(file, 1)
	if len(ar) > 0 {
		filePath = ar[0]
	}
	arr = reImage.FindAllStringSubmatch(content, -1)
	for _, p := range arr {
		if p[1] != "" && reImageExt.MatchString(p[1]) {
			if reHTTP.MatchString(p[1]) {
				files = append(files, interfaces.Link{Name: p[1], To: p[1], Type: ImageLinkType})
			} else if url != "" {
				if reSlash.MatchString(p[1]) {
					files = append(files, interfaces.Link{Name: domain + p[1], To: p[1], Type: ImageLinkType})
				} else {
					files = append(files, interfaces.Link{Name: baseURL + "/" + p[1], To: p[1], Type: ImageLinkType})
				}
			} else {
				// get image as local
				files = append(files, interfaces.Link{Name: filePath + p[1], To: p[1], Type: ImagePathType})
			}
		}
	}
	return files
}

func clearNotAllowTags(content string) string {
	var allTags = map[string]int{}
	arr := reTag.FindAllStringSubmatch(content, -1)
	for i := range arr {
		allTags[arr[i][1]] = 1
	}
	for t := range allTags {
		_, ok := allowTagsWithAttributes[t]
		if !ok {
			//log.Println("IS tag for del:", t)
			reClearTag := regexp.MustCompile(`(?ims)<` + t + `[^>]*>`)
			content = reClearTag.ReplaceAllString(content, "")
			reClearTag = regexp.MustCompile(`</` + t + `>`)
			content = reClearTag.ReplaceAllString(content, "")
		}
	}

	return content
}

func clearNotAllowAttributes(content string) string {

	for t, a := range allowTagsWithAttributes {
		if a != "nothing" {
			reAttributes := regexp.MustCompile(`<` + t + `.*?` + a + `\s*=\s*["']([^"']+)[^>]*>`)
			content = reAttributes.ReplaceAllString(content, `<`+t+" "+a+"=\"${1}\">")
		} else {
			reAttributes := regexp.MustCompile(`<` + t + `[^>]*>`)
			content = reAttributes.ReplaceAllString(content, `<`+t+`>`)
		}
	}
	return content
}

func clearSpaces(content string) string {

	content = reClearAllowTags.ReplaceAllString(content, "\n")
	content = reEnter.ReplaceAllString(content, "\n")
	content = reEnters.ReplaceAllString(content, "\n")
	//for {
	//	if reEmptyDIV.MatchString(content){
	//		content = reEmptyDIV.ReplaceAllString(content, "")
	//	} else {
	//		break
	//	}
	//}
	content = basis.ReSpaces.ReplaceAllString(content, " ")
	return content
}
