package sqllib_sqlite

//var reIndexLimit = regexp.MustCompile(`\(\d+\)`)

//func (_ SQLite) CreateSQL(table config.SQLTable) (string, error) {
//
//	//reText := regexp.MustCompile("text")
//	//reTime := regexp.MustCompile("(datetime|timestamp)")
//
//	sqlQuery := "create table `" + table.Name + "` ( "
//	firsFieldAdded := false
//
//	for _, f := range table.Fields {
//		if firsFieldAdded {
//			sqlQuery += ", \n"
//		}
//		sqlQuery += "`" + f.Name + "` " + f.Type + " "
//		if !f.Null {
//			sqlQuery += " NOT NULL "
//		} else {
//			sqlQuery += " NULL "
//		}
//
//		if f.Default != "" {
//			sqlQuery += " DEFAULT '" + f.Default + "' "
//		}
//
//		if f.Extra != "" {
//			sqlQuery += f.Extra + " "
//		}
//		firsFieldAdded = true
//	}
//	for _, i := range table.Indexes {
//		sqlQuery += ", \n"
//		if strings.ToUpper(i.Type) == "PRIMARY" {
//			sqlQuery += "PRIMARY KEY  ("
//			liF := 0
//			for _, f := range i.Fields {
//				if liF > 0 {
//					sqlQuery += ", "
//				}
//				sqlQuery += "`" + f + "`"
//				liF++
//			}
//			sqlQuery += ")"
//		}
//		if strings.ToUpper(i.Type) == "UNIQUE" || i.Type == "" {
//			if i.Type == strings.ToUpper("UNIQUE") {
//				sqlQuery += "UNIQUE KEY `" + i.Name + "` ("
//			} else {
//				sqlQuery += "KEY `" + i.Name + "` USING BTREE ("
//			}
//			liF := 0
//			for _, f := range i.Fields {
//				if liF > 0 {
//					sqlQuery += ", "
//				}
//				if reIndexLimit.MatchString(f) {
//					sqlQuery += f
//				} else {
//					sqlQuery += "`" + f + "`"
//				}
//				liF++
//			}
//			sqlQuery += ")"
//		}
//	}
//	sqlQuery += ");"
//
//	return sqlQuery, nil
//}
//
//func (_ SQLite) TableExistsSQL(tableName string) (string, error) {
//	tableName = strings.TrimSpace(tableName)
//	if tableName == "" {
//		return "", errors.New("empty table name")
//	}
//
//	return "SELECT name FROM sqlite_master WHERE type ='table' AND name = '" + tableName + "'", nil
//
//	// AND name NOT LIKE 'sqlite_%'
//}
