package main

import (
	"flag"
	"log"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/apps/node/node_settings"
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

	cfgService, l := config.Prepare("_environments/")
	label := "NODE/REST BUILD"
	joinerOp, err := starter.Run(node_settings.Components(true, false), &cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	node_settings.WG.Wait()
}
