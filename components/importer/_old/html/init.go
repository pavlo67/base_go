package htmlimporter

import (
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces/starter"
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/interfaces/founts"
)

// Starter ...
func Starter() starter.Operator {
	return &htmlComponent{}
}

type htmlComponent struct {
}

var fountOp founts.Operator

const InterfaceKey = "importer.htmlimporter"

func (fl *htmlComponent) Name() string {
	return InterfaceKey
}

func (h *htmlComponent) Check(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string) ([]starter.Info, error) {
	return nil, nil
}

func (h *htmlComponent) Setup(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string, data map[string]string) error {
	return nil
}

func (h *htmlComponent) Init() error {

	var ok bool
	fountOp, ok = program.GetInterfaceBySignature((*founts.Operator)(nil), "").(founts.Operator)
	if !ok {
		return errors.New("can't get interface for fount component")
	}

	importer := &ImporterHTML{}
	err := program.JoinInterface(importer, InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join ruthenia as application")
	}

	return nil
}
