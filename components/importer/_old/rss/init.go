package rss

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces/starter"
)

// Starter ...
func Starter() starter.Operator {
	return &rssComponent{}
}

type rssComponent struct {
	conf   config.PunctumConfig
	params map[string]string
}

const InterfaceKey = "importer.rss"

func (fl *rssComponent) Name() string {
	return InterfaceKey
}

func (rc *rssComponent) Check(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string) ([]starter.Info, error) {
	rc.conf = conf
	return nil, nil
}

func (rc *rssComponent) Setup(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string, data map[string]string) error {
	rc.conf = conf
	return nil
}

func (rc *rssComponent) Init() error {
	rssOp := &RSS{}

	err := program.JoinInterface(rssOp, InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join rss importer")
	}

	return nil
}
