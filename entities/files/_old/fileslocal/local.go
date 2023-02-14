package fileslocal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"

	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/basis/filelib/inspector"

	"github.com/pavlo67/punctum/confidenter/groups"

	"github.com/pavlo67/punctum/things_old/fileinfo"
	"github.com/pavlo67/punctum/things_old/files"
)

var _ files.Operator = &filesLocal{}

type filesLocal struct {
	repositoryPath string

	fileinfoOp fileinfo.Operator
}

func New(repositoryPath string, fileinfoOp fileinfo.Operator) (*filesLocal, error) {
	err := filelib.Dir(repositoryPath)
	if err != nil {
		return nil, err
	}

	if fileinfoOp == nil {
		return nil, errors.New("no fileinfo.Operator for fileslocal")
	}

	f := filesLocal{
		repositoryPath: repositoryPath,

		fileinfoOp: fileinfoOp,
	}

	return &f, nil
}

const maxFilenameTries = 50

func (fOp *filesLocal) createUserPath(userIS auth.ID, filename string) (string, string, error) {

	if len(filename) >= len(files.RepoSchema) && filename[:len(files.RepoSchema)] == files.RepoSchema {
		filename = filename[len(files.RepoSchema):]
	}

	ext := filepath.Ext(filename)
	if len(ext) <= len(filename) && ext != "" {
		filename = filename[:len(filename)-len(ext)]
	}

	pathPrefix := fOp.repositoryPath + userIS.Identity().Domain + "_"

	var fullFilename string
	for i := 0; i < maxFilenameTries; i++ {
		if i > 0 {
			fullFilename = userIS.Identity().Path + "_" + userIS.Identity().ID + "/" + filename + "." + strconv.Itoa(i) + ext
		} else {
			fullFilename = userIS.Identity().Path + "_" + userIS.Identity().ID + "/" + filename + ext
		}
		if _, err := os.Stat(pathPrefix + fullFilename); err != nil {
			return files.RepoSchema + fullFilename, pathPrefix + fullFilename, nil
		}
	}

	return "", "", errors.New("file already exists")
}

func (fOp *filesLocal) checkUserDir(userIS auth.ID) error {
	if userIS == "" {
		return basis.ErrBadIdentity
	}

	return filelib.Dir(fOp.repositoryPath + userIS.Identity().Domain + "_" + userIS.Identity().Path + "_" + userIS.Identity().ID)
}

const onCreate = "on fileslocal.Create()"

func (fOp *filesLocal) Create(userIS auth.ID, file *files.Item) (string, error) {
	if file == nil {
		return "", errors.Wrap(basis.ErrNull, onCreate)
	}
	file.Name = filelib.CorrectFileName(file.Name)
	if file.Name == "" {
		return "", fmt.Errorf(onCreate+": can't create file - empty .Label: %#v", file)
	}

	err := fOp.checkUserDir(userIS)
	if err != nil {
		return "", errors.Wrap(err, onCreate)
	}

	var localPath string

	file.Name, localPath, err = fOp.createUserPath(userIS, file.Name)
	if err != nil {
		return "", errors.Wrap(err, onCreate)
	}

	err = ioutil.WriteFile(localPath, file.Content, 0644)
	if err != nil {
		return "", errors.Wrapf(err, onCreate+": can't create file: %s", localPath)
	}

	if file.MIMEType == "" {
		file.MIMEType, err = inspector.MIME(localPath, nil)
		if err != nil {
			log.Println(onCreate+": can't read MIMEType for file: ", localPath, err)
		}
	}

	fi := &fileinfo.Info{
		Data:      file.Data,
		LocalName: localPath,
	}

	filenameVirtual, err := fOp.fileinfoOp.Create(userIS, fi)
	if err != nil {
		os.Remove(localPath)
		return "", errors.Wrapf(err, onCreate+": can't create virtual file name: %#v", file)
	}

	// TODO: create preview

	return filenameVirtual, nil
}

const onRead = "on filesLocal.Read()"

func (fOp *filesLocal) Read(userIS auth.ID, filename string) (*files.Item, error) {
	fi, _, err := fOp.fileinfoOp.Read(userIS, filename)
	if err != nil {
		return nil, errors.Wrapf(err, onRead+": %s", filename)
	}

	file, err := fOp.readFile(fi)
	if err != nil {
		return nil, errors.Wrapf(err, onRead+": %s", filename)
	}

	return file, nil
}

func (fOp *filesLocal) readFile(fi *fileinfo.Info) (*files.Item, error) {
	if fi == nil {
		return nil, errors.New("no fileinfo")
	}

	localFilename := fi.LocalPath + fi.LocalName

	stat, err := os.Stat(localFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "can't stat file: %s", localFilename)
	}

	file := &files.Item{
		Data:     fi.Data,
		Size:     stat.Size(),
		IsDir:    stat.IsDir(),
		Modified: stat.ModTime(),
	}

	if file.IsDir {
		return file, nil
	}

	openFile, err := os.Open(localFilename)
	defer openFile.Close()

	file.Content, err = ioutil.ReadAll(openFile)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read file: %s", localFilename)
	}

	switch fi.MIMEType {
	case "text/html; charset=windows-1251":
		contentUTF8, err := inspector.Win1251ToUTF8(file.Content)
		if err == nil {
			file.Content = inspector.ChangeCharsetWin1251(contentUTF8)
			file.MIMEType = "text/html; charset=utf-8"
		} else {
			log.Println("can't convert Win1251ToUTF8() for file:", file.Name, err)
		}
	}

	return file, nil
}

const onReadList = "on fileslocal.ReadList()"

func (fOp *filesLocal) ReadList(userIS auth.ID, options *content.ListOptions) ([]files.Item, uint64, error) {
	fiAll, cntAll, err := fOp.fileinfoOp.ReadList(userIS, options)
	if err != nil {
		return nil, 0, errors.Wrapf(err, onReadList)
	}

	var filesAll []files.Item
	var errs basis.Errors

	for _, fi := range fiAll {
		file, err := fOp.readFile(&fi)
		if err != nil {
			errs = append(errs, err)
		}
		if file != nil {
			filesAll = append(filesAll, *file)
		}

	}
	return filesAll, cntAll, errs.Err()
}

const onUpdate = "on fileslocal.Update()"

func (fOp *filesLocal) Update(userIS auth.ID, file *files.Item) (crud.Result, error) {
	if file == nil {
		return crud.Result{}, basis.ErrNull
	}
	content := file.Content

	fi, ctrlOp, err := fOp.fileinfoOp.Read(userIS, file.Name)
	if err != nil {
		return crud.Result{}, err
	}

	// TODO: check fi.Managers also
	if !groups.OneOf(userIS, ctrlOp, fi.ROwner) {
		return crud.Result{}, errors.New("no rights to change non-own files")
	}

	localFilename := fi.LocalPath + fi.LocalName

	err = ioutil.WriteFile(localFilename, content, 0644)
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, "can't update file: ", localFilename)
	}

	if file.MIMEType == "" {
		file.MIMEType, err = inspector.MIME(localFilename, nil)
		if err != nil {
			log.Println(onCreate+": can't read MIMEType for file: ", localFilename, err)
		}
	}

	fi.Data = file.Data
	_, err = fOp.fileinfoOp.Update(userIS, fi)
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, "can't update file file: ", fi.LocalPath)
	}

	return crud.Result{1}, nil
}

func (fOp *filesLocal) Delete(userIS auth.ID, filename string) (crud.Result, error) {

	if filename == "" {
		return crud.Result{}, errors.New("error delete filelib, file name is empty")
	}

	fi, ctrlOp, err := fOp.fileinfoOp.Read(userIS, filename)
	if err != nil {
		return crud.Result{}, err
	}

	// TODO: check fi.Managers also
	if !groups.OneOf(userIS, ctrlOp, fi.ROwner) {
		return crud.Result{}, errors.New("no rights to delete non-own files")
	}

	localFilename := fi.LocalPath + fi.LocalName

	err = os.Remove(localFilename)
	if err != nil {
		return crud.Result{}, err
	}
	_, err = fOp.fileinfoOp.Delete(userIS, filename)
	if err != nil {
		log.Println(err, "can't delete file_info for filelib", filename)
	}

	// TODO: delete preview
	return crud.Result{1}, nil
}

func (fOp *filesLocal) Close() {
}

//fileExt := strings.ToLower(filepath.Ext(filename))
//if str_json.ReImageExt.MatchString(fileExt) {
//	previewPath := str_json.ReFileExt.ReplaceAllString(info.LocalPath, pxPreviewStr+"px"+"${1}")
//	err = os.Remove(previewPath)
//	if err != nil {
//		//os.Exit(0)
//		log.Println(err, "can't delete preview file: ", previewPath)
//	} else {
//		previewName := str_json.ReFileExt.ReplaceAllString(filename, pxPreviewStr+"px"+"${1}")
//		_, err = fOp.fileinfoOp.DeleteList(userIS, previewName)
//		if err != nil {
//			log.Println(err, "can't delete file_info for previe filelib", previewName)
//		}
//	}
//
//}

//pathPreview, err := makePreviewFile(file.LocalPath)
//if err == nil && pathPreview != "" {
//	_, err := fOp.fileinfoOp.Create(
//		userIS,
//		&fileinfo.Info{
//			LocalPath: pathPreview,
//			Item: items.Item{
//				Label:   str_json.ReFileExt.ReplaceAllString(file.Label, pxPreviewStr+"px"+"${1}"),
//				RView:  file.RView,
//				ROwner: file.ROwner,
//			},
//		})
//	if err != nil {
//		return "", nil.Wrapf(err, "can't create virtual file name for file: %v", pathPreview)
//	}
//} else if err != nil {
//	log.Println("can't make previe for file: ", file.LocalPath, err)
//}
