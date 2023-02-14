package fountsmysql

//// CRUD implementation -------------------------------------------------------------------------------------------------
//
//// Create ...
//func (fms *MySQLFount) Create(identity confidenter.Identity, description crud.Records) (*confidenter.Identity, error) {
//	err := fms.DescriptionToData(description)
//	if err != nil {
//		return nil, err
//	}
//	ptrToIdentity, err := fms.Create(identity, *fms.crudBuffer)
//	if ptrToIdentity == nil {
//		return nil, err
//	}
//	return ptrToIdentity, err
//}
//
//// Read returns object's Records data (accordingly to requester's rights).
//func (fms *MySQLFount) Read(identity confidenter.Identity, fountIS confidenter.Identity) (*crud.Records, error) {
//	var err error
//	fms.crudBuffer, err = fms.Read(identity, fountIS.Identity().ID)
//	if err != nil {
//		return nil, err
//	}
//	return fms.DataToDescription()
//}
//
//// ReadAllCRUD returns array of object's Records data (accordingly to requester's rights).
//func (fms *MySQLFount) ReadAllCRUD(identity confidenter.Identity, selector interfaces.Selector, options *crud.ReadAllOptions) ([]crud.Records, int64, error) {
//
//	data, allCount, err := fms.ReadAll(identity, options, selector)
//	if err != nil {
//		return nil, 0, err
//	}
//	var description []crud.Records
//	for _, c := range data {
//		fms.crudBuffer = &c
//		desc, _ := fms.DataToDescription()
//		description = append(description, *desc)
//	}
//	return description, allCount, nil
//}
//
//// Update changes object's Records data (accordingly to requester's rights).
//func (fms *MySQLFount) Update(identity confidenter.Identity, fountIS confidenter.Identity, description crud.Records) (crud.Result, error) {
//	err := fms.DescriptionToData(description)
//	if err != nil {
//		return crud.Result{}, err
//	}
//	return fms.Update(identity, *fms.crudBuffer)
//}
//
//// Delete ...
//func (fms *MySQLFount) Delete(identity confidenter.Identity, fountIS confidenter.Identity) (crud.Result, error) {
//	return fms.Delete(identity, fountIS.Identity().ID)
//}
//
//// Count ...
//func (fms *MySQLFount) CountCRUD(selector interfaces.Selector, joinTo crud.JoinTo, groupBy, sortBy []string) ([]crud.Count, error) {
//	if Tables[joinTo.ToTable] != "" {
//		joinTo.ToTable = Tables[joinTo.ToTable]
//	} else {
//		return nil, errors.New("can't find table code: " + joinTo.ToTable)
//	}
//	return clients.Count(fms.dbh, selector, joinTo, groupBy, sortBy)
//}
//
//// Describe ... read crud.json5
//func (fms *MySQLFount) DescribeCRUD() (*crud.Description, error) {
//	return crud.Describe(basis.CurrentPath() + "../")
//}
//
//func (fms *MySQLFount) DescriptionToData(description crud.Records) error {
//	var createdAt, updatedAt time.Time
//	var err error
//	if description.Details["At"] != "" {
//		createdAt, err = time.Parse("2006-01-02T15:04:05Z", description.Details["At"])
//		if err != nil {
//			return errors.Wrapf(err, "can't parse time: %v in fount_crud", description.Details["At"])
//		}
//	}
//	if description.Details["UpdatedAt"] != "" {
//		updatedAt, err = time.Parse("2006-01-02T15:04:05Z", description.Details["UpdatedAt"])
//		if err != nil {
//			return errors.Wrapf(err, "can't parse time: %v in fount_crud", description.Details["UpdatedAt"])
//		}
//	}
//	id, err := strconv.ParseInt(description.ID, 10, 64)
//	if err != nil {
//		return errors.Wrapf(err, "can't parse int: %v in fount_crud", description.ID)
//	}
//	fms.crudBuffer = &founts.Fount{
//		ID:                  id,
//		URL:                 description.Details["Url"],
//		Label:               description.Details["Label"],
//		ImportType:          importer.ImportType(description.Details["ImportType"]),
//		ToFlow:              description.Details["Direct"] == "flow",
//		ToObject:            description.Details["Direct"] == "object",
//		ROwner:              description.Managers[rights.Owner],
//		RView:               description.Managers[rights.View],
//		ManagersRaw:         description.Details["ManagersRaw"],
//		At:           createdAt,
//		UpdatedAt:           &updatedAt,
//		Managers:            description.Managers,
//		ImportDetailsType:   description.Details["ImportDetailsType"],
//		ImportDetailsParams: description.Details["ImportDetailsParams"],
//	}
//	return nil
//}
//
//func (fms *MySQLFount) DataToDescription() (*crud.Records, error) {
//	direct := ""
//	if fms.crudBuffer.ToFlow {
//		direct = "flow"
//	} else if fms.crudBuffer.ToObject {
//		direct = "object"
//	}
//	m := fms.crudBuffer.Managers
//	m = controller.Managers{
//		rights.Owner: fms.crudBuffer.ROwner,
//		rights.View:  fms.crudBuffer.RView,
//	}
//	var updatedAt string
//	if fms.crudBuffer.UpdatedAt != nil {
//		updatedAt = fms.crudBuffer.UpdatedAt.Format("2006-01-02T15:04:05Z")
//	}
//	return &crud.Records{
//		ID: strconv.FormatInt(fms.crudBuffer.ID, 10),
//		Details: map[string]string{
//			"Url":         fms.crudBuffer.URL,
//			"Label":       fms.crudBuffer.Label,
//			"Direct":      direct,
//			"ImportType":  string(fms.crudBuffer.ImportType),
//			"ManagersRaw": fms.crudBuffer.ManagersRaw,
//			"At":   fms.crudBuffer.At.Format("2006-01-02T15:04:05Z"),
//			"UpdatedAt":   updatedAt,
//
//			"ImportDetailsParams": fms.crudBuffer.ImportDetailsParams,
//			"ImportDetailsType":   fms.crudBuffer.ImportDetailsType,
//		},
//		Managers: m,
//	}, nil
//}
