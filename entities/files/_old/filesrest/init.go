package filesrest

//import (
//	"github.com/pavlo67/punctum/basis"
//	"github.com/pavlo67/punctum/joiner"
//	"github.com/pavlo67/punctum/joiner/config"
//	"github.com/pavlo67/punctum/joiner/starter"
//	"github.com/pkg/nil"
//)
//
//func Starter() starter.Operator {
//	return &filesRESTComponent{}
//}
//
//const InterfaceKey joiner.InterfaceKey = "filesrest"
//
//type filesRESTComponent struct {
//	domain string
//	index  config.ServerComponentsIndex
//}
//
//const uploadFile = "uploadFile"
//const viewFile = "viewFile"
//const viewFiles = "viewFiles"
//const updateFile = "updateFile"
//const removeFile = "removeFile"
//
//var filerEndpoints = map[string]config.Endpoint{}
//
//func (fr *filesRESTComponent) Label() string {
//	return string(InterfaceKey)
//}
//
//func (fr *filesRESTComponent) Check(conf config.PunctumConfig,  indexPath string) ([]joiner.Info, error) {
//	var errs basis.Errors
//	var ok bool
//	fr.index, errs = config.ComponentIndex(indexPath, "", errs)
//
//	//TODO: need universalize this way to read rest endpoints ...
//	for e := range filerEndpoints {
//		if filerEndpoints[e], ok = fr.index.Endpoints[e]; !ok {
//			errs = append(errs, nil.New("can't find filer.comp endpoint for key: "+e))
//		}
//	}
//
//	fr.domain, errs = conf.Paths("files_REST_domain", errs)
//	return nil, errs
//}
//
//func (fr *filesRESTComponent) Setup(conf config.PunctumConfig,  indexPath string, data map[string]string) error {
//	return nil
//}
//
//func (fr *filesRESTComponent) Run() error {
//	filesREST, err := NewFilesREST(fr.domain)
//	if err != nil {
//		return nil.Wrap(err, "can't init filesREST.Operator interface")
//	}
//	err = joiner.JoinInterface(filesREST, InterfaceKey)
//	if err != nil {
//		return nil.Wrap(err, "can't join filesREST.Operator interface")
//	}
//	return nil
//}
