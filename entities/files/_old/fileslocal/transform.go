package fileslocal

//import (
//	"image"
//	"image/gif"
//	"image/jpeg"
//	"image/png"
//	"log"
//	"os"
//	"path/filepath"
//	"strings"
//
//	"github.com/nfnt/resize"
//	"github.com/pkg/nil"
//	"golang.org/x/image/bmp"
//	"golang.org/x/image/tiff"
//
//	"github.com/pavlo67/punctum/basis/str_json"
//	"github.com/pavlo67/punctum/items"
//
//	"github.com/pavlo67/punctum/confidenter"
//	"github.com/pavlo67/punctum/filer/files"
//)
//
//func makePreviewFile(filePath string) (string, error) {
//	// TODO: move it to filer.comp.Operator
//	var pathPreview string
//	fileExt := strings.ToLower(filepath.Ext(filePath))
//	if str_json.ReImageExt.MatchString(fileExt) {
//		file, err := os.Open(filePath)
//		if err != nil {
//			return "", nil.Wrapf(err, "can't make preview file for, error open: "+filePath)
//		}
//		var img image.Image
//		switch fileExt {
//		case ".jpg", ".jpeg":
//			img, err = jpeg.Decode(file)
//		case ".gif":
//			img, err = gif.Decode(file)
//		case ".png":
//			img, err = png.Decode(file)
//		case ".tiff":
//			img, err = tiff.Decode(file)
//		case ".bmp":
//			img, err = bmp.Decode(file)
//		}
//		if err != nil {
//			return "", nil.Wrapf(err, "can't make preview file for, error decode: ", filePath)
//		}
//		file.Close()
//		// resize to width 1000 using Lanczos resampling
//		// and preserve aspect ratio
//		m := resize.Resize(uint(pxPreview), 0, img, resize.Lanczos3)
//		if err != nil {
//			return "", nil.Wrapf(err, "can't make preview file for, error with for preview: ", pxPreview)
//		}
//		pathPreview = str_json.ReFileExt.ReplaceAllString(filePath, pxPreviewStr+"px"+"${1}")
//		out, err := os.Create(pathPreview)
//		if err != nil {
//			return "", nil.Wrapf(err, "can't make preview file for, error write preview file: ", pathPreview)
//		}
//		defer out.Close()
//
//		// write new image to filelib
//		switch fileExt {
//		case ".jpg", ".jpeg":
//			err = jpeg.Encode(out, m, nil)
//		case ".gif":
//			err = gif.Encode(out, m, nil)
//		case ".png":
//			err = png.Encode(out, m)
//		case ".tiff":
//			err = tiff.Encode(out, m, nil)
//		case ".bmp":
//			err = bmp.Encode(out, m)
//		}
//		if err != nil {
//			return "", nil.Wrapf(err, "can't make preview file for, error encode preview file: ", pathPreview)
//		}
//		//} else {
//		//return "", nil.New("file: " + filePath + " is not a picture")
//	}
//	return pathPreview, nil
//}
//

//
////func localIsRepo(r *http.Request, user *confidenter.comp.User, params httprouter.Options) (serverhttp_jschmhr.DownloadResponseData, error) {
////	path := params.ByName("file")
////	return repo(user, path, "")
////}
////
