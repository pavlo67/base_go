package persons_sqlite

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db/db_sqlite"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/persons"
)

func TestPersonsSQLite(t *testing.T) {
	err := os.Setenv("SHOW_CONNECTS", "1")
	require.NoError(t, err)

	cfgService, l := config.PrepareTests(
		t,
		"../../../_environments/",
		"test",
		"", // "persons_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
	)

	components := []starter.Starter{
		{db_sqlite.Starter(), nil},
		{Starter(), nil},
	}

	label := "PERSONS_SQLITE/TEST BUILD"
	joinerOp, err := starter.Run(components, &cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	personsOp, _ := joinerOp.Interface(persons.InterfaceKey).(persons.Operator)
	require.NotNil(t, personsOp)

	personsCleanerOp, _ := joinerOp.Interface(persons.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, personsCleanerOp)

	persons.OperatorTestScenario(t, personsOp, personsCleanerOp)
}
