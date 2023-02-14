package files

//var _ crud.Operator = &OperatorCRUD{}
//
//type OperatorCRUD struct {
//	Operator
//}
//
//func (op OperatorCRUD) Describe() (crud.Description, error) {
//	return crud.Description{
//		Fields: []crud.Field{
//			{Key: "content", Creatable: true, Editable: true},
//
//			{Key: "name", NotEmpty: true, Unique: true, AutoUnique: true},
//
//			{Key: "mimetype", Creatable: true, Editable: true},
//			{Key: "links", Creatable: true, Editable: true},
//
//			{Key: "r_view", Creatable: true, Editable: true, NotEmpty: true},
//			{Key: "r_owner", Creatable: true, Editable: true, NotEmpty: true},
//			{Key: "managers", Creatable: true, Editable: true},
//
//			{Key: "global_is", Creatable: true},
//		},
//	}, nil
//}
//
//func FileFromData(data crud.Contentus) (*items.Item, error) {
//	if data == nil {
//		return nil, basis.ErrNullItem
//	}
//
//	var err error
//
//	var linksList []items.Link
//	if data["links"] != "" {
//		err = json.Unmarshal([]byte(data["links"]), &linksList)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Unmarshal([]byte(data["links"]: %s)`, data["linksList"])
//		}
//	}
//
//	var managers rights.Managers
//	if data["managers"] != "" {
//		err = json.Unmarshal([]byte(data["managers"]), &managers)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Unmarshal([]byte(data["managers"]: %s)`, data["managers"])
//		}
//	}
//
//	return &items.Item{
//		Contentus: items.Contentus{
//			Name:     data["name"],
//			MIMEType: data["mimetype"],
//			RView:    basis.UserIS(data["r_view"]),
//			ROwner:   basis.UserIS(data["r_owner"]),
//			GlobalIS: data["global_is"],
//			Links:    linksList,
//			Managers: managers,
//		},
//
//		Contentus: []byte(data["content"]),
//	}, nil
//}
//
//func DataFromFile(file *items.Item) (crud.Contentus, error) {
//	if file == nil {
//		return nil, basis.ErrNullItem
//	}
//
//	var err error
//
//	var jsonLinks []byte
//	if file.Links != nil {
//		jsonLinks, err = json.Marshal(file.Links)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Marshal(obj.Links): %#v)`, file.Links)
//		}
//	}
//
//	var jsonManagers []byte
//	if len(file.Managers) > 0 {
//		jsonManagers, err = json.Marshal(file.Managers)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Marshal(obj.Managers): %#v)`, file.Managers)
//		}
//	}
//
//	return crud.Contentus{
//		"content": string(file.Contentus),
//
//		"name":     file.Name,
//		"mimetype": file.MIMEType,
//		"links":    string(jsonLinks),
//
//		"r_view":   string(file.RView),
//		"r_owner":  string(file.ROwner),
//		"managers": string(jsonManagers),
//
//		"global_is": file.GlobalIS,
//	}, nil
//}
//
//func (op OperatorCRUD) Create(userIS basis.UserIS, data crud.Contentus) (id string, err error) {
//	file, err := FileFromData(data)
//	if err != nil {
//		return "", err
//	}
//
//	return op.Operator.Create(userIS, file)
//}
//
//func (op OperatorCRUD) Read(userIS basis.UserIS, id string) (crud.Contentus, error) {
//	file, err := op.Operator.Read(userIS, id)
//	if err != nil {
//		return nil, err
//	}
//
//	return DataFromFile(file)
//}
//
//func (op OperatorCRUD) ReadList(userIS basis.UserIS, options *content.ListOptions) ([]crud.Contentus, uint64, error) {
//	filesList, allCnt, err := op.Operator.ReadList(userIS, options)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	var dataList []crud.Contentus
//	for _, fi := range filesList {
//		data, err := DataFromFile(&fi)
//		if err != nil {
//			return dataList, allCnt, err
//		}
//		dataList = append(dataList, data)
//	}
//
//	return dataList, allCnt, nil
//}
//
//func (op OperatorCRUD) Update(userIS basis.UserIS, data crud.Contentus) (crud.Result, error) {
//	fi, err := FileFromData(data)
//	if err != nil {
//		return crud.Result{}, err
//	}
//	return op.Operator.Update(userIS, fi)
//}
//
//func (op OperatorCRUD) TestCases(cleaner crud.Cleaner) ([]crud.OperatorTestCase, error) {
//
//	is := basis.UserIS("a/b/c")
//	isAnother := basis.UserIS("d/e/f")
//	userISNil := basis.UserIS("")
//
//	name := "Nick One"
//
//	repo := "repo:"
//
//	linksToCreatePrivate, _ := json.Marshal([]items.Link{
//		{TargetID: "1", Type: LinkType, Name: "Nick file test", ROwner: is, RView: is},
//		{TargetID: "Name2 file2 test2", Type: "tag", Name: "Name2 file2 test2", ROwner: is, RView: is},
//	})
//
//	toCreatePrivate := crud.Contentus{
//		"content": "1234",
//
//		"name":     repo + name + ".0.txt",
//		"mimetype": "text/plain",
//		"links":    string(linksToCreatePrivate),
//
//		"r_view":  string(is),
//		"r_owner": string(is),
//
//		"global_is": "3456reyt",
//	}
//
//	toUpdatePrivate := crud.Contentus{
//		"content": "56",
//
//		"name":     repo + name + ".1.txt",
//		"mimetype": "text/html",
//
//		"r_view":  string(is),
//		"r_owner": string(is),
//	}
//
//	linksToCreatePublic, _ := json.Marshal([]items.Link{
//		{TargetID: "1", Type: LinkType, Name: "Nick file test", ROwner: is, RView: basis.Anyone},
//		{TargetID: "Name2 file2 test2", Type: "tag", Name: "Name2 file2 test2", ROwner: is, RView: basis.Anyone},
//	})
//
//	toCreatePublic := crud.Contentus{
//		"content": "78",
//
//		"name":     repo + name + ".2.txt",
//		"mimetype": "text/html",
//		"links":    string(linksToCreatePublic),
//
//		"r_view":  string(basis.Anyone),
//		"r_owner": string(is),
//
//		"global_is": "rwereyt",
//	}
//
//	toUpdatePublic := crud.Contentus{
//		"content": "90",
//
//		"name":     repo + name + ".3.txt",
//		"mimetype": "text/plain",
//
//		"r_view":  string(basis.Anyone),
//		"r_owner": string(is),
//	}
//
//	testCases := []crud.OperatorTestCase{
//
//		// 0. all ok for private record,
//		// can't create with identityNil,
//		// can't read, update or delete with identityAnother
//		{
//			Operator: op,
//			Cleaner:  cleaner,
//
//			KeyField: "name",
//
//			ISToCreate:        is,
//			ISToCreateBad:     &userISNil,
//			ToSave:          toCreatePrivate,
//			ExpectedSaveErr: nil,
//
//			ISToRead:        is,
//			ISToReadBad:     &isAnother,
//			ExpectedReadErr: nil,
//
//			ISToUpdate:        is,
//			ISToUpdateBad:     &isAnother,
//			ToUpdate:          toUpdatePrivate,
//			ExpectedUpdateErr: nil,
//
//			ISToDelete:        is,
//			ISToDeleteBad:     &isAnother,
//			ExpectedRemoveErr: nil,
//		},
//
//		// 1. all ok for private record,
//		// can't create with identityNil,
//		// can't read, update or delete with identityNil
//		{
//			Operator: op,
//			Cleaner:  cleaner,
//
//			KeyField: "name",
//
//			ISToCreate:        is,
//			ISToCreateBad:     &userISNil,
//			ToSave:          toCreatePrivate,
//			ExpectedSaveErr: nil,
//
//			ISToRead:        is,
//			ISToReadBad:     &userISNil,
//			ExpectedReadErr: nil,
//
//			ISToUpdate:        is,
//			ISToUpdateBad:     &userISNil,
//			ToUpdate:          toUpdatePrivate,
//			ExpectedUpdateErr: nil,
//
//			ISToDelete:        is,
//			ISToDeleteBad:     &userISNil,
//			ExpectedRemoveErr: nil,
//		},
//
//		// 2. all ok for public record,
//		// can't create with identityNil,
//		// can read with identityAnother
//		// can't update or delete with identityAnother
//		{
//			Operator: op,
//			Cleaner:  cleaner,
//
//			KeyField: "name",
//
//			ISToCreate:        is,
//			ISToCreateBad:     &userISNil,
//			ToSave:          toCreatePublic,
//			ExpectedSaveErr: nil,
//
//			ISToRead:        isAnother,
//			ExpectedReadErr: nil,
//
//			ISToUpdate:        is,
//			ISToUpdateBad:     &isAnother,
//			ToUpdate:          toUpdatePublic,
//			ExpectedUpdateErr: nil,
//
//			ISToDelete:        is,
//			ISToDeleteBad:     &isAnother,
//			ExpectedRemoveErr: nil,
//		},
//
//		// 3. all ok for public record,
//		// can't create with identityNil,
//		// can read with identityNil,
//		// can't update or delete with identityNil
//		// close database
//		{
//			Operator: op,
//			Cleaner:  cleaner,
//
//			KeyField: "name",
//
//			ISToCreate:        is,
//			ISToCreateBad:     &userISNil,
//			ToSave:          toCreatePublic,
//			ExpectedSaveErr: nil,
//
//			ISToRead:        userISNil,
//			ExpectedReadErr: nil,
//
//			ISToUpdate:        is,
//			ISToUpdateBad:     &userISNil,
//			ToUpdate:          toUpdatePublic,
//			ExpectedUpdateErr: nil,
//
//			ISToDelete:        is,
//			ISToDeleteBad:     &userISNil,
//			ExpectedRemoveErr: nil,
//		},
//	}
//
//	return testCases, nil
//}
