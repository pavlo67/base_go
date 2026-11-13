package rbac

import (
	"encoding/json"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Roles []Role

func (roles Roles) Has(role ...Role) bool {
	for _, oneOf := range role {
		for _, r := range roles {
			if r == oneOf {
				//logPrintf("checked (%#v) as (%#v): true", role, roles)
				return true
			}
		}
	}

	//logPrintf("checked (%#v) as (%#v): false", role, roles)

	return false
}

func (roles Roles) Filter(role ...Role) Roles {
	var rolesFiltered Roles
	for _, r := range roles {
		for _, oneOf := range role {
			if r == oneOf {
				rolesFiltered = append(rolesFiltered, r)
			}
		}
	}

	return rolesFiltered
}

func (roles Roles) FilterNot(role ...Role) Roles {
	var rolesFiltered Roles

ROLES:
	for _, r := range roles {
		for _, oneOf := range role {
			if r == oneOf {
				continue ROLES
			}
		}
		rolesFiltered = append(rolesFiltered, r)
	}

	return rolesFiltered
}

func (roles Roles) MarshalJSON() ([]byte, error) {
	if len(roles) < 1 {
		return []byte("[]"), nil
	}

	var rs []string

	for _, r := range roles {
		rs = append(rs, string(r))
	}

	bytes, err := json.Marshal(rs)
	return bytes, err
}

func (roles Roles) ToStringList() []string {
	l := make([]string, 0, len(roles))
	for _, role := range roles {
		l = append(l, string(role))
	}
	return l
}
