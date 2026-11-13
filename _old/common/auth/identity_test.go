package auth

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/rbac"
)

func TestJSON(t *testing.T) {
	identity := Identity{
		ID:       "1",
		Nickname: "2",
		Roles:    rbac.Roles{rbac.RoleUser},
	}

	bytes, err := json.Marshal(identity)
	require.NoError(t, err)

	var identity1 Identity
	err = json.Unmarshal(bytes, &identity1)
	require.NoError(t, err)

	require.Equal(t, identity.ID, identity1.ID)
	require.Equal(t, identity.Nickname, identity1.Nickname)
	require.Equal(t, identity.Roles, identity1.Roles)
	// require.Equal(t, identity.Creds(CredsPassword), identity1.Creds(CredsPassword))

	// log.Printf("%s / %s", bytes, err)

}
