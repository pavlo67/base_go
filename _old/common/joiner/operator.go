package joiner

type InterfaceKey string

type Component struct {
	InterfaceKey
	Interface interface{}
}

type Operator interface {
	Join(interface{}, InterfaceKey) error
	Interface(InterfaceKey) interface{}
	InterfacesAll(ptrToInterface interface{}) []Component
	CloseAll()
}

type Closer interface {
	Close() error
}
