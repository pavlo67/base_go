package main

//func main() {
//	data0 := crud.Data{
//		Key: crud.Key{
//			Type: "type1111",
//			ID:   crud.NewID("id111111"),
//		},
//		Value: persons01.Item{
//			Person01: entities.Person01{
//				Firstnames: []string{"wqer", "eeeee"},
//				Middlename: "qwe",
//				Lastname:   "eeeeeeee",
//			},
//		},
//	}
//
//	dataJSON, err := json.Marshal(data0)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	data := crud.Data{
//		Value: json.RawMessage{},
//	}
//	if err := json.Unmarshal(dataJSON, &data); err != nil {
//		log.Fatal(err)
//	}
//
//	log.Printf("111111111111: %#v", data)
//
//	type DataRaw struct {
//		crud.Key
//		Value json.RawMessage
//	}
//
//	var dataRaw DataRaw
//	if err := json.Unmarshal(dataJSON, &dataRaw); err != nil {
//		log.Fatal(err)
//	}
//
//	log.Printf("111111111111: %#v", dataRaw)
//
//}
