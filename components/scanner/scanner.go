package scanner

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/data/components/importer"

	"github.com/pavlo67/data/entities/records"
	"github.com/pavlo67/data/entities/sources"
)

type Stat struct {
	Start     time.Time
	FountsNum int
}

var reMYSQL = regexp.MustCompile(`^mysql://(\w+)\?table=(\w+)`)

func GetDataFromFounts(fountOp sources.Operator, importerOp importer.Operator) Stat {

	fmt.Println("Scan fountOp ...")

	sources, err := fountOp.List()
	if err != nil {
		// log.Printf("can't exec fountOp.ReadAll(%v, nil, %v): %s", program.Identity(), sel, err)
	}
	log.Println("IS are sources: ", len(sources))
	scannerStat := Stat{Start: time.Now()}

	lastURL := ""
	for i, fount := range sources {
		if lastURL == fount.URL {
			continue
		}

		log.Print(fount.URL)

		lastURL = fount.URL
		scannerStat.FountsNum++

		// fountStat := Stat{FountID: fount.ID, Start: time.Now(), ScannerStart: scannerStat.Start}

		// TODO!!!
		//for j := i + 1; j < len(sources); j++ {
		//	if sources[j].URL == fount.URL && sources[j].ROwner != fount.ROwner {
		//		identityStrings[sources[j].ROwner.Identity()] = strconv.FormatInt(sources[j].ID, 10)
		//	} else {
		//		break
		//	}
		//}

		key := string(fount.ImportType)
		importOp, ok := program.GetInterfaceBySignature((*importer.Operator)(nil), key).(importer.Operator)
		if !ok {
			err = errors.New("no " + key + " interface found for scanner")
			log.Println(err)
			fountStat.ResponseError = err.Error()
			addFountStat(fountOp, fountStat)
			continue
		}

		// scanned data save to importer
		//err = importOp.Init(fount.URL, fount.ImportDetailsParams, false)
		//if err != nil {
		//	scannerStat.ErrorsNum++
		//	log.Printf("error init fount importer (%s): %s", fount.URL, err)
		//
		//	fountStat.ResponseError = err.Error()
		//	addFountStat(fountOp, fountStat)
		//	continue
		//}

		var importOp importer.Operator
		var recordsOp records.Operator

		dataSeries, err := importOp.Get()

		for _, data := range dataSeries {

			fountStat.ItemsTaken++

			//item.FountIS = confidenter.IdentityString(domain + "/fount/" + fountID)
			//// TODO: use fount.Identity()
			//
			//item.RView = rOwner.String() // TODO!!!
			////item.RView = fount.RView
			//item.ROwner = rOwner.String()

			isNew, err := importerOp.IsNew(data)
			if err != nil {
				log.Printf("error checking importerOp.IsNew(%v, %d, %s): %s", rOwner, fount.ID, item.OriginalID, err)
				fountStat.LastItemError = err.Error()
			}

			runes := []rune(item.Summary)
			if len(runes) > 255 {
				item.Summary = string(runes[0:255])
			}
			_, err = importerOp.Create(&rOwner, *item)
			if err != nil {
				log.Printf("error importerOp.Create(%s, %v): %s", rOwner, item, err)
			}

			if isNew {
				fountStat.ItemsNew++
				_, err := recordsOp.Save(data, nil)
			}

			// if fount.ToFlow {
			//} else if fount.ToObject {
			//	mySQLData := reMYSQL.FindStringSubmatch(fount.URL)
			//	if len(mySQLData) > 0 {
			//		err = importOp.Init(mySQLData[2], mySQLData[1], false)
			//	} else {
			//		err = errors.New("can't parse url")
			//	}
			//	if err != nil {
			//		scannerStat.ErrorsNum++
			//		log.Printf("error init fount importer (%s): %s", fount.URL, err)
			//
			//		fountStat.ResponseError = err.Err()
			//		addFountStat(fountOp, fountStat)
			//		continue
			//	}
			//	for {
			//		entity, err := importOp.Next()
			//		if err == importer.ErrNoMoreItems {
			//			break
			//		}
			//		if err != nil {
			//			fountStat.ItemErrors++
			//			log.Printf("error reading imported item: %s", err)
			//			fountStat.LastItemError = err.Err()
			//			continue
			//		}
			//		fountStat.ItemsTaken++
			//		o, err := entity.Object()
			//		if err != nil {
			//			fountStat.LastItemError = err.Err()
			//			log.Printf("can't get Object() for imported row: %v; %v", entity, err)
			//			continue
			//		}
			//		o.RView = fount.RView
			//		fID := fount.ROwner.Identity()
			//		_, err = objectsOp.Create(&fID, o)
			//		if err != nil {
			//			me, ok := errors.Cause(err).(*mysql.MySQLError)
			//			if !ok || me.Number != 1062 {
			//				fountStat.LastItemError = err.Err()
			//				log.Printf("can't create import object for imported object: %v; %v", o, err)
			//			}
			//
			//			continue
			//		}
			//		fountStat.ItemsNew++
			//	}
		}

		addFountStat(fountOp, fountStat)
		if fountStat.LastItemError != "" {
			scannerStat.ErrorsNum++
		}
		scannerStat.ItemsTaken += fountStat.ItemsTaken
		scannerStat.ItemsNew += fountStat.ItemsNew
	}

	addScannerStat(fountOp, scannerStat)
	return scannerStat
}

func addScannerStat(fountOp sources.Operator, scannerStat Stat) {
	scannerStat.Duration = time.Now().UnixNano() - scannerStat.Start.UnixNano()
	identity := program.Identity()
	err := fountOp.AddScannerStat(&identity, scannerStat)
	if err != nil {
		log.Println(errors.Wrapf(err, "error write fount_stat: %v", scannerStat))
	}
}

func addFountStat(fountOp sources.Operator, fountStat sources.FountStat) {
	fountStat.Duration = time.Now().UnixNano() - fountStat.Start.UnixNano()
	identity := program.Identity()
	err := fountOp.AddFountStat(&identity, fountStat)
	if err != nil {
		log.Println(errors.Wrapf(err, "error write fount_stat: %v", fountStat))
	}
}
