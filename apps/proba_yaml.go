package main

import (
	"log"

	"github.com/pavlo67/data/components/crud"

	"gopkg.in/yaml.v3"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"
)

func main() {
	actors := []auth.Actor{
		{
			Identity: &auth.Identity{
				Nickname: "pavlo",
				Roles:    rbac.Roles{rbac.RoleAdmin},
			},
			Creds: auth.Creds{
				auth.CredsNickname: "pavlo",
				auth.CredsPassword: "fotlM1mn",
			},
		},
		{
			Identity: &auth.Identity{
				Nickname: "crud_tester",
				Roles:    rbac.Roles{crud.RoleTester},
			},
			Creds: auth.Creds{
				auth.CredsNickname: string(crud.RoleTester),
				auth.CredsPassword: string(crud.RoleTester),
			},
		},
	}

	actorsYAML, _ := yaml.Marshal(actors)

	log.Printf("%s", actorsYAML)

}
