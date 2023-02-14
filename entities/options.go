package entities

import "github.com/pavlo67/common/common/rbac"

const RoleTester rbac.Role = "crud_tester"

type OptionsKey string

//type Item = Term

type Term struct {
	Key    OptionsKey
	Values interface{}
}

type Options struct {
	Selector *Term
	Ranges   *Ranges
}

type Ranges struct {
	GroupBy []string
	OrderBy []string
	JoinTo  string
	Values  []interface{}
	Offset  uint64
	Limit   uint64
}

func (options *Options) GetRanges() *Ranges {
	if options == nil {
		return nil
	}
	return options.Ranges
}

//func (options *Options) GetIdentity() *auth.Identity {
//	if options == nil {
//		return nil
//	}
//	return options.Identity
//}

func (options *Options) GetSelector() *Term {
	if options == nil {
		return nil
	}
	return options.Selector
}

func (options *Options) WithRanges(Ranges *Ranges) *Options {
	if options == nil {
		return &Options{Ranges: Ranges}
	}
	optionsCopied := *options
	options.Ranges = Ranges

	return &optionsCopied
}

//func (options *Options) HasRole(oneOfRoles ...rbac.Role) bool {
//	if options == nil || options.Identity == nil {
//		return false
//	}
//
//	return options.Identity.Roles.Has(oneOfRoles...)
//}

//func OptionsWithRoles(roles ...rbac.Role) *Options {
//	return &Options{
//		Identity: &auth.Identity{
//			Roles: roles,
//		},
//	}
//}
