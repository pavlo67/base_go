package vcs

import (
	"encoding/json"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
)

const maxHistoryLength = 5

func ModifyHistoryStr(history History, params common.Map) (historyChanged History, historyStr, historyChangedStr string, err error) {
	if len(history) > 0 {
		historyBytes, err := json.Marshal(history)
		if err != nil {
			return nil, "", "", errors.Wrapf(err, "can't marshal history (%#v)", history)
		}
		historyStr = string(historyBytes)

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

	historyChangedBytes, err := json.Marshal(historyChanged)
	if err != nil {
		return nil, "", "", errors.Wrapf(err, "can't marshal changed history (%#v)", historyChanged)
	}

	return historyChanged, string(historyChangedBytes), historyStr, nil
}
