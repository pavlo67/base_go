package sources_stub

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/workshop/components/sources"
)

func Starter() starter.Operator {
	return &sourcesStubStarter{}
}

var l logger.Operator
var _ starter.Operator = &sourcesStubStarter{}

type sourcesStubStarter struct {
	interfaceKey joiner.InterfaceKey
	// cleanerInterfaceKey joiner.HandlerKey

}

func (rs *sourcesStubStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (rs *sourcesStubStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	rs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(sources.InterfaceKey)))
	// rs.cleanerInterfaceKey = joiner.HandlerKey(options.StringDefault("cleaner_interface_key", string(sources.CleanerInterfaceKey)))

	return nil, nil
}

func (rs *sourcesStubStarter) Setup() error {
	return nil
}

func (rs *sourcesStubStarter) Run(joinerOp joiner.Operator) error {
	urls := []string{
		"https://rss.unian.net/site/news_ukr.rss",
		"https://censor.net.ua/includes/news_uk.xml",
		"http://texty.org.ua/mod/news/?view=rss",
		"http://texty.org.ua/mod/article/?view=rss&ed=1",
		"http://texty.org.ua/mod/blog/blog_list.php?view=rss",
		"https://www.pravda.com.ua/rss/",
		"http://k.img.com.ua/rss/ua/all_news2.0.xml",
		"https://www.obozrevatel.com/rss.xml",
		"https://lenta.ru/rss",
		"https://www.gazeta.ru/export/rss/first.xml",
		"https://www.gazeta.ru/export/rss/lenta.xml",
	}

	sourcesOp, _, err := New(urls)
	if err != nil {
		return errors.Wrap(err, "can't init *sourcesStub")
	}

	err = joinerOp.Join(sourcesOp, rs.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *sourcesStub as sources.ActorKey with key '%s'", rs.interfaceKey)
	}

	return nil
}
