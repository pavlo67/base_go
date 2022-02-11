package vcs

import (
	"encoding/json"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
)

const maxHistoryLength = 5

func ModifyHistory(history History, params common.Map) (historyChanged History, historyBytes, historyChangedBytes []byte, err error) {
	if len(history) > 0 {
		historyBytes, err = json.Marshal(history)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "can't marshal history (%#v)", history)
		}

		if len(history) >= maxHistoryLength {
			history = history[:maxHistoryLength-1]
		}
	}

	historyChanged = append(history, Action{
		// Actor:  "", // TODO!!!
		// Key:    "", // TODO!!!
		DoneAt: time.Now().UTC(),
		// Error:  nil, // TODO!!!
	})

	historyChangedBytes, err = json.Marshal(historyChanged)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "can't marshal changed history (%#v)", historyChanged)
	}

	return historyChanged, historyChangedBytes, historyBytes, nil
}
