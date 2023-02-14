package htmlimporter

import (
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/processor/sources"
	"github.com/pavlo67/punctum/starter/config"
)

// Starter ...
func Starter() starter.Operator {
	return &htmlComponent{}
}

type htmlComponent struct {
}

var fountOp sources.Operator

const InterfaceKey joiner.InterfaceKey = "importer.htmlimporter"

func (fl *htmlComponent) Name() string {
	return string(InterfaceKey)
}

func (h *htmlComponent) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {
	return nil, nil
}

func (h *htmlComponent) Setup(conf config.Config, indexPath string, data map[string]string) error {
	return nil
}

func (h *htmlComponent) Init() error {

	var ok bool
	fountOp, ok = joiner.GetInterfaceBySignature((*sources.Operator)(nil), "").(sources.Operator)
	if !ok {
		return errors.New("can't get interface for fount component")
	}

	importer := &ImporterHTML{}
	err := joiner.JoinInterface(importer, InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join ruthenia as application")
	}

	return nil
}
