package auth

type Creds map[CredsType]string

type CredsType string

const CredsEmail CredsType = "email"

const CredsToSet CredsType = "to_set"
const CredsRealm CredsType = "realm"

const CredsIP CredsType = "ip"

const CredsJWT CredsType = "jwt"
const CredsJWTRefresh CredsType = "jwt_refresh"

const CredsToken CredsType = "token"
const CredsPartnerToken CredsType = "partner_token"

const CredsID CredsType = "id"
const CredsLogin CredsType = "login"
const CredsNickname CredsType = "nickname"
const CredsPassword CredsType = "password"
const CredsTemporaryKey CredsType = "temporary_key"

const CredsQuestion CredsType = "question"
const CredsQuestionAnswer CredsType = "question_answer"

const CredsAllowedID CredsType = "allowed_id"

const CredsKeyToSignature CredsType = "key_to_signature"
const CredsSignature CredsType = "signature"
const CredsPublicKeyBase58 CredsType = "public_key_base58"
const CredsPublicKeyEncoding CredsType = "public_key_encoding"
const CredsPrivateKey CredsType = "private_key"

const CredsGroupID CredsType = "group_id"

const CredsPasshash CredsType = "passhash"
const CredsPasshashCryptype CredsType = "passhash_cryptype"

const CredsRole CredsType = "roles"
const CredsRolesJSON CredsType = "roles_json"

//func CheckCode() {
//	if toType == CredsSentCode {
//		// TODO: use "rector's" personIS to prevent second read
//		person, err := authOp.readByField("control_token", toFirstValue())
//		if err == nil {
//			err = authOp.personsOp.SetVerified(IS(person.ID), person.ID)
//		}
//		if err != nil {
//			return nil, nil, errors.Wrap(err, onAuthenticateWithCode)
//		}
//
//		return &User{
//			ID:       person.ID,
//			Nickname: "",
//			Accesses: nil,
//		}, nil, nil
//	}
//for _, ts := range toSet {
//if ts.Type == CredsSentCode {
//person.ControlToken = strlib.RandomString(persons.ControlTokenLength)
//_, err = authOp.personsOp.Update(IS(person.ID), *person)
//if err != nil {
//return nil, nil, errors.Wrap(err, onSendCode+": can't set control token")
//}
//
//if !person.Verified {
//return nil, nil, authOp.sendCode(Confirm, person.Email, person.ControlToken)
//}
//
//return nil, nil, authOp.sendCode(SendCode, person.Email, person.ControlToken)
//}
//}
//}

//const onAuthenticateWithCode = "on authPersons.RegisterWithCode"
//const onSendCode = "on authPersons.SendCode"

//const createMessageSubject = "Підтвердження реєстрації / Registration confirmation: "
//const createMessageContent = "Щоб завершити реєстрацію, будь ласка, пройдіть за ланкою / To complete registration please follow the link: "
//
//const sendCodeMessageSubject = "Підтвердження зміни паролю / Password updating confirmation: "
//const sendCodeMessageContent = "Щоб змінити пароль, будь ласка, пройдіть за ланкою / To change your password please follow the link: "
//

//func (authOp *authPersons) sendCode(cbKey Callback, email, token string) error {
//	cb, ok := authOp.callbacks[cbKey]
//	if !ok {
//		return errors.Errorf("wrong confirmation callback type: '%s'", cbKey)
//	}
//
//	message := connector.Message{
//		To:   email,
//		Body: token,
//	}
//
//	switch cbKey {
//	case Confirm:
//		message.Subject = createMessageSubject
//		if authOp.useMessages {
//			message.Body = createMessageContent + ` <a href="` + cb + url.QueryEscape(token) + `">підтвердити</a>`
//		}
//	case SendCode:
//		message.Subject = sendCodeMessageSubject
//		if authOp.useMessages {
//			message.Body = sendCodeMessageContent + ` <a href="` + cb + url.QueryEscape(token) + `">підтвердити</a>`
//		}
//	}
//
//	return authOp.senderOp.Set(message)
//}

//var reEmailToLogin1 = regexp.MustCompile(`@.*`)
//var reEmailToLogin2 = regexp.MustCompile(`(\.|-)`)

//type respFB struct {
//	TargetID    string `json:"id"`
//	Name  string `json:"name"`
//	Email string `json:"email"`
//}
//
//const onQueryPartnerUser = "on authPersons.QueryPartnerUser"
//
//func (authOp *authPersons) QueryPartnerUser(partnerKey, partnerToken string) (partnerUser *User, err error) {
//
//	// TODO: implement it as an interface call to "QueryPartnerUser" implementation according to partnerURL
//
//	partnerURL := "https://graph.facebook.com/me?fields=id,name,email&access_token="
//	req, err := http.NewRequest("GET", partnerURL+partnerToken, nil)
//
//	var resp *http.Response
//	if err == nil {
//		var client = &http.Client{}
//		resp, err = client.Do(req)
//		defer resp.Body.Close()
//	}
//	var data respFB
//	decoder := json.NewDecoder(resp.Body)
//	err = decoder.Decode(&data)
//	if err != nil {
//		return nil, errors.Wrapf(err, onQueryPartnerUser+": can't decode JSON from FB")
//	}
//
//	// log.Println("UserIS FB data:", data)
//
//	// TODO: use respFB.Label to control person if his email is changed
//
//	return &User{IdentityNamed: IDentityNamed{Nick: data.Name}, Contacts: data.Email}, nil
//}
//
//const onAuthenticateWithPartnerUser = "on authPersons.AuthenticateWithPartnerUser"
//
//func (authOp *authPersons) AuthenticateWithPartnerUser(partnerUser *User) (*User, error) {
//	person, err := authOp.readByLogin(partnerUser.Contacts)
//	if err != nil {
//		return nil, errors.Wrap(err, onAuthorize)
//	}
//
//	if person != nil {
//		person.Contacts = partnerUser.Contacts
//		return person, nil
//	}
//
//	partnerUser.TargetID, err = authOp.personsOp.Create(nil, partnerUser, true)
//	if err != nil {
//		return nil, errors.Wrap(err, onAuthenticateWithPartnerUser+": can't create person")
//	}
//
//	return partnerUser, nil
//}
//
//func (authOp *authPersons) Clean(selector selectors.Selector) error {
//	return authOp.personsOp.Clean(selector)
//}
//
//func (authOp *authPersons) ConsoleDo(action string, values map[string]string) error {
//
//	//// TODO:  refactor it
//	//
//	//var sqlQuery string
//	//var err error
//	//var stmt *sql.Stmt
//	//var res sql.Result
//	//var valuesAll []interface{}
//	//if action == "new_person" {
//	//	valuesAll = []interface{}{values["login"], values["password"], values["email"]}
//	//	sqlQuery = "insert into `" + u.personTable + "` (nickname, passhash, email, verified, contacts, history) values (?,?,?,1,'','')"
//	//} else if action == "new_password" {
//	//	valuesAll = []interface{}{values["password"], values["login"]}
//	//	sqlQuery = "update `" + u.personTable + "` set passhash=? where nickname=?"
//	//}
//	//
//	//if stmt, err = u.dbh.Init(sqlQuery); err != nil {
//	//	return nil.Wrapf(err, "can't prepare sql:%v", sqlQuery)
//	//}
//	//defer stmt.Close()
//	//
//	//if res, err = stmt.Exec(valuesAll...); err != nil {
//	//	return nil.Wrapf(err, "can't exec sql:%v, values=%v", sqlQuery, valuesAll)
//	//}
//	//iA, err := res.RowsAffected()
//	//
//	//if iA < 1 {
//	//	return nil.New("!!! RowsAffected() = ZERO !")
//	//}
//	//
//	return nil
//}

//
//const onUpdatePasswordWithCode = "on authPersons.UpdatePasswordWithCode"
//
//func (authOp *authPersons) UpdatePasswordWithCode(controlToken string, encodedNewPassword encrlib.EncryptedPass) error {
//
//	if controlToken == "" {
//		return errors.Wrap(confidenter.ErrTokenNotFound, onUpdatePasswordWithCode)
//	}
//	if encodedNewPassword.EncryptedPass == "" {
//		return errors.Wrap(ErrBadPassword, onUpdatePasswordWithCode)
//	}
//	if encodedNewPassword.Cryptype != encrlib.SHA256 {
//		return errors.Wrap(encrlib.ErrBadCryptype, onUpdatePasswordWithCode)
//	}
//
//	usrs, _, err := authOp.personsOp.ReadList(nil, &content.ListOptions{Selector: selectors.FieldEqual(persons.ControlTokenField, selectors.Value{controlToken})})
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": error finding person")
//	}
//	if len(usrs) != 1 {
//		return fmt.Errorf(onUpdatePasswordWithCode+": can't find person (find %d)", len(usrs))
//	}
//
//	personIS := &basis.UserIS{joiner.SystemDomain(), "/person/", usrs[0].TargetID}
//
//	// TODO: use correct domain/path/, prevent duplicate reading
//	person, err := authOp.personsOp.Read(personIS, usrs[0].TargetID)
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": error reading person")
//	}
//
//	person.Cryptype = encodedNewPassword.Cryptype
//	person.EncryptedPass = encodedNewPassword.EncryptedPass
//	_, err = authOp.personsOp.Update(personIS, person)
//	if err != nil {
//		return errors.Wrap(err, onUpdatePasswordWithCode+": can't set new password")
//	}
//
//	return nil
//}
