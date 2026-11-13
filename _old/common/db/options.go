package db

//type Identity struct {
//	Identity *auth.Identity
//	Selector *selectors.Term
//	// Tx       interface{}
//}
//
//func (options *Identity) GetIdentity() *auth.Identity {
//	if options == nil {
//		return nil
//	}
//	return options.Identity
//}
//
//func (options *Identity) GetSelector() *selectors.Term {
//	if options == nil {
//		return nil
//	}
//	return options.Selector
//}
//
//func (options *Identity) HasRole(oneOfRoles ...rbac.Role) bool {
//	if options == nil || options.Identity == nil {
//		return false
//	}
//
//	return options.Identity.Roles.Has(oneOfRoles...)
//}

//func (options *Identity) WithSelector(selector selectors.Term) *Identity {
//	if options == nil {
//		return &Identity{Selector: &selector}
//	}
//	optionsCopied := *options
//	optionsCopied.Selector = &selector
//
//	return &optionsCopied
//}
//func OptionsWithRoles(roles ...rbac.Role) *Identity {
//	return &Identity{
//		Identity: &auth.Identity{
//			Roles: roles,
//		},
//	}
//}
