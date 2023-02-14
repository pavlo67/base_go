package scanner

import (
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/data/entities/records"
	"github.com/pavlo67/data/entities/sources"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/starter"
	founts "github.com/pavlo67/data/entities/sources/_senex"
)

const InterfaceKey = "scanner.comp"

func Starter() starter.Operator {
	return &scanComponent{}
}

var sourcesOp sources.Operator
var recordsOp records.Operator
var joinerOp joiner.Operator,
var l logger.Operator

type scanComponent struct{}

var _ starter.Operator = &scanComponent{}

func (sc *scanComponent) Name() string {
	return InterfaceKey
}

func (sc *scanComponent) Init() error {
	identity = program.Identity()

	var ok bool

	fountOp, ok = program.Interface(founts.InterfaceKey).(founts.Operator)
	if !ok {
		return errors.New("no fount interface found for scanner")
	}

	flowOp, ok = program.Interface(flow.InterfaceKey).(flow.Operator)
	if !ok {
		return errors.New("no fragment interface found for scanner")
	}

	return nil
}
