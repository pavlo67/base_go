package founts

import (
	"strconv"
	"time"

	"github.com/pavlo67/punctum/interfaces/importer"

	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces"
	"github.com/pavlo67/punctum/interfaces/confidenter"
	"github.com/pavlo67/punctum/interfaces/controller"
	"github.com/pavlo67/punctum/interfaces/crud"
)

const InterfaceKey = "founts"

type Fount struct {
	ID          int64
	URL         string
	Title       string
	ImportType  importer.ImportType
	ToFlow      bool
	ToObject    bool
	Tags        string
	ROwner      confidenter.IdentityString
	RView       confidenter.IdentityString
	ManagersRaw string
	Managers    controller.Managers

	ImportDetailsType   string
	ImportDetailsParams string

	CreatedAt time.Time
	UpdatedAt *time.Time
}

type FountTag struct {
	ID      int64
	FountID int64
	Tag     string
	RView   confidenter.Identity
}

func (f *Fount) IdentityString() confidenter.Identity {
	if f == nil {
		return confidenter.Identity{}
	}
	return confidenter.Identity{
		Domain: program.Domain(),
		Path:   "fount",
		ID:     strconv.FormatInt(f.ID, 10),
	}
}

type FountStat struct {
	ScannerStart  time.Time
	FountID       int64
	Start         time.Time
	Duration      int64
	ResponseError string
	LastItemError string
	ItemErrors    uint32
	ItemsTaken    uint32
	ItemsNew      uint32
}

type ScannerStat struct {
	Start      time.Time
	Duration   int64
	FountsNum  uint32
	ErrorsNum  uint32
	ItemsTaken uint32
	ItemsNew   uint32
}

type Operator interface {

	// Fount CRUD

	Create(identity *confidenter.Identity, fount Fount) (*confidenter.Identity, error)
	Read(identity *confidenter.Identity, toRead string) (*Fount, error)
	ReadAll(identity *confidenter.Identity, options *crud.ReadAllOptions, selector interfaces.Selector) ([]Fount, int64, error)
	Update(identity *confidenter.Identity, fount Fount) (crud.Result, error)
	Delete(identity *confidenter.Identity, toDelete string) (crud.Result, error)

	// Fount statistics

	AddFountStat(identity *confidenter.Identity, fountStat FountStat) error
	ReadAllFountStat(identity *confidenter.Identity, options *crud.ReadAllOptions, selector interfaces.Selector) ([]FountStat, int64, error)
	DeleteAllFountStat(identity *confidenter.Identity, selector interfaces.Selector) (crud.Result, error)

	// Scanner statistics

	AddScannerStat(identity *confidenter.Identity, scannerStat ScannerStat) error
	ReadAllScannerStat(identity *confidenter.Identity, options *crud.ReadAllOptions, selector interfaces.Selector) ([]ScannerStat, int64, error)
	DeleteAllScannerStat(identity *confidenter.Identity, selector interfaces.Selector) (crud.Result, error)

	// Fount tags
	//GetFountsForTag(identity confidenter.Identity, tag string) ([]int64, error)
	ReadTags(identity *confidenter.Identity, sel interfaces.Selector) ([]FountTag, error)

	ExportSettings(url, fountType, importParams string) (int64, error)

	// Closer

	Close()
}
