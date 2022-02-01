package main

import (
	"flag"
	"log"

	crud01_app2 "github.com/pavlo67/data/components/crud01/crud01_app"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/starter"
)

var BuildDate, BuildTag, BuildCommit string
var versionOnly bool

func main() {
	log.Printf("builded: %s, tag: %s, commit: %s\n", BuildDate, BuildTag, BuildCommit)
	flag.BoolVar(&versionOnly, "v", false, "show build vars only")
	flag.Parse()

	if versionOnly {
		return
	}

	cfgService, l := config.PrepareApp("_environments/")

	cfgTest, err := config.Get("_environments/test.yaml", config.MarshalerYAML)
	if err != nil || cfgTest == nil {
		l.Fatal(err)
	}

	starters, err := crud01_app2.Components(cfgService, *cfgTest, false)
	if err != nil {
		l.Fatal(err)
	}

	label := "NODE_CRUD/HTML/REST BUILD"
	joinerOp, err := starter.Run(starters, &cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	crud01_app2.WG.Wait()
}
