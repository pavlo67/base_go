package files_sqlite

import (
	"database/sql"
	"time"

	"github.com/pavlo67/base_go/entities/files"
)

type itemScanner interface {
	Scan(dest ...any) error
}

func scanItem(scanner itemScanner) (*files.Item, error) {
	var item files.Item
	var isDir int
	var cTime sql.NullString
	var mTime sql.NullString
	var crc sql.NullInt64
	var createdAt sql.NullString
	var updatedAt sql.NullString

	err := scanner.Scan(
		&item.Path,
		&isDir,
		&item.Size,
		&cTime,
		&mTime,
		&crc,
		&item.MimeType,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.IsDir = isDir != 0

	if cTime.Valid && cTime.String != "" {
		t, err := time.Parse(time.RFC3339Nano, cTime.String)
		if err != nil {
			return nil, err
		}
		item.CTime = t
	}

	if mTime.Valid && mTime.String != "" {
		t, err := time.Parse(time.RFC3339Nano, mTime.String)
		if err != nil {
			return nil, err
		}
		item.MTime = t
	}

	if crc.Valid {
		item.CRC = &crc.Int64
	}

	if createdAt.Valid && createdAt.String != "" {
		t, err := time.Parse(time.RFC3339Nano, createdAt.String)
		if err != nil {
			return nil, err
		}
		item.CreatedAt = t
	}

	if updatedAt.Valid && updatedAt.String != "" {
		t, err := time.Parse(time.RFC3339Nano, updatedAt.String)
		if err != nil {
			return nil, err
		}
		item.UpdatedAt = t
	}

	return &item, nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func timeToDB(t time.Time) any {
	//if t.IsZero() {
	//	return nil
	//}
	return t.UTC().Format(time.RFC3339Nano)
}

//func timePtrToDB(t *time.Time) any {
//	if t == nil || t.IsZero() {
//		return nil
//	}
//	return t.UTC().Format(time.RFC3339Nano)
//}
