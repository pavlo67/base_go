package users_stub

//func Starter() starter.Operator {
//	return &users_stubStarter{}
//}
//
//var UserStubDefault = UserStub{
//	Key:      "pavlo/1",
//	Nickname: "pavlo",
//	Password: "fsamunp",
//}
//
//type users_stubStarter struct {
//	interfaceKey joiner.InterfaceKey
//
//	users []UserStub
//	// salt  string
//}
//
//var _ starter.Operator = &users_stubStarter{}
//var l logger.Operator
//
//var credentialsConf = map[string]string{}
//
//func (sc *users_stubStarter) Name() string {
//	return logger.GetCallInfo().PackageName
//}
//
//func (sc *users_stubStarter) Init(_ *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
//	if lCommon == nil {
//		return nil, errors.New("no logger.Operator")
//	}
//	l = lCommon
//
//	sc.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(users.InterfaceKey)))
//	sc.users = []UserStub{UserStubDefault}
//
//	return nil, nil
//}
//
//func (sc *users_stubStarter) Setup() error {
//	return nil
//}
//
//func (sc *users_stubStarter) Run(joiner joiner.Operator) error {
//	u, err := New(sc.users, "") // sc.salt
//	if err != nil {
//		return errors.Wrapf(err, "can't users_stubStarter.Run()")
//	}
//
//	err = joiner.Join(u, sc.interfaceKey)
//	if err != nil {
//		return errors.Wrap(err, "can't join *usersStub as users.Operator interface")
//	}
//
//	return nil
//}
