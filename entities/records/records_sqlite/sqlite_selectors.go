package records_sqlite

import (
	"fmt"
	"strings"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/data/components/selectors"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/data/components/tags"
	"github.com/pavlo67/data/entities/records"
)

const onAddParent = "on recordsSQLite.AddParent()"

func (recordsOp *recordsSQLite) AddParent(ts []tags.Item, id records.ID) ([]tags.Item, error) {
	idStr := strings.TrimSpace(string(id))
	if idStr == "" {
		return nil, errors.New(onAddParent + ": no id to add parent record")
	}

	return append(ts, tags.Item(idStr+":")), nil
}

const onConditions = "on records_sqlite.Conditions()"

func Conditions(selector *selectors.Term, options *auth.Identity) (string, []interface{}, error) {
	var condition string
	var values []interface{}

	if selector != nil {

		valuesStr, ok := selector.Values.([]string)
		if !ok {
			return "", nil, fmt.Errorf(onConditions+": wrong selector values (should be []string): %#v / %#v", selector, selector.Values)
		}

		switch selector.Key {
		case records.HasTag:
			if len(valuesStr) != 1 {
				return "", nil, fmt.Errorf(onConditions+": wrong values list in selector: %#v / %#v", selector, valuesStr)
			}
			tagStr := strings.TrimSpace(valuesStr[0])
			if tagStr == "" {
				return "", nil, errors.New(onConditions + ": no tag to select records")
			}
			condition = `tags LIKE ?`
			values = []interface{}{`%"` + tagStr + `"%`}

		case records.HasNoTag:
			condition = `tags IN ('', '{}')`

		case records.HasParent:
			if len(valuesStr) != 1 {
				return "", nil, fmt.Errorf(onConditions+": wrong values list in selector: %#v / %#v", selector, valuesStr)
			}
			idStr := strings.TrimSpace(valuesStr[0])
			if idStr == "" {
				return "", nil, errors.New(onConditions + ": no id to select records")
			}
			condition = `tags LIKE ?`
			values = []interface{}{`%"` + idStr + `:%`}

		default:
			return "", nil, fmt.Errorf(onConditions+": wrong selector.Key: %#v", selector)
		}
	}

	return condition, values, nil

}
