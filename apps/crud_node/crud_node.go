package main

import (
	"flag"
	"log"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/components/crud/crud_node"
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

	starters, err := crud_node.Starters(cfgService, *cfgTest, false)
	if err != nil {
		l.Fatal(err)
	}

	label := "NODE_CRUD/HTML/REST BUILD"
	joinerOp, err := starter.Run(starters, &cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	crud_node.WG.Wait()
}
