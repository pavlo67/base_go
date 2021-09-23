package users_stub

//var _ users.Operator = &usersStub{}
//
//type UserStub struct {
//	Key          identity.Key
//	Nickname     string
//	Password     string
//	PasswordHash string
//}
//
//type usersStub struct {
//	users []UserStub
//	//crypter crypt.Crypter
//	//salt    string
//}
//
//const onNew = "on users_stub/New(): "
//
//func New(users []UserStub, salt string) (*usersStub, error) {
//	//crypter := crypt.SHA256.New()
//	//
//	//var err error
//	//for i, user := range users {
//	//	users[i].PasswordHash, err = crypter.Generate([]byte(user.Password), []byte(salt))
//	//	if err != nil {
//	//		return nil, errors.Wrap(err, onNew)
//	//	}
//	//}
//
//	for i, user := range users {
//		h := sha256.New()
//		h.Write([]byte(user.Password))
//		users[i].PasswordHash = string(h.Sum(nil))
//	}
//
//	return &usersStub{
//		users: users,
//		//crypter: crypter,
//		//salt:    salt,
//	}, nil
//}
//
//func (u *usersStub) CheckPassword(password, passHash string) bool {
//	h := sha256.New()
//	h.Write([]byte(password))
//	return passHash == string(h.Sum(nil))
//}
//
//const onSave = "on usersStub.Save(): "
//
//func (u *usersStub) Save(item users.Item, _ *crud.SaveOptions) (identity.Key, error) {
//	userStub := UserStub{
//		Nickname: item.Creds[auth.CredsNickname],
//		Password: item.Creds[auth.CredsPassword],
//	}
//
//	h := sha256.New()
//	h.Write([]byte(userStub.Password))
//	userStub.PasswordHash = string(h.Sum(nil))
//
//	//userStub.PasswordHash, err = u.crypter.Generate([]byte(userStub.Password), []byte(u.salt))
//	//if err != nil {
//	//	return "", errors.Wrap(err, onSave)
//	//}
//
//	if item.Key == "" {
//		userStub.Key = identity.Key(strconv.FormatInt(time.Now().UnixNano(), 10))
//		u.users = append(u.users, userStub)
//
//		return userStub.Key, nil
//	}
//	for i, us := range u.users {
//		if us.Key == userStub.Key {
//			u.users[i] = userStub
//			return userStub.Key, nil
//		}
//	}
//
//	return "", errors.Errorf(onSave+"no user with the same key as %#v", item)
//}
//
//const onList = "on usersStub.List(): "
//
//func (u *usersStub) List(selector *selectors.Term, _ *crud.GetOptions) ([]users.Item, error) {
//	var items []users.Item
//
//	if selector == nil {
//		for _, uu := range u.users {
//			items = append(items, users.Item{
//				User: auth.User{
//					Key:   uu.Key,
//					Creds: auth.Creds{auth.CredsNickname: uu.Nickname, auth.CredsPassword: uu.Password, auth.CredsPasshash: uu.PasswordHash},
//				},
//				Allowed: true,
//			})
//		}
//		return items, nil
//	}
//
//	if selector.Operation != selectors.Eq {
//		return nil, errors.Errorf(onList+"wrong selector.Operation: %#v", *selector)
//	}
//	if selector.Left != users.NicknameFieldName {
//		return nil, errors.Errorf(onList+"wrong selector.Left: %#v", *selector)
//	}
//
//	var value interface{}
//
//	switch v := selector.Right.(type) {
//	case selectors.Value:
//		value = v.V
//	case *selectors.Value:
//		value = v.V
//	default:
//		return nil, errors.Errorf(onList+"wrong selector.Right: %#v", *selector)
//	}
//
//	var valueStr string
//
//	switch v := value.(type) {
//	case string:
//		valueStr = v
//	case *string:
//		valueStr = *v
//	default:
//		return nil, errors.Errorf(onList+"wrong the right value of the selector: %#v", *selector)
//	}
//
//	for _, uu := range u.users {
//		if uu.Nickname == valueStr {
//			items = append(items, users.Item{
//				User: auth.User{
//					Key:   uu.Key,
//					Creds: auth.Creds{auth.CredsNickname: uu.Nickname, auth.CredsPassword: uu.Password, auth.CredsPasshash: uu.PasswordHash},
//				},
//				Allowed: true,
//			})
//		}
//	}
//
//	return items, nil
//}
//
//func (u *usersStub) Remove(identity.Key, *crud.RemoveOptions) error {
//	return common.ErrNotImplemented
//}
//
//func (u *usersStub) Read(identity.Key, *crud.GetOptions) (*users.Item, error) {
//	return nil, common.ErrNotImplemented
//}
//
//func (u *usersStub) Count(*selectors.Term, *crud.GetOptions) (uint64, error) {
//	return 0, common.ErrNotImplemented
//}
//
//func (u *usersStub) Allow() error {
//	return common.ErrNotImplemented
//}
//
//func (u *usersStub) SetVerification(auth.CredsType, string, bool) error {
//	return common.ErrNotImplemented
//}
//
//func (u *usersStub) Verify(auth.CredsType, string, common.Errors) error {
//	return common.ErrNotImplemented
//}
