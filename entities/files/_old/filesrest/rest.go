package filesrest

//import (
//	"bytes"
//	"encoding/json"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"strconv"
//	"strings"
//
//	"github.com/pavlo67/punctum/confidenter"
//	"github.com/pavlo67/punctum/fronthttp/serverhttp_jschmhr"
//	"github.com/pavlo67/punctum/crud"
//	"github.com/pavlo67/punctum/filer"
//	"github.com/pavlo67/punctum/filer/fileinfo"
//	"github.com/pavlo67/punctum/filer/files"
//	"github.com/pkg/nil"
//)
//
//var _ files.Operator = &FilesREST{}
//
//type FilesREST struct {
//	DomainREST string
//}
//
//func NewFilesREST(domainREST string) (*FilesREST, error) {
//	return &FilesREST{DomainREST: domainREST}, nil
//}
//
//// Create ...
//func (fr *FilesREST) Create(userIS basis.UserIS, info *items.File) (string, error) {
//	req, err := http.NewRequest(
//		filerEndpoints[uploadFile].Method,
//		fr.DomainREST+filerEndpoints[uploadFile].WithParams+info.File.Label,
//		bytes.NewBuffer(info.Contentus))
//	if err != nil {
//		return "", err
//	}
//	req.Header.Set("Content-ImporterInterfaceKey", "application/octet-stream")
//	//TODO: need get token
//	token := "???? "
//	//req.Header.Set("Authorization", token)
//	req.Header.Set("Cookie", "token="+token+";")
//
//	var client = &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//	body, _ := ioutil.ReadAll(resp.Body)
//	//fmt.Println("response Body:", string(body))
//	if resp.Status != "200 OK" {
//		return "", nil.New(string(body))
//	}
//	var r serverhttp_jschmhr.RESTData
//	err = json.Unmarshal(body, &r)
//	if err != nil {
//		return "", nil.New("can't unmarshal rest uploadFile body: " + string(body))
//	}
//	return r.Info, nil
//}
//
//// Read ...
//func (fr *FilesREST) Read(userIS basis.UserIS, filename string) (*fileinfo.Info, error) {
//
//	req, err := http.NewRequest(
//		filerEndpoints[viewFile].Method, fr.DomainREST+filerEndpoints[viewFile].WithParams+filename,
//		nil)
//	if err != nil {
//		return nil, err
//	}
//
//	//TODO: need get token
//	token := "???? "
//	//req.Header.Set("Authorization", token)
//	req.Header.Set("Cookie", "token="+token+";")
//	var client = &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	body, _ := ioutil.ReadAll(resp.Body)
//	//fmt.Println("response Body:", string(body))
//	if resp.Status != "200 OK" {
//		return nil, nil.New(string(body))
//	}
//
//	size, err := strconv.ParseInt(resp.Header.ReadList("Contentus-TokenLength"), 10, 64)
//	nameFile := strings.Replace(resp.Header.ReadList("Contentus-Disposition"), "attachment; filename=", "", 1)
//
//	fi := fileinfo.Info{
//		File: items.File{
//			MIMEType: resp.Header.ReadList("Content-ImporterInterfaceKey"),
//			Contentus:  body,
//			Size:     size,
//			Label:     nameFile,
//		},
//	}
//
//	return &fi, nil
//}
//
//// ReadList ...
//func (fr *FilesREST) ReadList(userIS basis.UserIS, options *content.ListOptions) ([]fileinfo.Info, uint64, error) {
//	return nil, 0, nil
//}
//
//// Update ...
//func (fr *FilesREST) Update(userIS basis.UserIS, info *fileinfo.Info) (crud.Result, error) {
//	req, err := http.NewRequest(
//		filerEndpoints[updateFile].Method,
//		fr.DomainREST+filerEndpoints[updateFile].WithParams+info.File.Label,
//		bytes.NewBuffer(info.Contentus))
//	if err != nil {
//		return crud.Result{}, err
//	}
//	req.Header.Set("Content-ImporterInterfaceKey", "application/octet-stream")
//	//TODO: need get token
//	token := "???? "
//	//req.Header.Set("Authorization", token)
//	req.Header.Set("Cookie", "token="+token+";")
//
//	var client = &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return crud.Result{}, err
//	}
//	defer resp.Body.Close()
//	body, _ := ioutil.ReadAll(resp.Body)
//	//fmt.Println("response Body:", string(body))
//	if resp.Status != "200 OK" {
//		return crud.Result{}, nil.New(string(body))
//	}
//	var r serverhttp_jschmhr.RESTData
//	err = json.Unmarshal(body, &r)
//	if err != nil {
//		return crud.Result{}, nil.New("can't unmarshal rest updateFile body: " + string(body))
//	}
//	cnt, err := strconv.ParseInt(r.Info, 10, 64)
//	if err != nil {
//		log.Println("filesrest.Update() error; can't ParseInt: ", r.Info)
//	}
//	return crud.Result{NumOk: cnt}, nil
//}
//
//// DeleteList ...
//func (fr *FilesREST) DeleteList(userIS basis.UserIS, fileName string) (crud.Result, error) {
//	req, err := http.NewRequest(
//		filerEndpoints[removeFile].Method, fr.DomainREST+filerEndpoints[removeFile].WithParams+fileName,
//		nil)
//	if err != nil {
//		return crud.Result{}, err
//	}
//
//	//TODO: need get token
//	token := "???? "
//	//req.Header.Set("Authorization", token)
//	req.Header.Set("Cookie", "token="+token+";")
//	var client = &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return crud.Result{}, err
//	}
//	defer resp.Body.Close()
//	body, _ := ioutil.ReadAll(resp.Body)
//	var r serverhttp_jschmhr.RESTData
//	err = json.Unmarshal(body, &r)
//	if err != nil {
//		return crud.Result{}, nil.New("can't unmarshal rest updateFile body: " + string(body))
//	}
//	if r.Info == fileName {
//		return crud.Result{NumOk: 1}, nil
//	}
//	return crud.Result{}, nil
//}
//
//func (fr *FilesREST) Clean() error {
//	return nil
//}
//
//func (fr *FilesREST) Close() {
//}
