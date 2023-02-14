package scanner

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/joiner"

	"github.com/pavlo67/data/entities/sources"

	"github.com/pavlo67/data/components/importer"
)

type Stat struct {
	Start           time.Time
	Duration        time.Duration
	SourcesNum      int
	RecordsTotalNum int
	RecordsSavedNum int
}

const onDataFromSources = "on scanner.DataFromSources()"

func DataFromSources(actor auth.Actor) (*Stat, error) {

	sourcesList, err := sourcesOp.List(nil, actor)
	if err != nil {
		return nil, errors.Wrap(err, onDataFromSources)
	}
	scannerStat := Stat{Start: time.Now(), SourcesNum: len(sourcesList)}

	importerOps := map[joiner.InterfaceKey]importer.Operator{}

	for i, source := range sourcesList {

		l.Infof("scanning #%d of %d sourcesList: %s", i+1, len(sourcesList), source.SourceURN)

		sourceStat := sources.Stat{Start: time.Now()}

		importerOp := importerOps[source.ImporterInterfaceKey]
		if importerOp == nil {
			if importerOp, _ = joinerOp.Interface(source.ImporterInterfaceKey).(importer.Operator); importerOp == nil {
				sourceStat.LastError = fmt.Errorf("interface '%s' isn't found", source.ImporterInterfaceKey)
				sourceStat.ErrorsNum++
				l.Error(sourceStat.LastError)

				addSourceStat(sourcesOp, sourceStat, actor)
				continue
			}
		}

		// scanned record save to importer
		// err = importOp.Init(source.URL, source.ImportDetailsParams, false)
		// if err != nil {
		//	scannerStat.ErrorsNum++
		//	log.Printf("error init source importer (%s): %s", source.URL, err)
		//
		//	sourceStat.ResponseError = err.Error()
		//	addSourceStat(sourcesOp, sourceStat)
		//	continue
		// }

		dataSeries, err := importerOp.Get(source.SourceURN, source.Params)
		if err != nil || dataSeries == nil {
			sourceStat.LastError = fmt.Errorf("got %#v / %s", dataSeries, err)
			sourceStat.ErrorsNum++
			l.Error(sourceStat.LastError)

			addSourceStat(sourcesOp, sourceStat, actor)
			continue
		}

		sourceStat.RecordsTotalNum = len(dataSeries.Records)

		for i, record := range dataSeries.Records {
			if i%100 == 0 {
				l.Infof("processing %d of %d records", i+1, dataSeries.Records)
			}

			isNew, err := importerOp.IsNew(record)
			if err != nil {
				sourceStat.LastError = fmt.Errorf("checking %d, %s got: %s", source.ID, record.SourceURN, err)
				sourceStat.ErrorsNum++
				l.Error(sourceStat.LastError)
			}

			if isNew {
				if _, _, err = recordsOp.Add(record, actor); err != nil {
					sourceStat.LastError = err
					sourceStat.ErrorsNum++
					l.Error(err)
				} else {
					sourceStat.RecordsSavedNum++
				}
			}

		}

		addSourceStat(sourcesOp, sourceStat, actor)

		scannerStat.RecordsTotalNum += sourceStat.RecordsTotalNum
		scannerStat.RecordsSavedNum += sourceStat.RecordsSavedNum
	}

	addScannerStat(scannerStat)

	return &scannerStat, nil
}

func addScannerStat(stat Stat) {
	stat.Duration = time.Now().Sub(stat.Start)
	l.Infof("TOTAL SCANNER STATISTICS: %#v", stat)

	//if err := sourceOp.AddScannerStat(stat, actor); err != nil {
	//	l.Error(err)
	//}
}

func addSourceStat(sourceOp sources.Operator, sourceStat sources.Stat, actor auth.Actor) {
	sourceStat.Duration = time.Now().Sub(sourceStat.Start)
	if err := sourceOp.AddStat(sourceStat, actor); err != nil {
		l.Error(err)
	}
}
