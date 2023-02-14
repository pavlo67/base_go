package flowmysql

//// CRUD implementation -------------------------------------------------------------------------------------------------
//
////type flowCRUD struct {
////	flow
////}
//
//// Create ...
//func (f *FlowMySQL) Create(identity confidenter.Identity, description crud.Records) (*confidenter.Identity, error) {
//	err := f.DescriptionToData(description)
//	if err != nil {
//		return nil, err
//	}
//	res, err := f.Create(identity, *f.crudBuffer)
//	return &res, err
//}
//
//// Read returns object's Records data (accordingly to requester's rights).
//func (f *FlowMySQL) Read(identity confidenter.Identity, fragmentIS confidenter.Identity) (*crud.Records, error) {
//	var err error
//	f.crudBuffer, err = f.Read(identity, fragmentIS)
//	if err != nil {
//		return nil, err
//	}
//	return f.DataToDescription()
//}
//
//// ReadAllCRUD returns array of object's Records data (accordingly to requester's rights).
//func (f *FlowMySQL) ReadAllCRUD(identity confidenter.Identity, selector interfaces.Selector, options *crud.ReadAllOptions) ([]crud.Records, int64, error) {
//
//	data, allCnt, err := f.ReadAll(identity, options, selector)
//	if err != nil {
//		return nil, 0, err
//	}
//	var description []crud.Records
//	for _, c := range data {
//		f.crudBuffer = &c
//		desc, _ := f.DataToDescription()
//		description = append(description, *desc)
//	}
//	return description, allCnt, nil
//}
//
//// Update changes object's Records data (accordingly to requester's rights).
//func (f *FlowMySQL) Update(identity confidenter.Identity, fragmentIS confidenter.Identity, description crud.Records) (crud.Result, error) {
//	err := f.DescriptionToData(description)
//	if err != nil {
//		return crud.Result{}, err
//	}
//	return f.Update(identity, fragmentIS, *f.crudBuffer)
//}
//
//// Delete ...
//func (f *FlowMySQL) Delete(identity confidenter.Identity, fragmentIS confidenter.Identity) (crud.Result, error) {
//	return f.Delete(identity, fragmentIS)
//}
//
//// Count ...
//func (f *FlowMySQL) CountCRUD(selector interfaces.Selector, joinTo crud.JoinTo, groupBy, sortBy []string) ([]crud.Count, error) {
//
//	//if Tables[joinTo.ToTable] != "" {
//	//	joinTo.ToTable = Tables[joinTo.ToTable]
//	//} else {
//	//	return nil, errors.New("can't find table code: " + joinTo.ToTable)
//	//}
//
//	return clients.Count(f.dbh, selector, joinTo, groupBy, sortBy)
//}
//
//// Describe ... read crud.json5
//func (f *FlowMySQL) DescribeCRUD() (*crud.Description, error) {
//	return crud.Describe(basis.CurrentPath() + "../")
//}
//
//func (f *FlowMySQL) DataToDescription() (*crud.Records, error) {
//	i := f.crudBuffer
//	var err error
//	var data []byte
//	if i.Original != nil {
//		data, err = json.Marshal(i.Original)
//		if err != nil {
//			return nil, errors.Wrapf(err, "can't marshal data:%v in flow_crud", i.Original)
//		}
//	}
//	od := crud.Records{
//		Details: map[string]string{
//			"FountIS":    string(i.FountIS),
//			"FountURL":   i.FountURL,
//			"Original":   string(data),
//			"OriginalID": i.OriginalID,
//			"Summary":    i.Summary,
//			"Content":    i.Content,
//			"At":  i.At.Format("2006-01-02T15:04:05Z"),
//		},
//		Managers: controller.Managers{rights.View: i.RView, rights.Owner: i.ROwner},
//	}
//
//	return &od, nil
//}
//
////Records       interface{}
//
//func (f *FlowMySQL) DescriptionToData(o crud.Records) error {
//
//	var createdAt time.Time
//	var err error
//	var data interface{}
//	if o.Details["At"] != "" {
//		createdAt, err = time.Parse("2006-01-02T15:04:05Z", o.Details["At"])
//		if err != nil {
//			return errors.Wrapf(err, "can't parse time: %v in flow_crud", o.Details["At"])
//		}
//	}
//	if o.Details["Original"] != "" {
//		err = json.Unmarshal([]byte(o.Details["Original"]), &data)
//		if err != nil {
//			return errors.Wrapf(err, "can't Unmarshal : %v in flow_crud", o.Details["Original"])
//		}
//	}
//	i := flow.Item{
//		FountIS:    confidenter.IdentityString(o.Details["FountIS"]),
//		FountURL:   o.Details["FountURL"],
//		Original:   data,
//		OriginalID: o.Details["OriginalID"],
//		Summary:    o.Details["Summary"],
//		Content:    o.Details["Content"],
//		RView:      o.Managers[rights.View],
//		ROwner:     o.Managers[rights.Owner],
//		At:  createdAt,
//	}
//	f.crudBuffer = &i
//	return nil
//}
