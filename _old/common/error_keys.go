package common

import (
	"errors"
)

type ErrorKey string

const CantPerformKey ErrorKey = "cant_perform"

const NoCredsKey ErrorKey = "no_creds"
const InvalidCredsKey ErrorKey = "invalid_creds"
const NoUserKey ErrorKey = "no_user"
const DuplicateUserKey ErrorKey = "duplicate_user"
const NoRightsKey ErrorKey = "no_rights"

const NotUniqueEmailKey ErrorKey = "not_unique_email"
const WrongPathKey ErrorKey = "wrong_path"
const WrongBodyKey ErrorKey = "wrong_body"
const WrongIDKey ErrorKey = "wrong_id"
const WrongJSONKey ErrorKey = "wrong_json"

const NotFoundKey ErrorKey = "not_found"
const NullItemKey ErrorKey = "null_item"
const NotImplementedKey ErrorKey = "not_implemented"
const NotSupportedKey ErrorKey = "not_supported"

// default errors --------------------------------------------------

var ErrNotImplemented = errors.New(string(NotImplementedKey))
var ErrNotFound = errors.New(string(NotFoundKey))
var ErrNullItem = errors.New(string(NullItemKey))
var ErrNotSupported = errors.New(string(NotSupportedKey))
