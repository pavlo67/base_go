package joiner_runtime

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
)

var _ joiner.Operator = &joinerRuntime{}

func New(options common.Map, l logger.Operator) joiner.Operator {
	return &joinerRuntime{
		l:          l,
		options:    options,
		components: map[joiner.InterfaceKey]interface{}{},
		mutex:      &sync.RWMutex{},
	}
}

type joinerRuntime struct {
	l          logger.Operator
	options    common.Map
	components map[joiner.InterfaceKey]interface{}
	mutex      *sync.RWMutex
}

var ErrJoiningOnEmptyKey = errors.New("can't join on empty interface key")
var ErrJoiningNil = errors.New("can't join nil interface")
var ErrJoiningDuplicate = errors.New("can't join interface over joined before")

func (j *joinerRuntime) Join(intrfc interface{}, interfaceKey joiner.InterfaceKey) error {
	if j == nil {
		return fmt.Errorf("got nil on .Join(%s)", interfaceKey)
	}
	if interfaceKey == "" {
		return errors.Wrapf(ErrJoiningNil, "on .Join(%s)", interfaceKey)
	}
	if intrfc == nil {
		return errors.Wrapf(ErrJoiningNil, "on .Join(%s)", interfaceKey)
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()

	if _, ok := j.components[interfaceKey]; ok {
		return errors.Wrapf(ErrJoiningDuplicate, "on .Join(%s)", interfaceKey)
	}

	j.components[interfaceKey] = intrfc

	if j.l != nil && !j.options.IsTrue("silent") {
		j.l.Infof("joined (%T) as %s", intrfc, interfaceKey)
	}

	return nil
}

func (j *joinerRuntime) Interface(interfaceKey joiner.InterfaceKey) interface{} {
	if j == nil {
		log.Printf("on ActorKey.Component(%s): null ActorKey item", interfaceKey)
	}

	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if intrfc, ok := j.components[interfaceKey]; ok {
		return intrfc
	}

	return nil
}

func (j *joinerRuntime) InterfacesAll(ptrToInterface interface{}) []joiner.Component {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	var components []joiner.Component

	for key, intrfc := range j.components {
		if CheckInterface(intrfc, ptrToInterface) {
			components = append(components, joiner.Component{InterfaceKey: key, Interface: intrfc})
		}
	}

	return components
}

func CheckInterface(intrfc interface{}, ptrToInterface interface{}) bool {
	defer func() {
		recover()
	}()

	// ??? reflect.TypeOf(ptrToInterface).Elem()
	// ??? if intrfc != nil
	if reflect.TypeOf(intrfc).Implements(reflect.TypeOf(ptrToInterface).Elem()) {
		return true
	}

	return false
}

func (j *joinerRuntime) CloseAll() {
	if j == nil {
		log.Print("on joinerRuntime.CloseAll(): j == nil")
		return
	}

	closerComponents := j.InterfacesAll((*joiner.Closer)(nil))

	for _, closerComponent := range closerComponents {
		if closer, _ := closerComponent.Interface.(joiner.Closer); closer != nil {
			err := closer.Close()
			if err != nil {
				log.Print("on joinerRuntime.CloseAll(): ", err)
			}
		}
	}

}
