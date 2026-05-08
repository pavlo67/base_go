package vcs

import (
	"encoding/json"
	"time"

	"github.com/pavlo67/base_go/lib/data"
	"github.com/pavlo67/base_go/lib/errors"
)

const maxHistoryLength = 5

func ModifyHistory(history History, params data.Map) (historyChanged History, historyOriginalStr, historyChangedStr string, err error) {
	if len(history) > 0 {
		historyOriginalBytes, err := json.Marshal(history)
		if err != nil {
			return nil, "", "", errors.Wrapf(err, "can't marshal history (%#v)", history)
		}

		if len(history) >= maxHistoryLength {
			history = history[:maxHistoryLength-1]
		}
		historyOriginalStr = string(historyOriginalBytes)
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

	return historyChanged, string(historyChangedBytes), historyOriginalStr, nil
}
